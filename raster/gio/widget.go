package gio

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
)

type Option = func(*option)

type option struct {
	Paint   PaintFunc
	Options []decode.DecodeOption
}

func WithImageBackend() Option {
	return func(o *option) {
		o.Paint = ImagePaint
	}
}

func WithColors(colors ...color.Color) Option {
	return func(o *option) {
		for idx, c := range colors {
			o.Options = append(o.Options, decode.WithColorAt(idx, c))
		}
	}
}

// Widget creates a layout widget for rendering IconVG vector graphics data. It supports two rendering
// backends: a default Gio clip.Path implementation and an optional image-based raster backend
// (enabled via WithImageBackend()). The widget handles aspect ratio preservation following the
// SVG specification's "xMidYMid meet" behavior, which scales the image to fit the viewport while
// maintaining proportions and centering it both horizontally and vertically.
//
// The data parameter accepts the raw IconVG bytes, while width and height specify the desired
// dimensions in device-independent pixels (Dp). Additional rendering options can be provided
// through the variadic options parameter.
func Widget(data []byte, width, height unit.Dp, options ...Option) (layout.Widget, error) {
	viewBox, err := decode.DecodeViewBox(data)
	if err != nil {
		return nil, err
	}
	o := &option{Paint: GioPaint}
	for _, f := range options {
		f(o)
	}
	lastSize := image.Point{}
	callOp := op.CallOp{}
	widget := func(gtx layout.Context) layout.Dimensions {
		newSize := gtx.Constraints.Constrain(image.Pt(gtx.Dp(width), gtx.Dp(height)))
		minx, miny, maxx, maxy := viewBox.AspectMeet(float32(newSize.X), float32(newSize.Y), ivg.Mid, ivg.Mid)
		rect := image.Rect(int(minx), int(miny), int(maxx), int(maxy))
		if newSize != lastSize {
			lastSize = newSize
			ops := new(op.Ops)
			macro := op.Record(ops)
			o.Paint(ops, data, rect, o.Options...)
			callOp = macro.Stop()
		}
		callOp.Add(gtx.Ops)
		return layout.Dimensions{Size: newSize}
	}
	return widget, nil
}
