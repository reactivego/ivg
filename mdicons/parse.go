// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mdicons

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"sort"
)

func Parse(mdicons string, pkg string, genData, genDataTest bool, outSize float32) error {
	f, err := os.Open(mdicons)
	if err != nil {
		return fmt.Errorf("%v\n\nDid you pass a directory to -mdicons <directory>?", err)
	}
	defer f.Close()
	infos, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	names := []string{}
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		name := info.Name()
		if name[0] == '.' {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)

	out := &bytes.Buffer{}

	// Generate data.go.
	if genData {
		out.WriteString(fmt.Sprintf("// generated by mdicons; DO NOT EDIT\n\npackage %s\n\n", pkg))

		stats := Statistics{}
		for _, name := range names {
			if stat, err := ParseDir(mdicons, name, outSize, out); err != nil {
				return err
			} else {
				stats = stats.Add(stat)
			}
		}

		fmt.Fprintf(out,
			"// In total, %d SVG bytes in %d files (%d PNG bytes at 24px * 24px,\n"+
				"// %d PNG bytes at 48px * 48px) converted to %d IconVG bytes.\n",
			stats.TotalSVGBytes, stats.TotalFiles, stats.TotalPNG24Bytes, stats.TotalPNG48Bytes, stats.TotalIVGBytes)

		if len(stats.Failures) != 0 {
			out.WriteString("\n/*\nFAILURES:\n\n")
			for _, failure := range stats.Failures {
				out.WriteString(failure)
				out.WriteByte('\n')
			}
			out.WriteString("\n*/")
		}

		raw := out.Bytes()
		formatted, err := format.Source(raw)
		if err != nil {
			return fmt.Errorf("gofmt failed: %v\n\nGenerated code:\n%s", err, raw)
		}
		if err := os.WriteFile("data.go", formatted, 0644); err != nil {
			return fmt.Errorf("WriteFile failed: %s", err)
		}

		// Generate data_test.go.
		if genDataTest {
			out.Reset()
			out.WriteString(fmt.Sprintf("// generated by mdicons; DO NOT EDIT\n\npackage %s\n\n", pkg))
			out.WriteString("var list = []struct{ name string; data []byte } {\n")
			for _, v := range stats.VarNames {
				fmt.Fprintf(out, "{%q, %s},\n", v, v)
			}
			out.WriteString("}\n\n")
			raw := out.Bytes()
			formatted, err := format.Source(raw)
			if err != nil {
				return fmt.Errorf("gofmt failed: %v\n\nGenerated code:\n%s", err, raw)
			}
			if err := os.WriteFile("data_test.go", formatted, 0644); err != nil {
				return fmt.Errorf("WriteFile failed: %s", err)
			}
		}
	}
	return nil
}