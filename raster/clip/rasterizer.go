// SPDX-License-Identifier: Unlicense OR MIT

package clip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Rasterizer struct {
	ops   *op.Ops
	size  image.Point
	first struct{ x, y float32 }
	pen   struct{ x, y float32 }
	path  *clip.Path
	clipOp  clip.Op
}

func NewRasterizer(w, h int, ops *op.Ops) *Rasterizer {
	v := &Rasterizer{ops: ops, size: image.Pt(w, h)}
	return v
}

func (v *Rasterizer) Path() *clip.Path {
	if v.path == nil {
		v.path = new(clip.Path)
		if v.ops == nil {
			v.ops = new(op.Ops)
		}
		v.path.Begin(v.ops)
	}
	return v.path
}

func (v *Rasterizer) Op() clip.Op {
	if v.path != nil {
		v.clipOp = v.path.End()
		v.path = nil
	}
	return v.clipOp
}

func (v *Rasterizer) Reset(w, h int) {
	v.size = image.Pt(w, h)
	v.first.x, v.first.y = 0, 0
	v.pen.x, v.pen.y = 0, 0
	v.Op()
}

func (v *Rasterizer) Size() image.Point {
	return v.size
}

func (v *Rasterizer) Bounds() image.Rectangle {
	return image.Rectangle{Max: v.size}
}

func (v *Rasterizer) Pen() (x, y float32) {
	return v.pen.x, v.pen.y
}

func (v *Rasterizer) MoveTo(ax, ay float32) {
	v.Path().Move(f32.Pt(ax-v.pen.x, ay-v.pen.y))
	v.first.x, v.first.y = ax, ay
	v.pen.x, v.pen.y = ax, ay
}

func (v *Rasterizer) LineTo(bx, by float32) {
	v.Path().Line(f32.Pt(bx-v.pen.x, by-v.pen.y))
	v.pen.x, v.pen.y = bx, by
}

func (v *Rasterizer) QuadTo(bx, by, cx, cy float32) {
	v.Path().Quad(f32.Pt(bx-v.pen.x, by-v.pen.y), f32.Pt(cx-v.pen.x, cy-v.pen.y))
	v.pen.x, v.pen.y = cx, cy
}

func (v *Rasterizer) CubeTo(bx, by, cx, cy, dx, dy float32) {
	v.Path().Cube(f32.Pt(bx-v.pen.x, by-v.pen.y), f32.Pt(cx-v.pen.x, cy-v.pen.y), f32.Pt(dx-v.pen.x, dy-v.pen.y))
	v.pen.x, v.pen.y = dx, dy
}

func (v *Rasterizer) ClosePath() {
	v.LineTo(v.first.x, v.first.y)
}

func (v *Rasterizer) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	clip := v.Op()
	stack := op.Push(v.ops)
	op.Offset(f32.Pt(float32(r.Min.X), float32(r.Min.Y))).Add(v.ops)
	clip.Add(v.ops)
	paint.NewImageOp(src).Add(v.ops)
	rect := f32.Rect(0, 0, float32(r.Max.X-r.Min.X), float32(r.Max.Y-r.Min.Y))
	paint.PaintOp{Rect: rect}.Add(v.ops)
	stack.Pop()
}
