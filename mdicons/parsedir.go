package mdicons

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var skippedFiles = map[[2]string]bool{
	// ic_play_circle_filled_white_48px.svg is just the same as
	// ic_play_circle_filled_48px.svg with an explicit fill="#fff".
	{"av", "ic_play_circle_filled_white_48px.svg"}: true,
}

func ParseDir(mdicons, dirName string, outSize float32, out *bytes.Buffer) (Statistics, error) {
	var stats Statistics

	fqPNGDirName := filepath.FromSlash(path.Join(mdicons, dirName, "1x_web"))
	fqSVGDirName := filepath.FromSlash(path.Join(mdicons, dirName, "svg/production"))
	f, err := os.Open(fqSVGDirName)
	if err != nil {
		return stats, nil
	}
	defer f.Close()

	infos, err := f.Readdir(-1)
	if err != nil {
		return stats, err
	}
	baseNames, fileNames, sizes := []string{}, map[string]string{}, map[string]int{}
	for _, info := range infos {
		name := info.Name()

		if !strings.HasPrefix(name, "ic_") || skippedFiles[[2]string{dirName, name}] {
			continue
		}
		size := 0
		switch {
		case strings.HasSuffix(name, "_12px.svg"):
			size = 12
		case strings.HasSuffix(name, "_18px.svg"):
			size = 18
		case strings.HasSuffix(name, "_24px.svg"):
			size = 24
		case strings.HasSuffix(name, "_36px.svg"):
			size = 36
		case strings.HasSuffix(name, "_48px.svg"):
			size = 48
		default:
			continue
		}

		baseName := name[3 : len(name)-9]
		if prevSize, ok := sizes[baseName]; ok {
			if size > prevSize {
				fileNames[baseName] = name
				sizes[baseName] = size
			}
		} else {
			fileNames[baseName] = name
			sizes[baseName] = size
			baseNames = append(baseNames, baseName)
		}
	}

	sort.Strings(baseNames)
	for _, baseName := range baseNames {
		fileName := fileNames[baseName]
		stat, err := ParseFile(filepath.Join(fqSVGDirName, fileName), dirName, baseName, float32(sizes[baseName]), outSize, out)
		if err == ErrSkip {
			continue
		}
		if err != nil {
			stat.Failures = append(stat.Failures, fmt.Sprintf("%v/svg/production/%v: %v", dirName, fileName, err))
			continue
		}
		stats = stats.Add(stat)
		totalPNG24Bytes, fail := pngSize(fqPNGDirName, dirName, baseName, 24)
		if fail != "" {
			stats.Failures = append(stats.Failures, fail)
		} else {
			stats.TotalPNG24Bytes += totalPNG24Bytes
		}
		totalPNG48Bytes, fail := pngSize(fqPNGDirName, dirName, baseName, 48)
		if fail != "" {
			stats.Failures = append(stats.Failures, fail)
		} else {
			stats.TotalPNG48Bytes += totalPNG48Bytes
		}
	}

	return stats, nil
}

func pngSize(fqPNGDirName, dirName, baseName string, targetSize int) (int, string) {
	for _, size := range [...]int{48, 24, 18} {
		if size > targetSize {
			continue
		}
		fInfo, err := os.Stat(filepath.Join(fqPNGDirName,
			fmt.Sprintf("ic_%s_black_%ddp.png", baseName, size)))
		if err != nil {
			continue
		}
		return int(fInfo.Size()), ""
	}
	return 0, fmt.Sprintf("no PNG found for %s/1x_web/ic_%s_black_{48,24,18}dp.png", dirName, baseName)
}
