// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decode

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/reactivego/ivg"
)

// overwriteTestdataFiles is temporarily set to true when adding new
// testdataTestCases.
const overwriteTestdataFiles = false

// disassemble returns a disassembly of an encoded IconVG graphic. Users of
// this package aren't expected to want to do this, so it lives in a _test.go
// file, but it can be useful for debugging.
func disassemble(src []byte) ([]byte, error) {
	w := new(bytes.Buffer)
	p := func(b []byte, format string, args ...interface{}) {
		const hex = "0123456789abcdef"
		var buf [14]byte
		for i := range buf {
			buf[i] = ' '
		}
		for i, x := range b {
			buf[3*i+0] = hex[x>>4]
			buf[3*i+1] = hex[x&0x0f]
		}
		w.Write(buf[:])
		fmt.Fprintf(w, format, args...)
	}
	m := ivg.Metadata{}
	if err := decode(nil, p, &m, false, buffer(src)); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

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
		ivgData, err := ioutil.ReadFile(filepath.FromSlash(tc.filename) + ".ivg")
		if err != nil {
			t.Errorf("%s: ReadFile: %v", tc.filename, err)
			continue
		}
		got, err := disassemble(ivgData)
		if err != nil {
			t.Errorf("%s: disassemble: %v", tc.filename, err)
			continue
		}
		wantFilename := filepath.FromSlash(tc.filename) + ".ivg.disassembly"
		if overwriteTestdataFiles {
			if err := ioutil.WriteFile(filepath.FromSlash(wantFilename), got, 0666); err != nil {
				t.Errorf("%s: WriteFile: %v", tc.filename, err)
			}
			continue
		}
		want, err := ioutil.ReadFile(wantFilename)
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
