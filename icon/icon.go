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
	"github.com/reactivego/ivg/raster/vector"
	"github.com/reactivego/ivg/render"
)

// PreserveAspectRatio is the SVG attribute 'PreserveAspectRatio' which
// determines how the ViewBox of an icon is scaled w.r.t. a bounding
// rectangle.
type PreserveAspectRatio int

const (
	// AspectNone stretches or squashes the ViewBox to meet the rect.
	AspectNone PreserveAspectRatio = iota
	// AspectMeet fits the ViewBox inside the rect maintaining its aspect ratio.
	AspectMeet
	// AspectSlice fills the rect maintaining the ViewBox's aspect ratio.
	AspectSlice
)

const (
	// Min aligns min of ViewBox with min of rect
	Min = 0.0
	// Mid aligns mid of ViewBox with mid of rect
	Mid = 0.5
	// Max aligns max of ViewBox with max of rect
	Max = 1.0
)

// Rasterizer specifies the rasterizer to use for rendering the icon.
type Rasterizer int

const (
	// GioRasterizer selects "gioui.org/op/clip" as rasterizer
	GioRasterizer Rasterizer = iota
	// VectorRasterizer selects "golang.org/x/image/vector" as rasterizer
	VectorRasterizer
)

type key struct {
	md5    [16]byte
	col      color.RGBA
	rect   image.Rectangle
	aspect PreserveAspectRatio
	ax     float32
	ay     float32
}

type Cache struct {
	item   map[key]op.CallOp
	raster Rasterizer
}

func NewCache(raster Rasterizer) *Cache {
	return &Cache{item: make(map[key]op.CallOp), raster: raster}
}

func (c *Cache) FromData(data []byte, col color.RGBA, rect image.Rectangle, aspect PreserveAspectRatio, ax, ay float32) (op.CallOp, error) {
	key := key{md5.Sum(data), col, rect, aspect, ax, ay}
	if callOp, present := c.item[key]; present {
		return callOp, nil
	}
	callOp, err := FromData(data, col, rect, aspect, ax, ay, c.raster)
	c.item[key] = callOp
	return callOp, err
}

// FromData returns a gio op.CallOp that paints the icon decoded from 'data'
// with the given color 'c' inside the given rectangle 'rect'.
//
// data is the ivg encoded data representation of the icon.
//
// c is the color.RGBA to render the icon in.
//
// rect is the rectangle in (pixel coordinates) in which the icon should be
// rendered. Use rect with Min at 0,0 for proper caching. Note that the icon
// rendering is NOT clipped to the rect.
//
// aspect is the SVG attribute 'PreserveAspectRatio' which determines how the
// ViewBox of an icon is scaled w.r.t. the bounding 'rect'. Valid values are
// AspectNone, AspectMeet or AspectSlice.
//
// ax is a value from 0.0 to 1.0 which determines how the ViewBox is
// positioned horizontally in the rect. Min (0.0) aligns the left side of both
// rectangles. Mid (0.5) aligns the centers of both rectangles. Max (1.0)
// aligns the right side of both rectangles.
//
// ay is a value from 0.0 to 1.0 which determines how the ViewBox is
// positioned vertically in the rect. Min (0.0) aligns the top of both
// rectangles. Mid (0.5) aligns the middle of both rectangles. Max (1.0)
// aligns the bottom of both rectangles
//
// raster specifies the rasterizer to use for rendering the icon.
// GioRasterizer selects "gioui.org/op/clip" as rasterizer and
// VectorRasterizer selects "golang.org/x/image/vector" as rasterizer
//
// The function returns an op.CallOp and nil on success or an empty
// op.CallOp and an error when the icon could not be renderdered.
func FromData(data []byte, c color.RGBA, rect image.Rectangle, aspect PreserveAspectRatio, ax, ay float32, raster Rasterizer) (op.CallOp, error) {
	var callOp op.CallOp
	viewbox := ivg.DefaultViewBox
	palette := &ivg.DefaultPalette
	if md, err := decode.DecodeMetadata(data); err == nil {
		viewbox = md.ViewBox
		palette = &md.Palette
	} else {
		return callOp, err
	}
	(*palette)[0] = c
	options := &decode.DecodeOptions{Palette: palette}
	rdx, rdy := float32(rect.Dx()), float32(rect.Dy())
	vdx, vdy := viewbox.AspectRatio()
	vbAR := vdx / vdy
	vdx, vdy = rdx, rdy
	switch aspect {
	case AspectMeet:
		if vdx/vdy < vbAR {
			vdy = vdx / vbAR
		} else {
			vdx = vdy * vbAR
		}
	case AspectSlice:
		if vdx/vdy < vbAR {
			vdx = vdy * vbAR
		} else {
			vdy = vdx / vbAR
		}
	}
	dx := (rdx - vdx)
	dy := (rdy - vdy)
	minX := float32(rect.Min.X) + dx*ax
	maxX := minX + vdx
	minY := float32(rect.Min.X) + dy*ay
	maxY := minY + vdy
	rect32 := f32.Rect(minX, minY, maxX, maxY)
	ops := new(op.Ops)
	macro := op.Record(ops)
	var z render.Renderer
	switch raster {
	case GioRasterizer:
		r := gio.NewRasterizer(int(rect32.Dx()), int(rect32.Dy()), ops)
		rect = image.Rect(int(rect32.Min.X), int(rect32.Min.Y), int(rect32.Max.X), int(rect32.Max.Y))
		z.SetRasterizer(r, rect)
		if err := decode.Decode(&z, data, options); err != nil {
			return callOp, err
		}
	case VectorRasterizer:
		r := vector.NewRasterizer(image.NewRGBA(image.Rect(0, 0, int(rect32.Dx()), int(rect32.Dy()))), draw.Src)
		z.SetRasterizer(r, r.Bounds())
		if err := decode.Decode(&z, data, options); err != nil {
			return callOp, err
		}
		paint.NewImageOp(r.Dst).Add(ops)
		// paint.ColorOp{Color: c}.Add(ops)
		paint.PaintOp{Rect: rect32}.Add(ops)
	}
	callOp = macro.Stop()
	return callOp, nil
}
