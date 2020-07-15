package generate

import (
	"image/color"
	"math"

	"github.com/reactivego/ivg"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	CSELUsedAsBothGradientAndStop = Error("ivg: CSEL used as both gradient and stop")
	TooManyGradientStops          = Error("ivg: too many gradient stops")
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
	// See the package documentation's appendix for a derivation of the
	// transformation matrix.
	invRSSR := 1 / (rx*sy - sx*ry)

	ma := +sy * invRSSR
	mb := -sx * invRSSR
	mc := -ma*cx - mb*cy
	md := -ry * invRSSR
	me := +rx * invRSSR
	mf := -md*cx - me*cy

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
