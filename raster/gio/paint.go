package gio

import (
	"image"
	"image/draw"

	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/raster/img"
	"github.com/reactivego/ivg/render"
)

type PaintFunc func(*op.Ops, []byte, image.Rectangle, ...decode.DecodeOption)

func GioPaint(ops *op.Ops, data []byte, rect image.Rectangle, opts ...decode.DecodeOption) {
	z := &Rasterizer{Ops: ops}

	r := &render.Renderer{}
	r.SetRasterizer(z, rect)
	decode.Decode(r, data, opts...)
}

func ImagePaint(ops *op.Ops, data []byte, rect image.Rectangle, opts ...decode.DecodeOption) {
	offset, bounds := rect.Min, image.Rectangle{Max: rect.Size()}
	z := &img.Rasterizer{Dst: image.NewRGBA(bounds), DrawOp: draw.Src}

	r := &render.Renderer{}
	r.SetRasterizer(z, bounds)
	decode.Decode(r, data, opts...)

	paint.NewImageOp(z.Dst).Add(ops)
	defer op.Offset(offset).Push(ops).Pop()
	paint.PaintOp{}.Add(ops)
}
