package generate

import (
	"fmt"
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
)

func UnrecognizedPathDataVerb(verb byte) Error {
	return Error(fmt.Sprintf("ivg: unrecognized path data verb (%c)", verb))
}

// Aff3 is a 3x3 affine transformation matrix in row major order, where the
// bottom row is implicitly [0 0 1].
//
// m[3*r + c] is the element in the r'th row and c'th column.
type Aff3 [6]float32

func Translate(x, y float32) Aff3 {
	return Aff3{
		1, 0, x,
		0, 1, y,
	}
}

func Scale(v ...float32) Aff3 {
	switch len(v) {
	case 0:
		return Aff3{1, 0, 0, 0, 1, 0}
	case 1:
		return Aff3{v[0], 0, 0, 0, v[0], 0}
	default:
		return Aff3{v[0], 0, 0, 0, v[1], 0}
	}
}

func Concat(affs ...Aff3) Aff3 {
	switch len(affs) {
	case 0:
		return Aff3{1, 0, 0, 0, 1, 0}
	case 1:
		return affs[0]
	default:
		a := Aff3{1, 0, 0, 0, 1, 0}
		for _, b := range affs {
			a = Aff3{
				a[0]*b[0] + a[3]*b[1], a[1]*b[0] + a[4]*b[1], a[2]*b[0] + a[5]*b[1] + b[2],
				a[0]*b[3] + a[3]*b[4], a[1]*b[3] + a[4]*b[4], a[2]*b[3] + a[5]*b[4] + b[5],
			}
		}
		return a
	}
}

func MulAff3(x, y float32, a Aff3) (X, Y float32) {
	return x*a[0] + y*a[1] + a[2], x*a[3] + y*a[4] + a[5]
}

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
	transforms []Aff3
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

func (e *Generator) SetTransform(transforms ...Aff3) {
	e.transforms = []Aff3{Concat(transforms...)}
}

func (e *Generator) SetPathData(d string, adj uint8) error {
	var args [7]float32
	prevN, prevVerb := 0, byte(0)
	for start := true; d != "z"; start = false {
		// The verb at the start of the path data in d must be either 'M' or 'm'
		// A terminating 'Z' or 'z' is optional and only makes a difference for stroking.
		n, verb, implicit := 0, d[0], false
		switch verb {
		case 'H', 'h', 'V', 'v':
			n = 1
		case 'L', 'l', 'M', 'm', 'T', 't':
			n = 2
		case 'Q', 'q', 'S', 's':
			n = 4
		case 'C', 'c':
			n = 6
		case 'A', 'a':
			n = 7
		case 'Z', 'z':
			n = 0
		default:
			if prevVerb == '\x00' {
				return UnrecognizedPathDataVerb(verb)
			}
			n, verb, implicit = prevN, prevVerb, true
		}
		prevN, prevVerb = n, verb
		if prevVerb == 'M' {
			prevVerb = 'L'
		} else if prevVerb == 'm' {
			prevVerb = 'l'
		}
		if start {
			verb = '@'
		}
		if !implicit {
			d = d[1:]
		}

		if dnext, err := scan(&args, d, n); err != nil {
			return err
		} else {
			d = dnext
		}
		normalize(&args, n, verb, e.transforms...)

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
		case '@':
			e.StartPath(adj, args[0], args[1])
		case 'M':
			e.ClosePathAbsMoveTo(args[0], args[1])
		case 'm':
			e.ClosePathRelMoveTo(args[0], args[1])
		case 'T':
			e.AbsSmoothQuadTo(args[0], args[1])
		case 't':
			e.RelSmoothQuadTo(args[0], args[1])
		case 'Q':
			e.AbsQuadTo(args[0], args[1], args[2], args[3])
		case 'q':
			e.RelQuadTo(args[0], args[1], args[2], args[3])
		case 'S':
			e.AbsSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 's':
			e.RelSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 'C':
			e.AbsCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'c':
			e.RelCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'A':
			e.AbsArcTo(args[0], args[1], args[2]/360, args[3] != 0, args[4] != 0, args[5], args[6])
		case 'a':
			e.RelArcTo(args[0], args[1], args[2]/360, args[3] != 0, args[4] != 0, args[5], args[6])
		case 'Z', 'z':
			// No-op.
		default:
			return UnrecognizedPathDataVerb(verb)
		}
	}
	e.ClosePathEndPath()
	return nil
}

