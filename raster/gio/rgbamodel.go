// SPDX-License-Identifier: Unlicense OR MIT

package gio

import (
	"image/color"
	"math"
)

// RGBAModel is a color.Model that can convert any color.Color to a color.RGBA
// that can be passed to Gio.
//
//	yellow := color.NRGBA{0xfd, 0xee, 0x74, 0x7f}
//	rgba := gio.RGBAModel.Convert(yellow).(color.RGBA)
var RGBAModel color.Model = color.ModelFunc(rGBAModel)

func rGBAModel(c color.Color) color.Color {
	var r, g, b, a uint32
	premul := float32(0xffff)
	switch c := c.(type) {
	case color.NRGBA64:
		r, g, b, a = uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
	case color.NRGBA:
		r = uint32(c.R) | uint32(c.R)<<8
		g = uint32(c.G) | uint32(c.G)<<8
		b = uint32(c.B) | uint32(c.B)<<8
		a = uint32(c.A) | uint32(c.A)<<8
	default:
		r, g, b, a = c.RGBA()
		premul = float32(a)
	}
	if a == 0 {
		return color.RGBA{0, 0, 0, 0}
	}
	if a == 0xffff {
		return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xff}
	}
	rf := sRGBToLinear(float32(r) / premul)
	rf = linearTosRGB(rf * float32(a) / 0xffff)
	gf := sRGBToLinear(float32(g) / premul)
	gf = linearTosRGB(gf * float32(a) / 0xffff)
	bf := sRGBToLinear(float32(b) / premul)
	bf = linearTosRGB(bf * float32(a) / 0xffff)
	return color.RGBA{uint8(rf*255 + .5), uint8(gf*255 + .5), uint8(bf*255 + .5), uint8(a >> 8)}
}

// linearTosRGB transforms color value from linear to sRGB.
func linearTosRGB(c float32) float32 {
	// Formula from EXT_sRGB.
	switch {
	case c <= 0:
		return 0
	case 0 < c && c < 0.0031308:
		return 12.92 * c
	case 0.0031308 <= c && c < 1:
		return 1.055*float32(math.Pow(float64(c), 0.41666)) - 0.055
	}

	return 1
}

// sRGBToLinear transforms color value from sRGB to linear.
func sRGBToLinear(c float32) float32 {
	// Formula from EXT_sRGB.
	if c <= 0.04045 {
		return c / 12.92
	} else {
		return float32(math.Pow(float64((c+0.055)/1.055), 2.4))
	}
}
