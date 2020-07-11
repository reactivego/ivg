// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ivg_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/render"
	"github.com/reactivego/ivg/raster/vector"
)

// overwriteTestdataFiles is temporarily set to true when adding new
// testdataTestCases.
const overwriteTestdataFiles = false

func encodePNG(dstFilename string, src image.Image) error {
	f, err := os.Create(dstFilename)
	if err != nil {
		return err
	}
	encErr := png.Encode(f, src)
	closeErr := f.Close()
	if encErr != nil {
		return encErr
	}
	return closeErr
}

func decodePNG(srcFilename string) (image.Image, error) {
	f, err := os.Open(srcFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

func checkApproxEqual(m0, m1 image.Image) error {
	diff := func(a, b uint32) uint32 {
		if a < b {
			return b - a
		}
		return a - b
	}

	bounds0 := m0.Bounds()
	bounds1 := m1.Bounds()
	if bounds0 != bounds1 {
		return fmt.Errorf("bounds differ: got %v, want %v", bounds0, bounds1)
	}
	for y := bounds0.Min.Y; y < bounds0.Max.Y; y++ {
		for x := bounds0.Min.X; x < bounds0.Max.X; x++ {
			r0, g0, b0, a0 := m0.At(x, y).RGBA()
			r1, g1, b1, a1 := m1.At(x, y).RGBA()

			// TODO: be more principled in picking this magic threshold, other
			// than what the difference is, in practice, in x/image/vector's
			// fixed and floating point rasterizer?
			const D = 0xffff * 12 / 100 // Diff threshold of 12%.

			if diff(r0, r1) > D || diff(g0, g1) > D || diff(b0, b1) > D || diff(a0, a1) > D {
				return fmt.Errorf("at (%d, %d):\n"+
					"got  RGBA %#04x, %#04x, %#04x, %#04x\n"+
					"want RGBA %#04x, %#04x, %#04x, %#04x",
					x, y, r0, g0, b0, a0, r1, g1, b1, a1)
			}
		}
	}
	return nil
}

var testdataTestCases = []struct {
	filename string
	variants string
}{
	{"testdata/action-info.lores", ""},
	{"testdata/action-info.hires", ""},
	{"testdata/arcs", ""},
	{"testdata/blank", ""},
	{"testdata/cowbell", ""},
	{"testdata/elliptical", ""},
	{"testdata/favicon", ";pink"},
	{"testdata/gradient", ""},
	{"testdata/lod-polygon", ";64"},
	{"testdata/video-005.primitive", ""},
}

func TestRasterizer(t *testing.T) {
	for _, tc := range testdataTestCases {
		ivgData, err := ioutil.ReadFile(filepath.FromSlash(tc.filename) + ".ivg")
		if err != nil {
			t.Errorf("%s: ReadFile: %v", tc.filename, err)
			continue
		}
		md, err := decode.DecodeMetadata(ivgData)
		if err != nil {
			t.Errorf("%s: DecodeMetadata: %v", tc.filename, err)
			continue
		}

		for _, variant := range strings.Split(tc.variants, ";") {
			length := 256
			if variant == "64" {
				length = 64
			}
			width, height := length, length
			if dx, dy := md.ViewBox.AspectRatio(); dx < dy {
				width = int(float32(length) * dx / dy)
			} else {
				height = int(float32(length) * dy / dx)
			}

			opts := &decode.DecodeOptions{}
			if variant == "pink" {
				pal := ivg.DefaultPalette
				pal[0] = color.RGBA{0xfe, 0x76, 0xea, 0xff}
				opts.Palette = &pal
			}

			bounds := image.Rect(0, 0, width, height)
			got := image.NewRGBA(bounds)
			var z render.Renderer
			z.SetRasterizer(vector.NewRasterizer(got, draw.Src), got.Bounds())
			if err := decode.Decode(&z, ivgData, opts); err != nil {
				t.Errorf("%s %q variant: Decode: %v", tc.filename, variant, err)
				continue
			}

			wantFilename := filepath.FromSlash(tc.filename)
			if variant != "" {
				wantFilename += "." + variant
			}
			wantFilename += ".png"
			if overwriteTestdataFiles {
				if err := encodePNG(filepath.FromSlash(wantFilename), got); err != nil {
					t.Errorf("%s %q variant: encodePNG: %v", tc.filename, variant, err)
				}
				continue
			}
			want, err := decodePNG(wantFilename)
			if err != nil {
				t.Errorf("%s %q variant: decodePNG: %v", tc.filename, variant, err)
				continue
			}
			if err := checkApproxEqual(got, want); err != nil {
				t.Errorf("%s %q variant: %v", tc.filename, variant, err)
				continue
			}
		}
	}
}
