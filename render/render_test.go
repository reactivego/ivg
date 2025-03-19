package render

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/raster/img"
)

func TestInvalidAlphaPremultipliedColor(t *testing.T) {
	// See http://golang.org/issue/39526 for some discussion.

	dst := image.NewRGBA(image.Rect(0, 0, 1, 1))
	rasterizer := img.NewRasterizer(dst)
	var z Renderer
	z.SetRasterizer(rasterizer, dst.Bounds())
	z.Reset(ivg.ViewBox{MinX: 0.0, MinY: 0.0, MaxX: 1.0, MaxY: 1.0}, ivg.DefaultPalette)

	// Fill the unit square with red.
	z.SetCReg(0, false, ivg.RGBAColor(color.RGBA{0x55, 0x00, 0x00, 0x66}))
	z.StartPath(0, 0.0, 0.0)
	z.AbsLineTo(1.0, 0.0)
	z.AbsLineTo(1.0, 1.0)
	z.AbsLineTo(0.0, 1.0)
	z.ClosePathEndPath()

	// Fill the unit square with an invalid (non-gradient) alpha-premultiplied
	// color (super-saturated green). This should be a no-op (and not crash).
	z.SetCReg(0, false, ivg.RGBAColor(color.RGBA{0x00, 0x99, 0x00, 0x88}))
	z.StartPath(0, 0.0, 0.0)
	z.AbsLineTo(1.0, 0.0)
	z.AbsLineTo(1.0, 1.0)
	z.AbsLineTo(0.0, 1.0)
	z.ClosePathEndPath()

	// We should see red.
	got := dst.Pix
	want := []byte{0x55, 0x00, 0x00, 0x66}
	if !bytes.Equal(got, want) {
		t.Errorf("got [% 02x], want [% 02x]", got, want)
	}
}
