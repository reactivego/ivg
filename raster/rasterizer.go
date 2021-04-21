// SPDX-License-Identifier: Unlicense OR MIT

// Package raster provides rasterizers for 2-D vector graphics. Sub-directory
// gio provides an implementation based on gioui.org/op, while sub-directory
// vec provides an implementation based on golang.org/x/image/vector.
package raster

import (
	"image"
)

// Rasterizer is a 2-D vector graphics rasterizer.
type Rasterizer interface {
	// Reset resets a Rasterizer as if it was just returned by NewRasterizer.
	// This includes setting z.DrawOp to draw.Over.
	Reset(w, h int)
	// Size returns the width and height passed to NewRasterizer or Reset.
	Size() image.Point
	// Bounds returns the rectangle from (0, 0) to the width and height passed to
	// Reset.
	Bounds() image.Rectangle
	// Pen returns the location of the path-drawing pen: the last argument to the
	// most recent XxxTo call.
	Pen() (x, y float32)
	// MoveTo starts a new path and moves the pen to (ax, ay). The coordinates
	// are allowed to be out of the Rasterizer's bounds.
	MoveTo(ax, ay float32)
	// LineTo adds a line segment, from the pen to (bx, by), and moves the pen to
	// (bx, by). The coordinates are allowed to be out of the Rasterizer's
	// bounds.
	LineTo(bx, by float32)
	// QuadTo adds a quadratic Bézier segment, from the pen via (bx, by) to (cx,
	// cy), and moves the pen to (cx, cy). The coordinates are allowed to be out
	// of the Rasterizer's bounds.
	QuadTo(bx, by, cx, cy float32)
	// CubeTo adds a cubic Bézier segment, from the pen via (bx, by) and (cx, cy)
	// to (dx, dy), and moves the pen to (dx, dy). The coordinates are allowed to
	// be out of the Rasterizer's bounds.
	CubeTo(bx, by, cx, cy, dx, dy float32)
	// ClosePath closes the current path.
	ClosePath()
	// Draw aligns r.Min in z with sp in src and then replaces the rectangle r in
	// z with the result of a Porter-Duff composition. The vector paths
	// previously added via the XxxTo calls become the mask for drawing src onto
	// z.
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}
