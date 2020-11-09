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

// Rectangle is defined by its minimum and maximum coordinates.
type Rectangle struct {
	MinX, MinY, MaxX, MaxY float32
}

// Rect returns a rectangle for the given minX, minY, maxX and maxY float32
// arguments.
func Rect(minX, minY, maxX, maxY float32) Rectangle {
	return Rectangle{minX, minY, maxX, maxY}
}

// Fields returns the individual fields of the rectangle as multiple
// float32 return arguments.
func (r Rectangle) Fields() (MinX, MinY, MaxX, MaxY float32) {
	return r.MinX, r.MinY, r.MaxX, r.MaxY
}

// Size returns the Rectangles's size in both dimensions.
func (r Rectangle) Size() (dx, dy float32) {
	return r.MaxX - r.MinX, r.MaxY - r.MinY
}

// IntFields returns the individual fields of the rectangle as multiple
// int return arguments by truncating the float32 to int.
func (r Rectangle) IntFields() (MinX, MinY, MaxX, MaxY int) {
	return int(r.MinX), int(r.MinY), int(r.MaxX), int(r.MaxY)
}

// PreserveAspectRatio is the SVG attribute 'PreserveAspectRatio' which
// determines how the ViewBox is sized w.r.t. a bounding rectangle.
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

// ViewBox is a Rectangle
type ViewBox Rectangle

// Size returns the ViewBox's size in both dimensions. An IconVG graphic is
// scalable; these dimensions do not necessarily map 1:1 to pixels.
func (v *ViewBox) Size() (dx, dy float32) {
	return v.MaxX - v.MinX, v.MaxY - v.MinY
}

// SizeToRect resizes and positions the viewbox in the given rect. The aspect
// argument determines how the ViewBox is positioned in the rect. The ax, ay
// argument determine the position of the resized viewbox in the given rect.
// For example ax = Mid, ay = Mid will position the resized viewbox always in
// the middle of the rect
func (v *ViewBox) SizeToRect(rect Rectangle, aspect PreserveAspectRatio, ax, ay float32) Rectangle {
	rdx, rdy := rect.Size()
	vdx, vdy := v.Size()
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
	rect.MinX += (rdx - vdx) * ax
	rect.MaxX = rect.MinX + vdx
	rect.MinY += (rdy - vdy) * ay
	rect.MaxY = rect.MinY + vdy
	return rect
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
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
}
