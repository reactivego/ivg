package ivg

import "image/color"

// Destination handles the actions decoded from an IconVG graphic's opcodes.
//
// When passed to Decode, the first method called (if any) will be Reset. No
// methods will be called at all if an error is encountered in the encoded form
// before the metadata is fully decoded.
type Destination interface {
	Reset(viewbox ViewBox, palette [64]color.RGBA)
	CSel() uint8
	SetCSel(cSel uint8)
	NSel() uint8
	SetNSel(nSel uint8)
	SetCReg(adj uint8, incr bool, c Color)
	SetNReg(adj uint8, incr bool, f float32)
	SetLOD(lod0, lod1 float32)

	StartPath(adj uint8, x, y float32)
	ClosePathEndPath()
	ClosePathAbsMoveTo(x, y float32)
	ClosePathRelMoveTo(x, y float32)

	AbsHLineTo(x float32)
	RelHLineTo(x float32)
	AbsVLineTo(y float32)
	RelVLineTo(y float32)
	AbsLineTo(x, y float32)
	RelLineTo(x, y float32)
	AbsSmoothQuadTo(x, y float32)
	RelSmoothQuadTo(x, y float32)
	AbsQuadTo(x1, y1, x, y float32)
	RelQuadTo(x1, y1, x, y float32)
	AbsSmoothCubeTo(x2, y2, x, y float32)
	RelSmoothCubeTo(x2, y2, x, y float32)
	AbsCubeTo(x1, y1, x2, y2, x, y float32)
	RelCubeTo(x1, y1, x2, y2, x, y float32)
	AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
	RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
}
