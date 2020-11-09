// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivg

import (
	"fmt"
	"image/color"
)

// ColorType distinguishes types of Colors.
type ColorType uint8

const (
	// ColorTypeRGBA is a direct RGBA color.
	ColorTypeRGBA ColorType = iota

	// ColorTypePaletteIndex is an indirect color, indexing the custom palette.
	ColorTypePaletteIndex

	// ColorTypeCReg is an indirect color, indexing the CREG color registers.
	ColorTypeCReg

	// ColorTypeBlend is an indirect color, blending two other colors.
	ColorTypeBlend
)

// Color is an IconVG color, whose RGBA values can depend on context. Some
// Colors are direct RGBA values. Other Colors are indirect, referring to an
// index of the custom palette, a color register of the decoder virtual
// machine, or a blend of two other Colors.
//
// See the "Colors" section in the package documentation for details.
type Color struct {
	typ  ColorType
	data color.RGBA
}

// RGBAColor returns a direct Color.
func RGBAColor(c color.RGBA) Color { return Color{ColorTypeRGBA, c} }

// PaletteIndexColor returns an indirect Color referring to an index of the
// custom palette.
func PaletteIndexColor(i uint8) Color { return Color{ColorTypePaletteIndex, color.RGBA{R: i & 0x3f}} }

// CRegColor returns an indirect Color referring to a color register of the
// decoder virtual machine.
func CRegColor(i uint8) Color { return Color{ColorTypeCReg, color.RGBA{R: i & 0x3f}} }

// BlendColor returns an indirect Color that blends two other Colors. Those two
// other Colors must both be encodable as a 1 byte color.
//
// To blend a Color that is not encodable as a 1 byte color, first load that
// Color into a CREG color register, then call CRegColor to produce a Color
// that is encodable as a 1 byte color. See testdata/favicon.ivg for an
// example.
//
// See the "Colors" section in the package documentation for details.
func BlendColor(t, c0, c1 uint8) Color { return Color{ColorTypeBlend, color.RGBA{R: t, G: c0, B: c1}} }

func (c Color) rgba() color.RGBA         { return c.data }
func (c Color) paletteIndex() uint8      { return c.data.R }
func (c Color) cReg() uint8              { return c.data.R }
func (c Color) blend() (t, c0, c1 uint8) { return c.data.R, c.data.G, c.data.B }

// RGBA returns the color as a color.RGBA when that is its color type and the
// color is a valid premultiplied color. If the color is of a different color
// type or invalid, it will return  a opaque black and false.
func (c Color) RGBA() (color.RGBA, bool) {
	if c.typ != ColorTypeRGBA || !ValidAlphaPremulColor(c.data) {
		return color.RGBA{0x00, 0x00, 0x00, 0xff}, false
	}
	return c.data, true
}

// Resolve resolves the Color's RGBA value, given its context: the custom
// palette and the color registers of the decoder virtual machine.
func (c Color) Resolve(palette *[64]color.RGBA, cReg *[64]color.RGBA) color.RGBA {
	switch c.typ {
	case ColorTypeRGBA:
		return c.rgba()
	case ColorTypePaletteIndex:
		return palette[c.paletteIndex()&0x3f]
	case ColorTypeCReg:
		return cReg[c.cReg()&0x3f]
	case ColorTypeBlend:
		t, c0, c1 := c.blend()
		p, q := uint32(255-t), uint32(t)
		rgba0 := DecodeColor1(c0).Resolve(palette, cReg)
		rgba1 := DecodeColor1(c1).Resolve(palette, cReg)
		return color.RGBA{
			uint8(((p * uint32(rgba0.R)) + q*uint32(rgba1.R) + 128) / 255),
			uint8(((p * uint32(rgba0.G)) + q*uint32(rgba1.G) + 128) / 255),
			uint8(((p * uint32(rgba0.B)) + q*uint32(rgba1.B) + 128) / 255),
			uint8(((p * uint32(rgba0.A)) + q*uint32(rgba1.A) + 128) / 255),
		}
	}
	return color.RGBA{}
}

func DecodeColor1(x byte) Color {
	if x >= 0x80 {
		if x >= 0xc0 {
			return CRegColor(x)
		} else {
			return PaletteIndexColor(x)
		}
	}
	if x >= 125 {
		switch x - 125 {
		case 0:
			return RGBAColor(color.RGBA{0xc0, 0xc0, 0xc0, 0xc0})
		case 1:
			return RGBAColor(color.RGBA{0x80, 0x80, 0x80, 0x80})
		case 2:
			return RGBAColor(color.RGBA{0x00, 0x00, 0x00, 0x00})
		}
	}
	blue := dc1Table[x%5]
	x = x / 5
	green := dc1Table[x%5]
	x = x / 5
	red := dc1Table[x]
	return RGBAColor(color.RGBA{red, green, blue, 0xff})
}

var dc1Table = [5]byte{0x00, 0x40, 0x80, 0xc0, 0xff}

func Is1(c color.RGBA) bool {
	is1 := func(u uint8) bool { return u&0x3f == 0 || u == 0xff }
	return is1(c.R) && is1(c.G) && is1(c.B) && is1(c.A)
}

