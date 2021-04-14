// SPDX-License-Identifier: Unlicense OR MIT

package gio

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/reactivego/ivg/raster"
)

var (
	MinInf = float32(math.Inf(-1))
	MaxInf = float32(math.Inf(1))
)

type Rasterizer struct {
	Ops *op.Ops

	size       image.Point
	path       *clip.Path
	clipOp     clip.Op
	minX, minY float32
	maxX, maxY float32
}

func NewRasterizer(ops *op.Ops, w, h int) *Rasterizer {
	v := &Rasterizer{
		Ops:  ops,
		size: image.Pt(w, h),
		minX: MaxInf,
		minY: MaxInf,
		maxX: MinInf,
		maxY: MinInf,
	}
	return v
}

func (v *Rasterizer) Path() *clip.Path {
	if v.path == nil {
		v.path = new(clip.Path)
		if v.Ops == nil {
			v.Ops = new(op.Ops)
		}
		v.path.Begin(v.Ops)
	}
	return v.path
}

func (v *Rasterizer) Op() clip.Op {
	if v.path != nil {
		v.clipOp = clip.Outline{Path: v.path.End()}.Op()
		v.path = nil
	}
	return v.clipOp
}

func (v *Rasterizer) Reset(w, h int) {
	v.size = image.Pt(w, h)
	v.minX, v.minY = MaxInf, MaxInf
	v.maxX, v.maxY = MinInf, MinInf
	v.Op()
}

func (v *Rasterizer) Size() image.Point {
	return v.size
}

func (v *Rasterizer) Bounds() image.Rectangle {
	return image.Rectangle{Max: v.size}
}

func (v *Rasterizer) Pen() (x, y float32) {
	pos := v.path.Pos()
	return pos.X, pos.Y
}

func (v *Rasterizer) To(x, y float32) f32.Point {
	if x < v.minX {
		v.minX = float32(math.Floor(float64(x)))
	}
	if x > v.maxX {
		v.maxX = float32(math.Ceil(float64(x)))
	}
	if y < v.minY {
		v.minY = float32(math.Floor(float64(y)))
	}
	if y > v.maxY {
		v.maxY = float32(math.Ceil(float64(y)))
	}
	return f32.Pt(x, y)
}

func (v *Rasterizer) MoveTo(ax, ay float32) {
	v.Path().MoveTo(v.To(ax, ay))
}

func (v *Rasterizer) LineTo(bx, by float32) {
	v.Path().LineTo(v.To(bx, by))
}

func (v *Rasterizer) QuadTo(bx, by, cx, cy float32) {
	v.Path().QuadTo(f32.Pt(bx, by), v.To(cx, cy))
}

func (v *Rasterizer) CubeTo(bx, by, cx, cy, dx, dy float32) {
	v.Path().CubeTo(f32.Pt(bx, by), f32.Pt(cx, cy), v.To(dx, dy))
}

func (v *Rasterizer) ClosePath() {
	v.Path().Close()
}

func (v *Rasterizer) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	clip := v.Op()
	state := op.Save(v.Ops)
	op.Offset(f32.Pt(float32(r.Min.X), float32(r.Min.Y))).Add(v.Ops)
	clip.Add(v.Ops)
	switch source := src.(type) {
	case raster.GradientConfig:
		// TODO: If the gradient contains translucent colors we probably still must
		// convert the pixels using the RGBAModel from this package.
		gradient := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
		destrect := image.Rect(int(v.minX), int(v.minY), int(v.maxX), int(v.maxY))
		draw.Draw(gradient, destrect, src, destrect.Min.Add(sp), draw.Src)
		paint.NewImageOp(gradient).Add(v.Ops)
	case *image.Uniform:
		c := color.NRGBAModel.Convert(source.C).(color.NRGBA)
		paint.ColorOp{Color: c}.Add(v.Ops)
	default:
		paint.NewImageOp(src).Add(v.Ops)
	}
	paint.PaintOp{}.Add(v.Ops)
	state.Load()
}
