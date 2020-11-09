// SPDX-License-Identifier: Unlicense OR MIT

package vec

import (
	"image"
	"image/draw"

	"golang.org/x/image/vector"
)

// Rasterizer that wraps an inner "golang.org/x/image/vector" Rasterizer. The
// dst image normally passed to a call Draw is set as a field so Draw does not
// have to take it as a parameter.
type Rasterizer struct {
	vector.Rasterizer

	// Dst is the image that the Draw call uses as destination to draw into.
	Dst draw.Image

	// DrawOp is a Porter-Duff compositing operator that will be used for the
	// next call to the Draw method. After that call finishes, DrawOp is set to
	// draw.Over.
	DrawOp draw.Op
}

// NewRasterizer returns a rasterizer for dst image, with the dst size used to
// reset the inner rasterizer.
func NewRasterizer(dst draw.Image) *Rasterizer {
	r := &Rasterizer{Dst: dst}
	s := dst.Bounds().Size()
	r.Rasterizer.Reset(s.X, s.Y)
	return r
}

// Draw aligns r.Min in field Dst with sp in src and then replaces the
// rectangle r in Dst with the result of drawing src on Dst. The current value
// of the DrawOp field is used for drawing. But note, after drawing the DrawOp
// is reset to draw.Over.
func (z *Rasterizer) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	z.Rasterizer.DrawOp = z.DrawOp
	z.Rasterizer.Draw(z.Dst, r, src, sp)
	z.DrawOp = draw.Over
}
