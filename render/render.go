// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render

import (
	"image"
	"image/color"
	"math"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/raster"
)

var (
	negativeInfinity = math.Float32frombits(0xff800000)
	positiveInfinity = math.Float32frombits(0x7f800000)
)

const (
	smoothTypeNone = iota
	smoothTypeQuad
	smoothTypeCube
)

// Renderer implements the ivg.Destination interface to render an IconVG graphic
// onto a ivg/raster.Rasterizer.
//
// Call SetRasterizer to change the rasterizer, before calling Decode or
// between calls to Decode.
type Renderer struct {
	z raster.Rasterizer
	r image.Rectangle

	// scale and bias transforms the viewBox rectangle to the (0, 0) - (r.Dx(),
	// r.Dy()) rectangle.
	scaleX float32
	biasX  float32
	scaleY float32
	biasY  float32

	viewBox ivg.ViewBox
	// Palette is a 64 color palette. When encoding, it is the suggested
	// palette to place within the IconVG graphic. When decoding, it is either
	// the optional palette passed to Decode, or if no optional palette was
	// given, the suggested palette within the IconVG graphic.
	palette [64]color.RGBA

	lod0 float32
	lod1 float32
	cSel uint8
	nSel uint8

	disabled bool

	prevSmoothType   uint8
	prevSmoothPointX float32
	prevSmoothPointY float32

	fill      image.Image
	flatColor color.RGBA
	flatImage image.Uniform
	gradient  Gradient

	cReg  [64]color.RGBA
	nReg  [64]float32
	stops [64]Stop
}

// SetRasterizer sets the rasterizer to draw into.
// The IconVG graphic (which does not have a fixed size in pixels) will be
// scaled in the X and Y dimensions to fit the rectangle r. The scaling factors
// may differ in the two dimensions.
func (z *Renderer) SetRasterizer(c raster.Rasterizer, r image.Rectangle) {
	z.z = c
	if r.Empty() {
		r = image.Rectangle{}
	}
	z.r = r
	z.recalcTransform()
}

// Reset resets the Destination for the given Metadata.
func (z *Renderer) Reset(viewbox ivg.ViewBox, palette *[64]color.RGBA) {
	z.viewBox = viewbox
	z.palette = *palette
	z.lod0 = 0
	z.lod1 = positiveInfinity
	z.cSel = 0
	z.nSel = 0
	z.prevSmoothType = smoothTypeNone
	z.prevSmoothPointX = 0
	z.prevSmoothPointY = 0
	z.cReg = *palette
	z.nReg = [64]float32{}
	z.recalcTransform()
}

func (z *Renderer) recalcTransform() {
	z.scaleX = float32(z.r.Dx()) / (z.viewBox.MaxX - z.viewBox.MinX)
	z.biasX = -z.viewBox.MinX
	z.scaleY = float32(z.r.Dy()) / (z.viewBox.MaxY - z.viewBox.MinY)
	z.biasY = -z.viewBox.MinY
}

func (z *Renderer) CSel() uint8 {
	return z.cSel
}

func (z *Renderer) SetCSel(cSel uint8) {
	z.cSel = cSel & 0x3f
}

func (z *Renderer) NSel() uint8 {
	return z.nSel
}

func (z *Renderer) SetNSel(nSel uint8) {
	z.nSel = nSel & 0x3f
}

func (z *Renderer) SetCReg(adj uint8, incr bool, c ivg.Color) {
	z.cReg[(z.cSel-adj)&0x3f] = c.Resolve(&z.palette, &z.cReg)
	if incr {
		z.cSel++
	}
}

func (z *Renderer) SetNReg(adj uint8, incr bool, f float32) {
	z.nReg[(z.nSel-adj)&0x3f] = f
	if incr {
		z.nSel++
	}
}

func (z *Renderer) SetLOD(lod0, lod1 float32) {
	z.lod0, z.lod1 = lod0, lod1
}

