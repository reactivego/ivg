// SPDX-License-Identifier: Unlicense OR MIT

// Package icon implements a renderer for icons in the ivg format for gioui.org.
// The rasterizer used by the renderer can be switched between "gioui.org/op/clip" and
// "golang.org/x/image/vector".
package icon

import (
	"crypto/md5"
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/raster/gio"
	"github.com/reactivego/ivg/raster/vec"
	"github.com/reactivego/ivg/render"
)

type Icon interface {
	// RenderOn is called by a rasterizer to let the icon render itself on the
	// 'dst' ivg.Destination with the 'col' color overrides.
	RenderOn(dst ivg.Destination, col ...color.RGBA) error
}

type CachableIcon interface {
	Icon

	// Rasterize will rasterize the icon using a default or internal rasterizer.
	Rasterize(rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error)

	// Name is the unique name of the icon
	Name() string
}

// Rasterizer can render an icon and returns an op.CallOp. There are 2
// concrete rasterizers that implement this interface. VecRasterizer,
// GioRasterizer.
type Rasterizer interface {
	// Name is the unique name of the rasterizer.
	Name() string

	// Rasterize returns a gio op.CallOp that uses the rasterizer to render
	// an icon inside a given rect with some optional override colors.
	//
	// icon is an icon that conforms to the Icon interface.
	//
	// rect is the rectangle in (pixel coordinates) in which the icon should be
	// rendered. Use rect with Min at 0,0 for proper caching. Note that the icon
	// rendering is NOT clipped to the rect.
	//
	// col are the RGBA color overrides to use when rendering the icon.
	//
	// The function returns an op.CallOp and nil on success or an empty
	// op.CallOp and an error when the icon could not be rasterized.
	Rasterize(icon Icon, rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error)
}

// IconVG is an icon that implements the gioui.org/widget.IconVG interface.
type IconVG struct {
	ViewBox *ivg.ViewBox
	Palette *[64]color.RGBA

	data       []byte
	rasterizer Rasterizer

	callOp   op.CallOp
	imgSize  int
	imgColor color.RGBA
}

// New creates a new IconVG (cachable) icon from the given data bytes. The
// argument rasterizers can be used to (optionally) pass in a rasterizer. If
// no rastorizer is passed the GioRasterizer is used to directly render the
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

func (i *IconVG) RenderOn(dst ivg.Destination, col ...color.RGBA) error {
	for idx, c := range col {
		i.Palette[idx] = c
	}
	return decode.Decode(dst, i.data, &decode.DecodeOptions{Palette: i.Palette})
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

// GioRasterizer is a clipping rasterizer based on "gioui.org/op/clip".
var GioRasterizer gioRasterizer

type gioRasterizer struct{}

func (g gioRasterizer) Name() string {
	return "Gio"
}

func (g gioRasterizer) Rasterize(icon Icon, rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error) {
	ops := new(op.Ops)
	macro := op.Record(ops)
	r := &render.Renderer{}
	z := &gio.Rasterizer{Ops: ops}
	r.SetRasterizer(z, image.Rect(int(rect.Min.X), int(rect.Min.Y), int(rect.Max.X), int(rect.Max.Y)))
	if err := icon.RenderOn(r, col...); err != nil {
		return op.CallOp{}, err
	}
	return macro.Stop(), nil
}

// VecRasterizer is an image rasterizer based on "golang.org/x/image/vector".
var VecRasterizer vecRasterizer

type vecRasterizer struct{}

func (v vecRasterizer) Name() string {
	return "Vec"
}

func (v vecRasterizer) Rasterize(icon Icon, rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error) {
	ops := new(op.Ops)
	macro := op.Record(ops)
	r := &render.Renderer{}
	offset := rect.Min
	bounds := image.Rect(0, 0, int(rect.Dx()), int(rect.Dy()))
	z := &vec.Rasterizer{Dst: image.NewRGBA(bounds), DrawOp: draw.Src}
	r.SetRasterizer(z, bounds)
	if err := icon.RenderOn(r, col...); err != nil {
		return op.CallOp{}, err
	}
	paint.NewImageOp(z.Dst).Add(ops)
	stack := op.Push(ops)
	op.Offset(offset).Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Pop()
	return macro.Stop(), nil
}

// Cache is an icon cache that caches op.CallOp values returned by a call to
// the Rasterize method.
type Cache struct {
	item map[key]op.CallOp
}

type key struct {
	checksum [md5.Size]byte
	rect     f32.Rectangle
}

// NewCache returns a new icon cache.
func NewCache() *Cache {
	return &Cache{make(map[key]op.CallOp)}
}

// Rasterize returns a gio op.CallOp that paints the 'icon' inside the given
// rectangle 'rect' overiding colors with the colors 'col'.
func (c *Cache) Rasterize(icon CachableIcon, rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error) {
	data := []byte(icon.Name())
	for _, c := range col {
		data = append(data, c.R, c.G, c.B, c.A)
	}
	key := key{md5.Sum(data), rect}
	if callOp, present := c.item[key]; present {
		return callOp, nil
	}
	if callOp, err := icon.Rasterize(rect, col...); err == nil {
		c.item[key] = callOp
		return callOp, nil
	} else {
		return op.CallOp{}, err
	}
}
