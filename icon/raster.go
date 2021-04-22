package icon

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/reactivego/ivg/raster/gio"
	"github.com/reactivego/ivg/raster/vec"
	"github.com/reactivego/ivg/render"
)

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
	state := op.Save(ops)
	op.Offset(offset).Add(ops)
	paint.PaintOp{}.Add(ops)
	state.Load()
	return macro.Stop(), nil
}
