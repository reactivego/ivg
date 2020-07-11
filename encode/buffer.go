// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encode

import (
	"math"

	"github.com/reactivego/ivg"
)

// buffer holds an encoded IconVG graphic.
//
// The encodeXxx methods append to the buffer, modifying the slice in place.
type buffer []byte

func (b *buffer) encodeNatural(u uint32) {
	if u < 1<<7 {
		u = (u << 1)
		*b = append(*b, uint8(u))
		return
	}
	if u < 1<<14 {
		u = (u << 2) | 1
		*b = append(*b, uint8(u), uint8(u>>8))
		return
	}
	u = (u << 2) | 3
	*b = append(*b, uint8(u), uint8(u>>8), uint8(u>>16), uint8(u>>24))
}

func (b *buffer) encodeReal(f float32) int {
	if u := uint32(f); float32(u) == f && u < 1<<14 {
		if u < 1<<7 {
			u = (u << 1)
			*b = append(*b, uint8(u))
			return 1
		}
		u = (u << 2) | 1
		*b = append(*b, uint8(u), uint8(u>>8))
		return 2
	}
	b.encode4ByteReal(f)
	return 4
}

func (b *buffer) encode4ByteReal(f float32) {
	u := math.Float32bits(f)

	// Round the fractional bits (the low 23 bits) to the nearest multiple of
	// 4, being careful not to overflow into the upper bits.
	v := u & 0x007fffff
	if v < 0x007ffffe {
		v += 2
	}
	u = (u & 0xff800000) | v

	// A 4 byte encoding has the low two bits set.
	u |= 0x03
	*b = append(*b, uint8(u), uint8(u>>8), uint8(u>>16), uint8(u>>24))
}

func (b *buffer) encodeCoordinate(f float32) int {
	if i := int32(f); -64 <= i && i < +64 && float32(i) == f {
		u := uint32(i + 64)
		u = (u << 1)
		*b = append(*b, uint8(u))
		return 1
	}
	if i := int32(f * 64); -128*64 <= i && i < +128*64 && float32(i) == f*64 {
		u := uint32(i + 128*64)
		u = (u << 2) | 1
		*b = append(*b, uint8(u), uint8(u>>8))
		return 2
	}
	b.encode4ByteReal(f)
	return 4
}

func (b *buffer) encodeAngle(f float32) int {
	// Normalize f to the range [0, 1).
	g := float64(f)
	g -= math.Floor(g)
	return b.encodeZeroToOne(float32(g))
}

func (b *buffer) encodeZeroToOne(f float32) int {
	if u := uint32(f * 15120); float32(u) == f*15120 && u < 15120 {
		if u%126 == 0 {
			u = ((u / 126) << 1)
			*b = append(*b, uint8(u))
			return 1
		}
		u = (u << 2) | 1
		*b = append(*b, uint8(u), uint8(u>>8))
		return 2
	}
	b.encode4ByteReal(f)
	return 4
}

func (b *buffer) encodeColor1(c ivg.Color) {
	if x, ok := c.Encode1(); ok {
		*b = append(*b, x)
		return
	}
	// Default to opaque black.
	*b = append(*b, 0x00)
}

func (b *buffer) encodeColor2(c ivg.Color) {
	if x, ok := c.Encode2(); ok {
		*b = append(*b, x[0], x[1])
		return
	}
	// Default to opaque black.
	*b = append(*b, 0x00, 0x0f)
}

func (b *buffer) encodeColor3Direct(c ivg.Color) {
	if x, ok := c.Encode3Direct(); ok {
		*b = append(*b, x[0], x[1], x[2])
		return
	}
	// Default to opaque black.
	*b = append(*b, 0x00, 0x00, 0x00)
}

func (b *buffer) encodeColor4(c ivg.Color) {
	if x, ok := c.Encode4(); ok {
		*b = append(*b, x[0], x[1], x[2], x[3])
		return
	}
	// Default to opaque black.
	*b = append(*b, 0x00, 0x00, 0x00, 0xff)
}

func (b *buffer) encodeColor3Indirect(c ivg.Color) {
	if x, ok := c.Encode3Indirect(); ok {
		*b = append(*b, x[0], x[1], x[2])
		return
	}
	// Default to opaque black.
	*b = append(*b, 0x00, 0x00, 0x00)
}
