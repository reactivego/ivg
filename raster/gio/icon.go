package gio

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type Option = func(*option)

type option struct {
	Icon   func(data []byte, width, height unit.Dp, colors ...color.Color) (layout.Widget, error)
	Colors []color.Color
}

func WithClipRasterizer() Option {
	return func(o *option) {
		o.Icon = Clip
	}
}

func WithVecRasterizer() Option {
	return func(o *option) {
		o.Icon = Vec
	}
}

func WithColors(colors ...color.Color) Option {
	return func(o *option) {
		o.Colors = colors
	}
}

func Icon(data []byte, width, height unit.Dp, options ...Option) (layout.Widget, error) {
	o := &option{Icon: Clip}
	for _, f := range options {
		f(o)
	}
	return o.Icon(data, width, height, o.Colors...)
}
