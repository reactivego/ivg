package raster

import (
	"image"
)

type Rasterizer interface {
	Reset(w, h int)
	Size() image.Point
	Bounds() image.Rectangle
	Pen() (x, y float32)
	MoveTo(ax, ay float32)
	LineTo(bx, by float32)
	QuadTo(bx, by, cx, cy float32)
	CubeTo(bx, by, cx, cy, dx, dy float32)
	ClosePath()
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}
