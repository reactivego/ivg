package generate

import (
	"image/color"
	"math"
	"strconv"

	"github.com/reactivego/ivg"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	CSELUsedAsBothGradientAndStop = Error("ivg: CSEL used as both gradient and stop")
	TooManyGradientStops          = Error("ivg: too many gradient stops")
	UnrecognizedPathDataVerb      = Error("ivg: unrecognized path data verb")
)

// Aff3 is a 3x3 affine transformation matrix in row major order, where the
// bottom row is implicitly [0 0 1].
//
// m[3*r + c] is the element in the r'th row and c'th column.
type Aff3 [6]float32

// GradientShape is the gradient shape.
type GradientShape uint8

const (
	GradientShapeLinear GradientShape = iota
	GradientShapeRadial
)

// GradientSpread is how to spread a gradient past its nominal bounds (from
// offset being 0.0 to offset being 1.0).
type GradientSpread uint8

const (
	GradientSpreadNone GradientSpread = iota
	GradientSpreadPad
	GradientSpreadReflect
	GradientSpreadRepeat
)

// GradientStop is a color/offset gradient stop.
type GradientStop struct {
	Offset float32
	Color  color.Color
}

type Generator struct {
	ivg.Destination
}

func (g *Generator) SetDestination(d ivg.Destination) {
	g.Destination = d
}

// SetLinearGradient is like SetGradient with shape=ShapeLinear except that the
// transformation matrix is implicitly defined by two boundary points (x1, y1)
// and (x2, y2).
func (g *Generator) SetLinearGradient(x1, y1, x2, y2 float32, spread GradientSpread, stops []GradientStop) error {
	// See the package documentation's appendix for a derivation of the
	// transformation matrix.
	dx, dy := x2-x1, y2-y1
	d := dx*dx + dy*dy
	ma := dx / d
	mb := dy / d
	vbx2grad := Aff3{
		ma, mb, -ma*x1 - mb*y1,
		0, 0, 0,
	}
	return g.SetGradient(GradientShapeLinear, spread, stops, vbx2grad)
}

// SetCircularGradient is like SetGradient with radial=true except that the
// transformation matrix is implicitly defined by a center (cx, cy) and a
// radius vector (rx, ry) such that (cx+rx, cy+ry) is on the circle.
func (g *Generator) SetCircularGradient(cx, cy, rx, ry float32, spread GradientSpread, stops []GradientStop) error {
	// See the package documentation's appendix for a derivation of the
	// transformation matrix.
	invR := float32(1 / math.Sqrt(float64(rx*rx+ry*ry)))
	vbx2grad := Aff3{
		invR, 0, -cx * invR,
		0, invR, -cy * invR,
	}
	return g.SetGradient(GradientShapeRadial, spread, stops, vbx2grad)
}

// SetEllipticalGradient is like SetGradient with radial=true except that the
// transformation matrix is implicitly defined by a center (cx, cy) and two
// axis vectors (rx, ry) and (sx, sy) such that (cx+rx, cy+ry) and (cx+sx,
// cy+sy) are on the ellipse.
func (d *Generator) SetEllipticalGradient(cx, cy, rx, ry, sx, sy float32, spread GradientSpread, stops []GradientStop) error {
	// Explicitly disable FMA in the floating-point calculations below
	// to get consistent results on all platforms, and in turn produce
	// a byte-identical encoding.
	// See https://golang.org/ref/spec#Floating_point_operators and issue 43219.

	// See the package documentation's appendix for a derivation of the
	// transformation matrix.
	invRSSR := 1 / (float32(rx*sy) - float32(sx*ry))

	ma := +sy * invRSSR
	mb := -sx * invRSSR
	mc := -float32(ma*cx) - float32(mb*cy)
	md := -ry * invRSSR
	me := +rx * invRSSR
	mf := -float32(md*cx) - float32(me*cy)

	vbx2grad := Aff3{
		ma, mb, mc,
		md, me, mf,
	}
	return d.SetGradient(GradientShapeRadial, spread, stops, vbx2grad)
}

