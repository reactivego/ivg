// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decode

import (
	"image/color"
	"math"

	"github.com/reactivego/ivg"
)

// buffer holds an encoded IconVG graphic.
//
// The decodeXxx methods return the decoded value and an integer n, the number
// of bytes that value was encoded in. They return n == 0 if an error occured.
type buffer []byte

func (b buffer) decodeNatural() (u uint32, n int) {
	if len(b) < 1 {
		return 0, 0
	}
	x := b[0]
	if x&0x01 == 0 {
		return uint32(x) >> 1, 1
	}
	if x&0x02 == 0 {
		if len(b) >= 2 {
			y := uint16(b[0]) | uint16(b[1])<<8
			return uint32(y) >> 2, 2
		}
		return 0, 0
	}
	if len(b) >= 4 {
		y := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
		return y >> 2, 4
	}
	return 0, 0
}

func (b buffer) decodeReal() (f float32, n int) {
	switch u, n := b.decodeNatural(); n {
	case 0:
		return 0, n
	case 1:
		return float32(u), n
	case 2:
		return float32(u), n
	default:
		return math.Float32frombits(u << 2), n
	}
}

func (b buffer) decodeCoordinate() (f float32, n int) {
	switch u, n := b.decodeNatural(); n {
	case 0:
		return 0, n
	case 1:
		return float32(int32(u) - 64), n
	case 2:
		return float32(int32(u)-64*128) / 64, n
	default:
		return math.Float32frombits(u << 2), n
	}
}

func (b buffer) decodeZeroToOne() (f float32, n int) {
	switch u, n := b.decodeNatural(); n {
	case 0:
		return 0, n
	case 1:
		return float32(u) / 120, n
	case 2:
		return float32(u) / 15120, n
	default:
		return math.Float32frombits(u << 2), n
	}
}

func (b buffer) decodeColor1() (c ivg.Color, n int) {
	if len(b) < 1 {
		return ivg.Color{}, 0
	}
	return ivg.DecodeColor1(b[0]), 1
}

func (b buffer) decodeColor2() (c ivg.Color, n int) {
	if len(b) < 2 {
		return ivg.Color{}, 0
	}
	return ivg.RGBAColor(color.RGBA{
		R: 0x11 * (b[0] >> 4),
		G: 0x11 * (b[0] & 0x0f),
		B: 0x11 * (b[1] >> 4),
		A: 0x11 * (b[1] & 0x0f),
	}), 2
}

func (b buffer) decodeColor3Direct() (c ivg.Color, n int) {
	if len(b) < 3 {
		return ivg.Color{}, 0
	}
	return ivg.RGBAColor(color.RGBA{
		R: b[0],
		G: b[1],
		B: b[2],
		A: 0xff,
	}), 3
}

func (b buffer) decodeColor4() (c ivg.Color, n int) {
	if len(b) < 4 {
		return ivg.Color{}, 0
	}
	return ivg.RGBAColor(color.RGBA{
		R: b[0],
		G: b[1],
		B: b[2],
		A: b[3],
	}), 4
}

func (b buffer) decodeColor3Indirect() (c ivg.Color, n int) {
	if len(b) < 3 {
		return ivg.Color{}, 0
	}
	return ivg.BlendColor(b[0], b[1], b[2]), 3
}
