// SPDX-License-Identifier: Unlicense OR MIT

package vector

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
	Dst    draw.Image
	DrawOp draw.Op
}

// NewRasterizer returns a rasterizer for dst image, with the dst width and
// height used to reset the inner rasterizer. The drawOp is set both as the
// inner drawOp and as the DrawOp field.
func NewRasterizer(dst draw.Image, drawOp draw.Op) *Rasterizer {
	r := &Rasterizer{Dst: dst, DrawOp: drawOp}
	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
	r.Rasterizer.Reset(w, h)
	r.Rasterizer.DrawOp = r.DrawOp
	return r
}

// Reset will reset the inner rasterizer to its initial state and then sets
// the inner DrawOp to the current DrawOp value. It will then set the current
// DrawOp value to draw.Over. So the next time Reset is called draw.Over will
// be set in the inner Rasterizer unless DrawOp is set before calling Reset.
func (r *Rasterizer) Reset(w, h int) {
	r.Rasterizer.Reset(w, h)
	r.Rasterizer.DrawOp = r.DrawOp
	r.DrawOp = draw.Over
}

func (r *Rasterizer) Draw(rect image.Rectangle, src image.Image, sp image.Point) {
	r.Rasterizer.Draw(r.Dst, rect, src, sp)
}