func (z *Renderer) unabsX(x float32) float32 { return x/z.scaleX - z.biasX }
func (z *Renderer) unabsY(y float32) float32 { return y/z.scaleY - z.biasY }
func (z *Renderer) absX(x float32) float32   { return z.scaleX * (x + z.biasX) }
func (z *Renderer) absY(y float32) float32   { return z.scaleY * (y + z.biasY) }
func (z *Renderer) relX(x float32) float32   { return z.scaleX * x }
func (z *Renderer) relY(y float32) float32   { return z.scaleY * y }

func (z *Renderer) absVec2(x, y float32) (zx, zy float32) {
	return z.absX(x), z.absY(y)
}

func (z *Renderer) relVec2(x, y float32) (zx, zy float32) {
	px, py := z.z.Pen()
	return px + z.relX(x), py + z.relY(y)
}

// implicitSmoothPoint returns the implicit control point for smooth-quadratic
// and smooth-cubic Bézier curves.
//
// https://www.w3.org/TR/SVG/paths.html#PathDataCurveCommands says, "The first
// control point is assumed to be the reflection of the second control point on
// the previous command relative to the current point. (If there is no previous
// command or if the previous command was not [a quadratic or cubic command],
// assume the first control point is coincident with the current point.)"
func (z *Renderer) implicitSmoothPoint(thisSmoothType uint8) (zx, zy float32) {
	px, py := z.z.Pen()
	if z.prevSmoothType != thisSmoothType {
		return px, py
	}
	return 2*px - z.prevSmoothPointX, 2*py - z.prevSmoothPointY
}

func (z *Renderer) initGradient(rgba color.RGBA) (ok bool) {
	cBase, nBase, shape, spread, nStops := ivg.DecodeGradient(rgba)

	prevN := negativeInfinity
	for i := uint8(0); i < nStops; i++ {
		c := z.cReg[(cBase+i)&0x3f]
		if !ivg.ValidAlphaPremulColor(c) {
			return false
		}
		n := z.nReg[(nBase+i)&0x3f]
		if !(0 <= n && n <= 1) || !(n > prevN) {
			return false
		}
		prevN = n
		z.stops[i] = Stop{
			Offset: float64(n),
			RGBA64: color.RGBA64{
				R: uint16(c.R) * 0x101,
				G: uint16(c.G) * 0x101,
				B: uint16(c.B) * 0x101,
				A: uint16(c.A) * 0x101,
			},
		}
	}

	// The affine transformation matrix in the IconVG graphic, stored in 6
	// contiguous NREG registers, goes from graphic coordinate space (i.e. the
	// viewBox) to the gradient coordinate space. We need it to start
	// in pixel space, not graphic coordinate space.

	invZSX := 1 / float64(z.scaleX)
	invZSY := 1 / float64(z.scaleY)
	zBX := float64(z.biasX)
	zBY := float64(z.biasY)

	a := float64(z.nReg[(nBase-6)&0x3f])
	b := float64(z.nReg[(nBase-5)&0x3f])
	c := float64(z.nReg[(nBase-4)&0x3f])
	d := float64(z.nReg[(nBase-3)&0x3f])
	e := float64(z.nReg[(nBase-2)&0x3f])
	f := float64(z.nReg[(nBase-1)&0x3f])

	pix2Grad := Aff3{
		a * invZSX,
		b * invZSY,
		c - a*zBX - b*zBY,
		d * invZSX,
		e * invZSY,
		f - d*zBX - e*zBY,
	}

	return z.gradient.Init(Shape(shape),Spread(spread),pix2Grad,z.stops[:nStops])
}

func (z *Renderer) StartPath(adj uint8, x, y float32) {
	z.flatColor = z.cReg[(z.cSel-adj)&0x3f]
	switch {
	case ivg.ValidAlphaPremulColor(z.flatColor):
		z.flatImage.C = &z.flatColor
		z.fill = &z.flatImage
		z.disabled = z.flatColor.A == 0
	case ivg.ValidGradient(z.flatColor):
		z.fill = &z.gradient
		z.disabled = !z.initGradient(z.flatColor)
	default:
		z.disabled = true
	}

	width, height := z.r.Dx(), z.r.Dy()
	h := float32(height)
	z.disabled = z.disabled || !(z.lod0 <= h && h < z.lod1)
	if z.disabled {
		return
	}

	z.z.Reset(width, height)
	z.prevSmoothType = smoothTypeNone
	z.z.MoveTo(z.absVec2(x, y))
}

