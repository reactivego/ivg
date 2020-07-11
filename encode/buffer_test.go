// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encode

import (
	"image/color"
	"math"
	"testing"

	"github.com/reactivego/ivg"
)

var naturalTestCases = []struct {
	in    buffer
	want  uint32
	wantN int
}{{
	buffer{},
	0,
	0,
}, {
	buffer{0x28},
	20,
	1,
}, {
	buffer{0x59},
	0,
	0,
}, {
	buffer{0x59, 0x83},
	8406,
	2,
}, {
	buffer{0x07, 0x00, 0x80},
	0,
	0,
}, {
	buffer{0x07, 0x00, 0x80, 0x3f},
	266338305,
	4,
}}

func TestEncodeNatural(t *testing.T) {
	for _, tc := range naturalTestCases {
		if tc.wantN == 0 {
			continue
		}
		var b buffer
		b.encodeNatural(tc.want)
		if got, want := string(b), string(tc.in); got != want {
			t.Errorf("value=%v:\ngot  % x\nwant % x", tc.want, got, want)
		}
	}
}

var realTestCases = []struct {
	in    buffer
	want  float32
	wantN int
}{{
	buffer{0x28},
	20,
	1,
}, {
	buffer{0x59, 0x83},
	8406,
	2,
}, {
	buffer{0x07, 0x00, 0x80, 0x3f},
	1.000000476837158203125,
	4,
}}

func TestEncodeReal(t *testing.T) {
	for _, tc := range realTestCases {
		var b buffer
		b.encodeReal(tc.want)
		if got, want := string(b), string(tc.in); got != want {
			t.Errorf("value=%v:\ngot  % x\nwant % x", tc.want, got, want)
		}
	}
}

var coordinateTestCases = []struct {
	in    buffer
	want  float32
	wantN int
}{{
	buffer{0x8e},
	7,
	1,
}, {
	buffer{0x81, 0x87},
	7.5,
	2,
}, {
	buffer{0x03, 0x00, 0xf0, 0x40},
	7.5,
	4,
}, {
	buffer{0x07, 0x00, 0xf0, 0x40},
	7.5000019073486328125,
	4,
}}

func TestEncodeCoordinate(t *testing.T) {
	for _, tc := range coordinateTestCases {
		if tc.want == 7.5 && tc.wantN == 4 {
			// 7.5 can be encoded in fewer than 4 bytes.
			continue
		}
		var b buffer
		b.encodeCoordinate(tc.want)
		if got, want := string(b), string(tc.in); got != want {
			t.Errorf("value=%v:\ngot  % x\nwant % x", tc.want, got, want)
		}
	}
}

func trunc(x float32) float32 {
	u := math.Float32bits(x)
	u &^= 0x03
	return math.Float32frombits(u)
}

var zeroToOneTestCases = []struct {
	in    buffer
	want  float32
	wantN int
}{{
	buffer{0x0a},
	1.0 / 24,
	1,
}, {
	buffer{0x41, 0x1a},
	1.0 / 9,
	2,
}, {
	buffer{0x63, 0x0b, 0x36, 0x3b},
	trunc(1.0 / 360),
	4,
}}

func TestEncodeZeroToOne(t *testing.T) {
	for _, tc := range zeroToOneTestCases {
		var b buffer
		b.encodeZeroToOne(tc.want)
		if got, want := string(b), string(tc.in); got != want {
			t.Errorf("value=%v:\ngot  % x\nwant % x", tc.want, got, want)
		}
	}
}

var colorTestCases = []struct {
	in     buffer
	encode func(*buffer, ivg.Color)
	want   ivg.Color
	wantN  int
}{{
	buffer{},
	(*buffer).encodeColor1,
	ivg.Color{},
	0,
}, {
	buffer{0x00},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0x00, 0x00, 0x00, 0xff}),
	1,
}, {
	buffer{0x30},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0x40, 0xff, 0xc0, 0xff}),
	1,
}, {
	buffer{0x7c},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0xff, 0xff, 0xff, 0xff}),
	1,
}, {
	buffer{0x7d},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0xc0, 0xc0, 0xc0, 0xc0}),
	1,
}, {
	buffer{0x7e},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0x80, 0x80, 0x80, 0x80}),
	1,
}, {
	buffer{0x7f},
	(*buffer).encodeColor1,
	ivg.RGBAColor(color.RGBA{0x00, 0x00, 0x00, 0x00}),
	1,
}, {
	buffer{0x80},
	(*buffer).encodeColor1,
	ivg.PaletteIndexColor(0x00),
	1,
}, {
	buffer{0xbf},
	(*buffer).encodeColor1,
	ivg.PaletteIndexColor(0x3f),
	1,
}, {
	buffer{0xc0},
	(*buffer).encodeColor1,
	ivg.CRegColor(0x00),
	1,
}, {
	buffer{0xff},
	(*buffer).encodeColor1,
	ivg.CRegColor(0x3f),
	1,
}, {
	buffer{0x01},
	(*buffer).encodeColor2,
	ivg.Color{},
	0,
}, {
	buffer{0x38, 0x0f},
	(*buffer).encodeColor2,
	ivg.RGBAColor(color.RGBA{0x33, 0x88, 0x00, 0xff}),
	2,
}, {
	buffer{0x00, 0x02},
	(*buffer).encodeColor3Direct,
	ivg.Color{},
	0,
}, {
	buffer{0x30, 0x66, 0x07},
	(*buffer).encodeColor3Direct,
	ivg.RGBAColor(color.RGBA{0x30, 0x66, 0x07, 0xff}),
	3,
}, {
	buffer{0x00, 0x00, 0x03},
	(*buffer).encodeColor4,
	ivg.Color{},
	0,
}, {
	buffer{0x30, 0x66, 0x07, 0x80},
	(*buffer).encodeColor4,
	ivg.RGBAColor(color.RGBA{0x30, 0x66, 0x07, 0x80}),
	4,
}, {
	buffer{0x00, 0x04},
	(*buffer).encodeColor3Indirect,
	ivg.Color{},
	0,
}, {
	buffer{0x40, 0x7f, 0x82},
	(*buffer).encodeColor3Indirect,
	ivg.BlendColor(0x40, 0x7f, 0x82),
	3,
}}

func TestEncodeColor(t *testing.T) {
	for _, tc := range colorTestCases {
		if tc.wantN == 0 {
			continue
		}
		var b buffer
		tc.encode(&b, tc.want)
		if got, want := string(b), string(tc.in); got != want {
			t.Errorf("value=%v:\ngot  % x\nwant % x", tc.want, got, want)
		}
	}
}