// SetGradient sets CREG[CSEL] to encode the gradient whose colors defined by
// spread and stops. Its geometry is either linear or radial, depending on the
// radial argument, and the given affine transformation matrix maps from
// graphic coordinate space defined by the metadata's viewBox (e.g. from (-32,
// -32) to (+32, +32)) to gradient coordinate space. Gradient coordinate space
// is where a linear gradient ranges from x=0 to x=1, and a radial gradient has
// center (0, 0) and radius 1.
//
// The colors of the n stops are encoded at CREG[cBase+0], CREG[cBase+1], ...,
// CREG[cBase+n-1]. Similarly, the offsets of the n stops are encoded at
// NREG[nBase+0], NREG[nBase+1], ..., NREG[nBase+n-1]. Additional parameters
// are stored at NREG[nBase-4], NREG[nBase-3], NREG[nBase-2] and NREG[nBase-1].
//
// The CSEL and NSEL selector registers maintain the same values after the
// method returns as they had when the method was called.
//
// See the package documentation for more details on the gradient encoding
// format and the derivation of common transformation matrices.
func (d *Generator) SetGradient(shape GradientShape, spread GradientSpread, stops []GradientStop, transform Aff3) error {
	cBase, nBase := uint8(10), uint8(10)

	nStops := uint8(len(stops))
	if nStops > uint8(64-len(transform)) {
		return TooManyGradientStops
	}
	if x, y := d.CSel(), d.CSel()+64; (cBase <= x && x < cBase+nStops) || (cBase <= y && y < cBase+nStops) {
		return CSELUsedAsBothGradientAndStop
	}

	oldCSel := d.CSel()
	oldNSel := d.NSel()
	d.SetCReg(0, false, ivg.RGBAColor(ivg.EncodeGradient(cBase, nBase, uint8(shape), uint8(spread), nStops)))
	d.SetCSel(cBase)
	d.SetNSel(nBase)
	for i, v := range transform {
		d.SetNReg(uint8(len(transform)-i), false, v)
	}
	for _, s := range stops {
		r, g, b, a := s.Color.RGBA()
		d.SetCReg(0, true, ivg.RGBAColor(color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}))
		d.SetNReg(0, true, s.Offset)
	}
	d.SetCSel(oldCSel)
	d.SetNSel(oldNSel)
	return nil
}

func (e *Generator) SetPathData(d string, adj uint8, normalizeTo64X64 bool) error {
	var args [7]float32
	prevN, prevVerb := 0, byte(0)
	for first := true; d != "z"; first = false {
		n, verb, implicit := 0, d[0], false
		switch d[0] {
		case 'H', 'h', 'V', 'v':
			n = 1
		case 'L', 'M', 'l', 'm':
			n = 2
		case 'S', 's':
			n = 4
		case 'C', 'c':
			n = 6
		case 'A', 'a':
			n = 7
		case 'z':
			n = 0
		default:
			if prevVerb == '\x00' {
				return UnrecognizedPathDataVerb
			}
			n, verb, implicit = prevN, prevVerb, true
		}
		prevN, prevVerb = n, verb
		if prevVerb == 'M' {
			prevVerb = 'L'
		} else if prevVerb == 'm' {
			prevVerb = 'l'
		}
		if !implicit {
			d = d[1:]
		}

		for i := 0; i < n; i++ {
			nDots := 0
			if d[0] == '.' {
				nDots = 1
			}
			j := 1
			for ; ; j++ {
				switch d[j] {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					continue
				case '.':
					nDots++
					if nDots == 1 {
						continue
					}
				}
				break
			}
			f, err := strconv.ParseFloat(d[:j], 64)
			if err != nil {
				return err
			}
			args[i] = float32(f)
			for ; d[j] == ' ' || d[j] == ','; j++ {
			}
			d = d[j:]
		}

		if normalizeTo64X64 {
			// The original SVG is 32x32 units, with the top left being (0, 0).
			// Normalize to 64x64 units, with the center being (0, 0).
			if verb == 'A' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[2] /= 360
				args[5] = 2*args[5] - 32
				args[6] = 2*args[6] - 32
			} else if verb == 'a' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[2] /= 360
				args[5] = 2 * args[5]
				args[6] = 2 * args[6]
			} else if first || ('A' <= verb && verb <= 'Z') {
				for i := range args {
					args[i] = 2*args[i] - 32
				}
			} else {
				for i := range args {
					args[i] = 2 * args[i]
				}
			}
		} else if verb == 'A' || verb == 'a' {
			args[2] /= 360
		}

		if first {
			first = false
			e.StartPath(adj, args[0], args[1])
			continue
		}
		switch verb {
		case 'H':
			e.AbsHLineTo(args[0])
		case 'h':
			e.RelHLineTo(args[0])
		case 'V':
			e.AbsVLineTo(args[0])
		case 'v':
			e.RelVLineTo(args[0])
		case 'L':
			e.AbsLineTo(args[0], args[1])
		case 'l':
			e.RelLineTo(args[0], args[1])
		case 'm':
			e.ClosePathRelMoveTo(args[0], args[1])
		case 'S':
			e.AbsSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 's':
			e.RelSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 'C':
			e.AbsCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'c':
			e.RelCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'A':
			e.AbsArcTo(args[0], args[1], args[2], args[3] != 0, args[4] != 0, args[5], args[6])
		case 'a':
			e.RelArcTo(args[0], args[1], args[2], args[3] != 0, args[4] != 0, args[5], args[6])
		case 'z':
			// No-op.
		default:
			return UnrecognizedPathDataVerb
		}
	}
	e.ClosePathEndPath()
	return nil
}
