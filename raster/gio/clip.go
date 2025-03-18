package gio

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/render"
)

// Clip returns a widget that renders the given IconVG data using a clip.Path.
// According to the SVG specification, default value when the preserveAspectRatio attribute
// is not specified is "xMidYMid meet". This means that the image is scaled to fit the viewport
// while preserving the aspect ratio. The image is centered in the viewport along the x and y axes.
func Clip(data []byte, width, height unit.Dp, colors ...color.Color) (layout.Widget, error) {
	viewBox, err := decode.DecodeViewBox(data)
	if err != nil {
		return nil, err
	}
	lastSize := image.Point{}
	callOp := op.CallOp{}
	widget := func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Constrain(image.Pt(gtx.Dp(width), gtx.Dp(height)))
		minx, miny, maxx, maxy := viewBox.AspectMeet(float32(size.X), float32(size.Y), ivg.Mid, ivg.Mid)
		rect := image.Rect(int(minx), int(miny), int(maxx), int(maxy))
		if size != lastSize {
			lastSize = size
			ops := new(op.Ops)
			macro := op.Record(ops)
			// gio ->
			z := &Rasterizer{Ops: ops}
			r := &render.Renderer{}
			r.SetRasterizer(z, rect)
			// <- gio
			opts := []decode.DecodeOption{}
			for idx, c := range colors {
				opts = append(opts, decode.WithColorAt(idx, c))
			}
			decode.Decode(r, data, opts...)
			callOp = macro.Stop()
		}
		callOp.Add(gtx.Ops)
		return layout.Dimensions{Size: size}
	}
	return widget, nil
}
