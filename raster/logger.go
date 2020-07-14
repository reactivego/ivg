// SPDX-License-Identifier: Unlicense OR MIT

package raster

import (
	"fmt"
	"image"
)

type RasterizerLogger struct {
	Rasterizer
}

func (r *RasterizerLogger) Reset(w, h int) {
	fmt.Printf("raster.Reset(w:%d, h:%d)\n", w, h)
	r.Rasterizer.Reset(w, h)
}

func (r *RasterizerLogger) Size() image.Point {
	s := r.Rasterizer.Size()
	// fmt.Printf("raster.Size() = %#v\n", s)
	return s
}

func (r *RasterizerLogger) Bounds() image.Rectangle {
	b := r.Rasterizer.Bounds()
	// fmt.Printf("raster.Bounds() = %#v\n", b)
	return b
}

func (r *RasterizerLogger) Pen() (x, y float32) {
	x, y = r.Rasterizer.Pen()
	// fmt.Printf("raster.Pen() = (x:%.2f, y:%.2f)\n", x, y)
	return x, y
}

func (r *RasterizerLogger) MoveTo(ax, ay float32) {
	fmt.Printf("raster.MoveTo(ax:%.2f, ay:%.2f)\n", ax, ay)
	r.Rasterizer.MoveTo(ax, ay)
}

func (r *RasterizerLogger) LineTo(bx, by float32) {
	fmt.Printf("raster.LineTo(bx:%.2f, by:%.2f)\n", bx, by)
	r.Rasterizer.LineTo(bx, by)
}

func (r *RasterizerLogger) QuadTo(bx, by, cx, cy float32) {
	fmt.Printf("raster.QuadTo(bx:%.2f, by:%.2f, cx:%.2f, cy:%.2f)\n", bx, by, cx, cy)
	r.Rasterizer.QuadTo(bx, by, cx, cy)
}

func (r *RasterizerLogger) CubeTo(bx, by, cx, cy, dx, dy float32) {
	fmt.Printf("raster.CubeTo(bx:%.2f, by:%.2f, cx:%.2f, cy:%.2f, dx:%.2f, dy:%.2f)\n", bx, by, cx, cy, dx, dy)
	r.Rasterizer.CubeTo(bx, by, cx, cy, dx, dy)
}

func (r *RasterizerLogger) ClosePath() {
	fmt.Printf("raster.ClosePath()\n")
	r.Rasterizer.ClosePath()
}

func (rl *RasterizerLogger) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	fmt.Printf("raster.Draw(r: %#v, src: %#v, sp: %#v)\n", r, src.Bounds(), sp)
	rl.Rasterizer.Draw(r, src, sp)
}