func (z *Renderer) ClosePathEndPath() {
	if z.disabled {
		return
	}
	z.z.ClosePath()
	z.z.Draw(z.r, z.fill, image.Pt(0, 0))
}

func (z *Renderer) ClosePathAbsMoveTo(x, y float32) {
	if z.disabled {
		return
	}
	z.prevSmoothType = smoothTypeNone
	z.z.ClosePath()
	z.z.MoveTo(z.absVec2(x, y))
}

func (z *Renderer) ClosePathRelMoveTo(x, y float32) {
	if z.disabled {
		return
	}
	z.prevSmoothType = smoothTypeNone
	z.z.ClosePath()
	z.z.MoveTo(z.relVec2(x, y))
}

func (z *Renderer) AbsHLineTo(x float32) {
	if z.disabled {
		return
	}
	_, py := z.z.Pen()
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(z.absX(x), py)
}

func (z *Renderer) RelHLineTo(x float32) {
	if z.disabled {
		return
	}
	px, py := z.z.Pen()
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(px+z.relX(x), py)
}

func (z *Renderer) AbsVLineTo(y float32) {
	if z.disabled {
		return
	}
	px, _ := z.z.Pen()
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(px, z.absY(y))
}

func (z *Renderer) RelVLineTo(y float32) {
	if z.disabled {
		return
	}
	px, py := z.z.Pen()
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(px, py+z.relY(y))
}

func (z *Renderer) AbsLineTo(x, y float32) {
	if z.disabled {
		return
	}
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(z.absVec2(x, y))
}

func (z *Renderer) RelLineTo(x, y float32) {
	if z.disabled {
		return
	}
	z.prevSmoothType = smoothTypeNone
	z.z.LineTo(z.relVec2(x, y))
}

func (z *Renderer) AbsSmoothQuadTo(x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 := z.implicitSmoothPoint(smoothTypeQuad)
	x, y = z.absVec2(x, y)
	z.prevSmoothType = smoothTypeQuad
	z.prevSmoothPointX, z.prevSmoothPointY = x1, y1
	z.z.QuadTo(x1, y1, x, y)
}

func (z *Renderer) RelSmoothQuadTo(x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 := z.implicitSmoothPoint(smoothTypeQuad)
	x, y = z.relVec2(x, y)
	z.prevSmoothType = smoothTypeQuad
	z.prevSmoothPointX, z.prevSmoothPointY = x1, y1
	z.z.QuadTo(x1, y1, x, y)
}

func (z *Renderer) AbsQuadTo(x1, y1, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 = z.absVec2(x1, y1)
	x, y = z.absVec2(x, y)
	z.prevSmoothType = smoothTypeQuad
	z.prevSmoothPointX, z.prevSmoothPointY = x1, y1
	z.z.QuadTo(x1, y1, x, y)
}

func (z *Renderer) RelQuadTo(x1, y1, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 = z.relVec2(x1, y1)
	x, y = z.relVec2(x, y)
	z.prevSmoothType = smoothTypeQuad
	z.prevSmoothPointX, z.prevSmoothPointY = x1, y1
	z.z.QuadTo(x1, y1, x, y)
}

func (z *Renderer) AbsSmoothCubeTo(x2, y2, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 := z.implicitSmoothPoint(smoothTypeCube)
	x2, y2 = z.absVec2(x2, y2)
	x, y = z.absVec2(x, y)
	z.prevSmoothType = smoothTypeCube
	z.prevSmoothPointX, z.prevSmoothPointY = x2, y2
	z.z.CubeTo(x1, y1, x2, y2, x, y)
}

func (z *Renderer) RelSmoothCubeTo(x2, y2, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 := z.implicitSmoothPoint(smoothTypeCube)
	x2, y2 = z.relVec2(x2, y2)
	x, y = z.relVec2(x, y)
	z.prevSmoothType = smoothTypeCube
	z.prevSmoothPointX, z.prevSmoothPointY = x2, y2
	z.z.CubeTo(x1, y1, x2, y2, x, y)
}

