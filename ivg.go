// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivg

import (
	"image/color"
)

const Magic = "\x89IVG"

var MagicBytes = []byte(Magic)

const (
	MidViewBox          = 0
	MidSuggestedPalette = 1
)

const (
	// Min aligns min of ViewBox with min of rect
	Min = 0.0
	// Mid aligns mid of ViewBox with mid of rect
	Mid = 0.5
	// Max aligns max of ViewBox with max of rect
	Max = 1.0
)

// ViewBox is a Rectangle
type ViewBox struct {
	MinX, MinY, MaxX, MaxY float32
}

// Size returns the ViewBox's size in both dimensions. An IconVG graphic is
// scalable; these dimensions do not necessarily map 1:1 to pixels.
func (v ViewBox) Size() (dx, dy float32) {
	return v.MaxX - v.MinX, v.MaxY - v.MinY
}

// AspectMeet fits the ViewBox inside a rectangle of size dx,dy maintaining its aspect ratio.
// The ax, ay argument determine the position of the resized viewbox in the
// given rectangle. For example ax = Mid, ay = Mid will position the resized
// viewbox always in the middle of the rectangle.
func (v ViewBox) AspectMeet(dx, dy float32, ax, ay float32) (MinX, MinY, MaxX, MaxY float32) {
	vdx, vdy := v.Size()
	vbAR := vdx / vdy
	vdx, vdy = dx, dy
	if vdx/vdy < vbAR {
		vdy = vdx / vbAR
	} else {
		vdx = vdy * vbAR
	}
	minX := (dx - vdx) * ax
	maxX := minX + vdx
	minY := (dy - vdy) * ay
	maxY := minY + vdy
	return minX, minY, maxX, maxY
}

// AspectSlice fills the rectangle of size dx,dy maintaining the ViewBox's aspect ratio.
// The ax,ay arguments determine the position of the resized viewbox in the given
// rectangle. For example ax = Mid, ay = Mid will position the resized viewbox
// always in the middle of the rectangle
func (v ViewBox) AspectSlice(dx, dy float32, ax, ay float32) (MinX, MinY, MaxX, MaxY float32) {
	vdx, vdy := v.Size()
	vbAR := vdx / vdy
	vdx, vdy = dx, dy
	if vdx/vdy < vbAR {
		vdx = vdy * vbAR
	} else {
		vdy = vdx / vbAR
	}
	minX := (dx - vdx) * ax
	maxX := minX + vdx
	minY := (dy - vdy) * ay
	maxY := minY + vdy
	return minX, minY, maxX, maxY
}

// Metadata is an IconVG's metadata.
type Metadata struct {
	ViewBox ViewBox

	// Palette is a 64 color palette. When encoding, it is the suggested
	// palette to place within the IconVG graphic. When decoding, it is either
	// the optional palette passed to Decode, or if no optional palette was
	// given, the suggested palette within the IconVG graphic.
	Palette [64]color.RGBA
}

// DefaultViewBox is the default ViewBox. Its values should not be modified.
var DefaultViewBox = ViewBox{
	MinX: -32, MinY: -32,
	MaxX: +32, MaxY: +32,
}

// DefaultPalette is the default Palette. Its values should not be modified.
var DefaultPalette = [64]color.RGBA{
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
	{0x00, 0x00, 0x00, 0xff},
}

// DefaultMetadata combines the default ViewBox and the default Palette.
var DefaultMetadata = Metadata{
	ViewBox: DefaultViewBox,
	Palette: DefaultPalette,
}

// Icon is an interface to an icon that can be drawn on a Destination
type Icon interface {
	// Name is a unique name of the icon inside your program. e.g. "favicon"
	// It is used to differentiate it from other icons in your program.
	Name() string

	// ViewBox is the ViewBox of the icon.
	ViewBox() ViewBox

	// RenderOn is called to let the icon render itself on
	// a Destination with a list of color.Color overrides.
	RenderOn(dst Destination, col ...color.Color) error
}