func scan(args *[7]float32, d string, n int) (string, error) {
	for i := 0; i < n; i++ {
		nDots := 0
		if d[0] == '.' {
			nDots = 1
		}
		j := 1 // skip over a '+' or '-' or any other character for that matter
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
			return d, err
		}
		args[i] = float32(f)
		for ; d[j] == ' ' || d[j] == ','; j++ {
		}
		d = d[j:]
	}
	return d, nil
}

func normalize(args *[7]float32, n int, verb byte, transforms ...Aff3) {
	if len(transforms) > 0 {
		// The original SVG is 32x32 units, with the top left being (0, 0).
		// Normalize to 64x64 units, with the center being (0, 0).
		if true {
			transform := Concat(transforms...)
			scale := Aff3{transform[0], 0, 0, 0, transform[4], 0}
			if 'a' <= verb && verb <= 'z' {
				transform = scale
			}
			switch n {
			case 7:
				args[0], args[1] = MulAff3(args[0], args[1], scale)
				args[5], args[6] = MulAff3(args[5], args[6], transform)
			case 6:
				args[4], args[5] = MulAff3(args[4], args[5], transform)
				fallthrough
			case 4:
				args[2], args[3] = MulAff3(args[2], args[3], transform)
				fallthrough
			case 2:
				args[0], args[1] = MulAff3(args[0], args[1], transform)
			case 1:
				if verb == 'H' || verb == 'h' {
					args[0], _ = MulAff3(args[0], 0, transform)
				} else if verb == 'V' || verb == 'v' {
					_, args[0] = MulAff3(0, args[0], transform)
				}
			}
		} else {
			if verb == 'A' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[5] = 2*args[5] - 32
				args[6] = 2*args[6] - 32
			} else if verb == 'a' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[5] = 2 * args[5]
				args[6] = 2 * args[6]
			} else if '@' <= verb && verb <= 'Z' {
				for i := range args {
					args[i] = 2*args[i] - 32
				}
			} else {
				for i := range args {
					args[i] = 2 * args[i]
				}
			}
		}
		// fmt.Println(args[:n])
	}
}

func (e *Generator) StartPath(adj uint8, x, y float32) {
	e.Destination.StartPath(adj, x, y)
}

func (e *Generator) ClosePathEndPath() {
	e.Destination.ClosePathEndPath()
}

func (e *Generator) ClosePathAbsMoveTo(x, y float32) {
	e.Destination.ClosePathAbsMoveTo(x, y)
}

func (e *Generator) ClosePathRelMoveTo(x, y float32) {
	e.Destination.ClosePathRelMoveTo(x, y)
}

func (e *Generator) AbsHLineTo(x float32) {
	e.Destination.AbsHLineTo(x)
}

func (e *Generator) RelHLineTo(x float32) {
	e.Destination.RelHLineTo(x)
}

func (e *Generator) AbsVLineTo(y float32) {
	e.Destination.AbsVLineTo(y)
}

func (e *Generator) RelVLineTo(y float32) {
	e.Destination.RelVLineTo(y)
}

func (e *Generator) AbsLineTo(x, y float32) {
	e.Destination.AbsLineTo(x, y)
}

func (e *Generator) RelLineTo(x, y float32) {
	e.Destination.RelLineTo(x, y)
}

func (e *Generator) AbsSmoothQuadTo(x, y float32) {
	e.Destination.AbsSmoothQuadTo(x, y)
}

func (e *Generator) RelSmoothQuadTo(x, y float32) {
	e.Destination.RelSmoothQuadTo(x, y)
}

func (e *Generator) AbsQuadTo(x1, y1, x, y float32) {
	e.Destination.AbsQuadTo(x1, y1, x, y)
}

func (e *Generator) RelQuadTo(x1, y1, x, y float32) {
	e.Destination.RelQuadTo(x1, y1, x, y)
}

func (e *Generator) AbsSmoothCubeTo(x2, y2, x, y float32) {
	e.Destination.AbsSmoothCubeTo(x2, y2, x, y)
}

func (e *Generator) RelSmoothCubeTo(x2, y2, x, y float32) {
	e.Destination.RelSmoothCubeTo(x2, y2, x, y)
}

func (e *Generator) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	e.Destination.AbsCubeTo(x1, y1, x2, y2, x, y)
}

func (e *Generator) RelCubeTo(x1, y1, x2, y2, x, y float32) {
	e.Destination.RelCubeTo(x1, y1, x2, y2, x, y)
}

func (e *Generator) AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	e.Destination.AbsArcTo(rx, ry, xAxisRotation, largeArc, sweep, x, y)
}

func (e *Generator) RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	e.Destination.RelArcTo(rx, ry, xAxisRotation, largeArc, sweep, x, y)
}