func (z *Renderer) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 = z.absVec2(x1, y1)
	x2, y2 = z.absVec2(x2, y2)
	x, y = z.absVec2(x, y)
	z.prevSmoothType = smoothTypeCube
	z.prevSmoothPointX, z.prevSmoothPointY = x2, y2
	z.z.CubeTo(x1, y1, x2, y2, x, y)
}

func (z *Renderer) RelCubeTo(x1, y1, x2, y2, x, y float32) {
	if z.disabled {
		return
	}
	x1, y1 = z.relVec2(x1, y1)
	x2, y2 = z.relVec2(x2, y2)
	x, y = z.relVec2(x, y)
	z.prevSmoothType = smoothTypeCube
	z.prevSmoothPointX, z.prevSmoothPointY = x2, y2
	z.z.CubeTo(x1, y1, x2, y2, x, y)
}

func (z *Renderer) AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	if z.disabled {
		return
	}
	z.prevSmoothType = smoothTypeNone

	// We follow the "Conversion from endpoint to center parameterization"
	// algorithm as per
	// https://www.w3.org/TR/SVG/implnote.html#ArcConversionEndpointToCenter

	// There seems to be a bug in the spec's "implementation notes".
	//
	// Actual implementations, such as
	//	- https://git.gnome.org/browse/librsvg/tree/rsvg-path.c
	//	- http://svn.apache.org/repos/asf/xmlgraphics/batik/branches/svg11/sources/org/apache/batik/ext/awt/geom/ExtendedGeneralPath.java
	//	- https://java.net/projects/svgsalamander/sources/svn/content/trunk/svg-core/src/main/java/com/kitfox/svg/pathcmd/Arc.java
	//	- https://github.com/millermedeiros/SVGParser/blob/master/com/millermedeiros/geom/SVGArc.as
	// do something slightly different (marked with a †).

	// (†) The Abs isn't part of the spec. Neither is checking that Rx and Ry
	// are non-zero (and non-NaN).
	Rx := math.Abs(float64(rx))
	Ry := math.Abs(float64(ry))
	if !(Rx > 0 && Ry > 0) {
		z.z.LineTo(x, y)
		return
	}

	// We work in IconVG coordinates (e.g. from -32 to +32 by default), rather
	// than destination image coordinates (e.g. the width of the dst image),
	// since the rx and ry radii also need to be scaled, but their scaling
	// factors can be different, and aren't trivial to calculate due to
	// xAxisRotation.
	//
	// We convert back to destination image coordinates via absX and absY calls
	// later, during arcSegmentTo.
	penX, penY := z.z.Pen()
	x1 := float64(z.unabsX(penX))
	y1 := float64(z.unabsY(penY))
	x2 := float64(x)
	y2 := float64(y)

	phi := 2 * math.Pi * float64(xAxisRotation)

	// Step 1: Compute (x1′, y1′)
	halfDx := (x1 - x2) / 2
	halfDy := (y1 - y2) / 2
	cosPhi := math.Cos(phi)
	sinPhi := math.Sin(phi)
	x1Prime := +cosPhi*halfDx + sinPhi*halfDy
	y1Prime := -sinPhi*halfDx + cosPhi*halfDy

	// Step 2: Compute (cx′, cy′)
	rxSq := Rx * Rx
	rySq := Ry * Ry
	x1PrimeSq := x1Prime * x1Prime
	y1PrimeSq := y1Prime * y1Prime

	// (†) Check that the radii are large enough.
	radiiCheck := x1PrimeSq/rxSq + y1PrimeSq/rySq
	if radiiCheck > 1 {
		c := math.Sqrt(radiiCheck)
		Rx *= c
		Ry *= c
		rxSq = Rx * Rx
		rySq = Ry * Ry
	}

	denom := rxSq*y1PrimeSq + rySq*x1PrimeSq
	step2 := 0.0
	if a := rxSq*rySq/denom - 1; a > 0 {
		step2 = math.Sqrt(a)
	}
	if largeArc == sweep {
		step2 = -step2
	}
	cxPrime := +step2 * Rx * y1Prime / Ry
	cyPrime := -step2 * Ry * x1Prime / Rx

	// Step 3: Compute (cx, cy) from (cx′, cy′)
	cx := +cosPhi*cxPrime - sinPhi*cyPrime + (x1+x2)/2
	cy := +sinPhi*cxPrime + cosPhi*cyPrime + (y1+y2)/2

	// Step 4: Compute θ1 and Δθ
	ax := (+x1Prime - cxPrime) / Rx
	ay := (+y1Prime - cyPrime) / Ry
	bx := (-x1Prime - cxPrime) / Rx
	by := (-y1Prime - cyPrime) / Ry
	// angle returns the angle between the u and v vectors.
	angle := func(ux, uy, vx, vy float64) float64 {
		uNorm := math.Sqrt(ux*ux + uy*uy)
		vNorm := math.Sqrt(vx*vx + vy*vy)
		norm := uNorm * vNorm
		cos := (ux*vx + uy*vy) / norm
		ret := 0.0
		if cos <= -1 {
			ret = math.Pi
		} else if cos >= +1 {
			ret = 0
		} else {
			ret = math.Acos(cos)
		}
		if ux*vy < uy*vx {
			return -ret
		}
		return +ret
	}
	theta1 := angle(1, 0, ax, ay)
	deltaTheta := angle(ax, ay, bx, by)
	if sweep {
		if deltaTheta < 0 {
			deltaTheta += 2 * math.Pi
		}
	} else {
		if deltaTheta > 0 {
			deltaTheta -= 2 * math.Pi
		}
	}

	// This ends the
	// https://www.w3.org/TR/SVG/implnote.html#ArcConversionEndpointToCenter
	// algorithm. What follows below is specific to this implementation.

	// We approximate an arc by one or more cubic Bézier curves.
	n := int(math.Ceil(math.Abs(deltaTheta) / (math.Pi/2 + 0.001)))
	// arcSegmentTo approximates an arc by a cubic Bézier curve. The mathematical
	// formulae for the control points are the same as that used by librsvg.
	arcSegmentTo := func(cx, cy, theta1, theta2, rx, ry, cosPhi, sinPhi float64) {
		halfDeltaTheta := (theta2 - theta1) * 0.5
		q := math.Sin(halfDeltaTheta * 0.5)
		t := (8 * q * q) / (3 * math.Sin(halfDeltaTheta))
		cos1 := math.Cos(theta1)
		sin1 := math.Sin(theta1)
		cos2 := math.Cos(theta2)
		sin2 := math.Sin(theta2)
		x1 := rx * (+cos1 - t*sin1)
		y1 := ry * (+sin1 + t*cos1)
		x2 := rx * (+cos2 + t*sin2)
		y2 := ry * (+sin2 - t*cos2)
		x3 := rx * (+cos2)
		y3 := ry * (+sin2)
		z.z.CubeTo(
			z.absX(float32(cx+cosPhi*x1-sinPhi*y1)),
			z.absY(float32(cy+sinPhi*x1+cosPhi*y1)),
			z.absX(float32(cx+cosPhi*x2-sinPhi*y2)),
			z.absY(float32(cy+sinPhi*x2+cosPhi*y2)),
			z.absX(float32(cx+cosPhi*x3-sinPhi*y3)),
			z.absY(float32(cy+sinPhi*x3+cosPhi*y3)),
		)
	}
	for i := 0; i < n; i++ {
		arcSegmentTo(cx, cy,
			theta1+deltaTheta*float64(i+0)/float64(n),
			theta1+deltaTheta*float64(i+1)/float64(n),
			Rx, Ry, cosPhi, sinPhi,
		)
	}
}

func (z *Renderer) RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	ax, ay := z.relVec2(x, y)
	z.AbsArcTo(rx, ry, xAxisRotation, largeArc, sweep, z.unabsX(ax), z.unabsY(ay))
}
