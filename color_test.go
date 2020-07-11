// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivg_test

import (
	"image/color"
	"testing"

	"github.com/reactivego/ivg"
)

func TestBlendColor(t *testing.T) {
	// This example comes from doc.go. Look for "orange" in the "Colors"
	// section.
	pal := [64]color.RGBA{
		2: {0xff, 0xcc, 0x80, 0xff}, // "Material Design Orange 200".
	}
	cReg := [64]color.RGBA{}
	got := ivg.BlendColor(0x40, 0x7f, 0x82).Resolve(&pal, &cReg)
	want := color.RGBA{0x40, 0x33, 0x20, 0x40} // 25% opaque "Orange 200", alpha-premultiplied.
	if got != want {
		t.Errorf("\ngot  %x\nwant %x", got, want)
	}
}
