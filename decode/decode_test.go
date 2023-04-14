// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decode

import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
)

// overwriteTestdataFiles is temporarily set to true when adding new
// testdataTestCases.
const overwriteTestdataFiles = false

func diffLines(t *testing.T, got, want string) {
	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")
	for i := 1; ; i++ {
		if len(gotLines) == 0 {
			t.Errorf("line %d:\ngot  %q\nwant %q", i, "", wantLines[0])
			return
		}
		if len(wantLines) == 0 {
			t.Errorf("line %d:\ngot  %q\nwant %q", i, gotLines[0], "")
			return
		}
		g, w := gotLines[0], wantLines[0]
		gotLines = gotLines[1:]
		wantLines = wantLines[1:]
		if g != w {
			t.Errorf("line %d:\ngot  %q\nwant %q", i, g, w)
			return
		}
	}
}

var testdataTestCases = []struct {
	filename string
	variants string
}{
	{"../testdata/action-info.lores", ""},
	{"../testdata/action-info.hires", ""},
	{"../testdata/arcs", ""},
	{"../testdata/blank", ""},
	{"../testdata/cowbell", ""},
	{"../testdata/elliptical", ""},
	{"../testdata/favicon", ";pink"},
	{"../testdata/gradient", ""},
	{"../testdata/lod-polygon", ";64"},
	{"../testdata/video-005.primitive", ""},
}

func TestDisassembly(t *testing.T) {
	for _, tc := range testdataTestCases {
		ivgData, err := os.ReadFile(filepath.FromSlash(tc.filename) + ".ivg")
		if err != nil {
			t.Errorf("%s: ReadFile: %v", tc.filename, err)
			continue
		}
		got, err := Disassemble(ivgData)
		if err != nil {
			t.Errorf("%s: disassemble: %v", tc.filename, err)
			continue
		}
		wantFilename := filepath.FromSlash(tc.filename) + ".ivg.disassembly"
		if overwriteTestdataFiles {
			if err := os.WriteFile(filepath.FromSlash(wantFilename), got, 0666); err != nil {
				t.Errorf("%s: WriteFile: %v", tc.filename, err)
			}
			continue
		}
		want, err := os.ReadFile(wantFilename)
		if err != nil {
			t.Errorf("%s: ReadFile: %v", tc.filename, err)
			continue
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s: got:\n%s\nwant:\n%s", tc.filename, got, want)
			diffLines(t, string(got), string(want))
		}
	}
}

// The IconVG decoder and encoder are expected to be completely deterministic,
// so check that we get the original bytes after a decode + encode round-trip.
func TestDecodeEncodeRoundTrip(t *testing.T) {
	for _, tc := range testdataTestCases {
		ivgData, err := os.ReadFile(filepath.FromSlash(tc.filename) + ".ivg")
		if err != nil {
			t.Errorf("%s: ReadFile: %v", tc.filename, err)
			continue
		}
		var e resolutionPreservingEncoder
		e.HighResolutionCoordinates = strings.HasSuffix(tc.filename, ".hires")
		if err := Decode(&e, ivgData); err != nil {
			t.Errorf("%s: Decode: %v", tc.filename, err)
			continue
		}
		got, err := e.Bytes()
		if err != nil {
			t.Errorf("%s: Encoder.Bytes: %v", tc.filename, err)
			continue
		}
		if want := ivgData; !bytes.Equal(got, want) {
			t.Errorf("%s:\ngot  %d bytes (on GOOS=%s GOARCH=%s, using compiler %q):\n% x\nwant %d bytes:\n% x",
				tc.filename, len(got), runtime.GOOS, runtime.GOARCH, runtime.Compiler, got, len(want), want)
			gotDisasm, err1 := Disassemble(got)
			wantDisasm, err2 := Disassemble(want)
			if err1 == nil && err2 == nil {
				diffLines(t, string(gotDisasm), string(wantDisasm))
			}
		}
	}
}

// resolutionPreservingEncoder is an Encoder
// whose Reset method keeps prior resolution.
type resolutionPreservingEncoder struct {
	encode.Encoder
}

// Reset resets the Encoder for the given Metadata.
//
// Unlike Encoder.Reset, it leaves the value
// of e.HighResolutionCoordinates unmodified.
func (e *resolutionPreservingEncoder) Reset(viewbox ivg.ViewBox, palette [64]color.RGBA) {
	orig := e.HighResolutionCoordinates
	e.Encoder.Reset(viewbox, palette)
	e.HighResolutionCoordinates = orig
}
