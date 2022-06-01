package gio

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/raster/vec"
	"github.com/reactivego/ivg/render"
)

// Draw is a drawer based on "gioui.org/op/clip".
func Draw(ops *op.Ops, icon ivg.Icon, rect image.Rectangle, col ...color.Color) error {
	r := &render.Renderer{}
	z := &Rasterizer{Ops: ops}
	r.SetRasterizer(z, rect)
	return icon.RenderOn(r, col...)
}

// DrawGio is a shortcut to the default Draw func based on "gioui.org/op/clip".
var DrawGio = Draw

// DrawVec is a drawer based on "golang.org/x/image/vector".
func DrawVec(ops *op.Ops, icon ivg.Icon, rect image.Rectangle, col ...color.Color) error {
	r := &render.Renderer{}
	offset := rect.Min
	bounds := image.Rect(0, 0, rect.Dx(), rect.Dy())
	z := &vec.Rasterizer{Dst: image.NewRGBA(bounds), DrawOp: draw.Src}
	r.SetRasterizer(z, bounds)
	if err := icon.RenderOn(r, col...); err != nil {
		return err
	}
	paint.NewImageOp(z.Dst).Add(ops)
	tstack := op.Offset(f32.Pt(float32(offset.X), float32(offset.Y))).Push(ops)
	paint.PaintOp{}.Add(ops)
	tstack.Pop()
	return nil
}

type Options struct {
	Colors []color.Color
	Draw   Drawer
}

type Option func(*Options)

func WithColors(colors ...color.Color) Option {
	return func(options *Options) {
		options.Colors = colors
	}
}

type Drawer func(ops *op.Ops, icon ivg.Icon, rect image.Rectangle, col ...color.Color) error

func WithDrawer(drawer Drawer) Option {
	return func(options *Options) {
		options.Draw = drawer
	}
}

func Rasterize(icon ivg.Icon, rect image.Rectangle, options ...Option) (op.CallOp, error) {
	opts := Options{Draw: Draw}
	for _, option := range options {
		option(&opts)
	}
	ops := new(op.Ops)
	macro := op.Record(ops)
	err := opts.Draw(ops, icon, rect, opts.Colors...)
	return macro.Stop(), err
}

// IconCache is an icon cache that caches op.CallOp values returned by a call to
// the Rasterize method.
type IconCache struct {
	item map[key]op.CallOp
}

type key struct {
	checksum [md5.Size]byte
	rect     image.Rectangle
}

// NewIconCache returns a new icon cache.
func NewIconCache() *IconCache {
	return &IconCache{make(map[key]op.CallOp)}
}

// Rasterize returns a gio op.CallOp that paints the 'icon' inside the given
// rectangle 'rect' overiding colors with the colors 'col'.
func (c *IconCache) Rasterize(icon ivg.Icon, rect image.Rectangle, options ...Option) (op.CallOp, error) {
	data := []byte(icon.Name())
	opts := Options{Draw: Draw}
	for _, option := range options {
		option(&opts)
	}
	for _, col := range opts.Colors {
		c := color.NRGBAModel.Convert(col).(color.NRGBA)
		data = append(data, c.R, c.G, c.B, c.A)
	}
	data = append(data, fmt.Sprintf("%v", opts.Draw)...)
	key := key{md5.Sum(data), rect}
	if callOp, present := c.item[key]; present {
		return callOp, nil
	}
	ops := new(op.Ops)
	macro := op.Record(ops)
	err := opts.Draw(ops, icon, rect, opts.Colors...)
	callOp := macro.Stop()
	if err == nil {
		c.item[key] = callOp
	}
	return callOp, err
}

// Icon is an icon that implements the ivg.Icon interface.
type Icon struct {
	Data    []byte
	ViewBox ivg.ViewBox

	imgSize  int
	imgColor color.RGBA
	callOp   op.CallOp
}

// New creates a new IconVG (cachable) icon from the given data bytes.
func NewIcon(data []byte) (icon *Icon, err error) {
	i := &Icon{Data: data}
	if i.ViewBox, err = decode.DecodeViewBox(data); err != nil {
		return nil, err
	}
	return i, nil
}

func (i *Icon) Name() string {
	return string(i.Data)
}

func (i *Icon) RenderOn(dst ivg.Destination, col ...color.Color) error {
	opts := []decode.DecodeOption{}
	for idx, c := range col {
		opts = append(opts, decode.WithColorAt(idx, c))
	}
	return decode.Decode(dst, i.Data, opts...)
}

func (i *Icon) AspectMeet(size image.Point, ax, ay float32) image.Rectangle {
	minx, miny, maxx, maxy := i.ViewBox.AspectMeet(float32(size.X), float32(size.Y), ax, ay)
	return image.Rect(int(minx), int(miny), int(maxx), int(maxy))
}

func (i *Icon) AspectSlice(size image.Point, ax, ay float32) image.Rectangle {
	minx, miny, maxx, maxy := i.ViewBox.AspectSlice(float32(size.X), float32(size.Y), ax, ay)
	return image.Rect(int(minx), int(miny), int(maxx), int(maxy))
}

func (i *Icon) Layout(ops *op.Ops, sz int, c color.RGBA) image.Point {
	rect := i.AspectMeet(image.Pt(sz, sz), ivg.Mid, ivg.Mid)
	if sz != i.imgSize || c != i.imgColor {
		if callOp, err := Rasterize(i, rect, WithColors(c)); err != nil {
			return image.Pt(0, 0)
		} else {
			i.callOp = callOp
			i.imgSize = sz
			i.imgColor = c
		}
	}
	i.callOp.Add(ops)
	return image.Pt(sz, rect.Max.Y)
}
