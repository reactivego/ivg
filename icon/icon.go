// SPDX-License-Identifier: Unlicense OR MIT

// Package icon implements a renderer for icons in the ivg format for gioui.org.
// The rasterizer used by the renderer can be switched between "gioui.org/op/clip" and
// "golang.org/x/image/vector".
package icon

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
)

type Icon interface {
	// RenderOn is called by a rasterizer to let the icon render itself on the
	// 'dst' ivg.Destination with the 'col' color overrides.
	RenderOn(dst ivg.Destination, col ...color.RGBA) error
}

// IconVG is an icon that implements the gioui.org/widget.IconVG interface.
type IconVG struct {
	ViewBox *ivg.ViewBox
	Palette *[64]color.RGBA
	data    []byte

	rasterizer Rasterizer
	imgSize    int
	imgColor   color.RGBA
	callOp     op.CallOp
}

// New creates a new IconVG (cachable) icon from the given data bytes. The
// argument rasterizers can be used to (optionally) pass in a rasterizer. If
// no rasterizer is passed the GioRasterizer is used to directly render the
// icon using gio clip operations.
func New(data []byte, rasterizers ...Rasterizer) (icon *IconVG, err error) {
	i := &IconVG{data: data}
	if metadata, err := decode.DecodeMetadata(data); err != nil {
		return nil, err
	} else {
		i.ViewBox = &metadata.ViewBox
		i.Palette = &metadata.Palette
	}
	if len(rasterizers) > 0 {
		i.rasterizer = rasterizers[0]
	}
	return i, nil
}

func (i *IconVG) RenderOn(dst ivg.Destination, col ...color.RGBA) error {
	for idx, c := range col {
		i.Palette[idx] = c
	}
	return decode.Decode(dst, i.data, &decode.DecodeOptions{Palette: i.Palette})
}

func (i *IconVG) Name() string {
	if i.rasterizer != nil {
		return string(i.data) + i.rasterizer.Name()
	} else {
		return string(i.data) + GioRasterizer.Name()
	}
}

func (i *IconVG) Rasterize(rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error) {
	if i.rasterizer != nil {
		return i.rasterizer.Rasterize(i, rect, col...)
	}
	return GioRasterizer.Rasterize(i, rect, col...)
}

func (i *IconVG) AspectMeet(rect f32.Rectangle, ax, ay float32) f32.Rectangle {
	return f32.Rect(i.ViewBox.AspectMeet(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, ax, ay))
}

func (i *IconVG) AspectSlice(rect f32.Rectangle, ax, ay float32) f32.Rectangle {
	return f32.Rect(i.ViewBox.AspectSlice(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, ax, ay))
}

func (i *IconVG) Layout(ops *op.Ops, sz int, c color.RGBA) image.Point {
	rect := i.AspectMeet(f32.Rect(0, 0, float32(sz), float32(sz)), ivg.Mid, ivg.Mid)
	if sz != i.imgSize || c != i.imgColor {
		if callOp, err := i.Rasterize(rect, c); err != nil {
			return image.Pt(0, 0)
		} else {
			i.callOp = callOp
			i.imgSize = sz
			i.imgColor = c
		}
	}
	i.callOp.Add(ops)
	return image.Pt(sz, int(rect.Max.Y))
}