func Is2(c color.RGBA) bool {
	is2 := func(u uint8) bool { return u%0x11 == 0 }
	return is2(c.R) && is2(c.G) && is2(c.B) && is2(c.A)
}

func Is3(c color.RGBA) bool {
	return c.A == 0xff
}

func ValidAlphaPremulColor(c color.RGBA) bool {
	return c.R <= c.A && c.G <= c.A && c.B <= c.A
}

// ValidGradient returns true if the RGBA color is non-sensical
func ValidGradient(c color.RGBA) bool {
	return c.A == 0 && c.B&0x80 != 0
}

// EncodeGradient returns a non-sensical RGBA color encoding gradient
// parameters.
func EncodeGradient(cBase, nBase, shape, spread, nStops uint8) color.RGBA {
	cBase &= 0x3f
	nBase &= 0x3f
	shape = 0x02 | shape&0x01
	spread &= 0x03
	nStops &= 0x3f
	return color.RGBA{
		R: nStops,
		G: cBase | spread<<6,
		B: nBase | shape<<6,
		A: 0x00,
	}
}

// DecodeGradient returns the gradient parameters from a non-sensical RGBA
// color encoding a gradient.
func DecodeGradient(c color.RGBA) (cBase, nBase, shape, spread, nStops uint8) {
	cBase = c.G & 0x3f
	nBase = c.B & 0x3f
	shape = (c.B >> 6) & 0x01
	spread = (c.G >> 6) & 0x03
	nStops = c.R & 0x3f
	return
}

func (c Color) Is1() bool {
	return c.typ == ColorTypeRGBA && Is1(c.data)
}

func (c Color) Encode1() (x byte, ok bool) {
	switch c.typ {
	case ColorTypeRGBA:
		if c.data.A != 0xff {
			switch c.data {
			case color.RGBA{0x00, 0x00, 0x00, 0x00}:
				return 127, true
			case color.RGBA{0x80, 0x80, 0x80, 0x80}:
				return 126, true
			case color.RGBA{0xc0, 0xc0, 0xc0, 0xc0}:
				return 125, true
			}
		} else if Is1(c.data) {
			r := c.data.R / 0x3f
			g := c.data.G / 0x3f
			b := c.data.B / 0x3f
			return 25*r + 5*g + b, true
		}
	case ColorTypePaletteIndex:
		return c.data.R | 0x80, true
	case ColorTypeCReg:
		return c.data.R | 0xc0, true
	}
	return 0, false
}

func (c Color) Is2() bool {
	return c.typ == ColorTypeRGBA && Is2(c.data)
}

func (c Color) Encode2() (x [2]byte, ok bool) {
	if c.Is2() {
		return [2]byte{
			(c.data.R/0x11)<<4 | (c.data.G / 0x11),
			(c.data.B/0x11)<<4 | (c.data.A / 0x11),
		}, true
	}
	return [2]byte{}, false
}

func (c Color) Is3() bool {
	return c.typ == ColorTypeRGBA && Is3(c.data)
}

func (c Color) Encode3Direct() (x [3]byte, ok bool) {
	if c.Is3() {
		return [3]byte{c.data.R, c.data.G, c.data.B}, true
	}
	return [3]byte{}, false
}

func (c Color) Encode4() (x [4]byte, ok bool) {
	if c.typ == ColorTypeRGBA {
		return [4]byte{c.data.R, c.data.G, c.data.B, c.data.A}, true
	}
	return [4]byte{}, false
}

func (c Color) Encode3Indirect() (x [3]byte, ok bool) {
	if c.typ == ColorTypeBlend {
		return [3]byte{c.data.R, c.data.G, c.data.B}, true
	}
	return [3]byte{}, false
}

func (c Color) String() string {
	switch c.typ {
	case ColorTypeRGBA:
		rgba := c.rgba()
		switch {
		case ValidAlphaPremulColor(rgba):
			return fmt.Sprintf("RGBA %02x%02x%02x%02x", rgba.R, rgba.G, rgba.B, rgba.A)
		case ValidGradient(rgba):
			gradientShapeNames := [2]string{"linear", "radial"}
			gradientSpreadNames := [4]string{"none", "pad", "reflect", "repeat"}
			return fmt.Sprintf("gradient (NSTOPS=%d, CBASE=%d, NBASE=%d, %s, %s)",
				rgba.R&0x3f,
				rgba.G&0x3f,
				rgba.B&0x3f,
				gradientShapeNames[(rgba.B>>6)&0x01],
				gradientSpreadNames[rgba.G>>6],
			)
		}
	case ColorTypePaletteIndex:
		return fmt.Sprintf("customPalette[%d]", c.paletteIndex())
	case ColorTypeCReg:
		return fmt.Sprintf("CREG[%d]", c.cReg())
	case ColorTypeBlend:
		// old
		// 40                blend 191:64 c0:c1
		// ff                    c0: CREG[63]
		// 80                    c1: customPalette[0]

		// new
		// 40 ff 80          blend (191:64) (CREG[63]:customPalette[0])
		t, c0, c1 := c.blend()
		return fmt.Sprintf("blend (%d:%d) (%v:%v)", 0xff-t, t, DecodeColor1(c0), DecodeColor1(c1))
	}
	return fmt.Sprintf("nonsensical color")
}
