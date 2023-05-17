package mdicons

type SVG struct {
	Width   float32 `xml:"width,attr"`
	Height  float32 `xml:"height,attr"`
	ViewBox string  `xml:"viewBox,attr"`
	Paths   []Path  `xml:"path"`
	// Some of the SVG files contain <circle> elements, not just <path>
	// elements. IconVG doesn't have circles per se. Instead, we convert such
	// circles to be paired arcTo commands, tacked on to the first path.
	//
	// In general, this isn't correct if the circles and the path overlap, but
	// that doesn't happen in the specific case of the Material Design icons.
	Circles []Circle `xml:"circle"`
}

type Path struct {
	D           string   `xml:"d,attr"`
	Fill        string   `xml:"fill,attr"`
	FillOpacity *float32 `xml:"fill-opacity,attr"`
	Opacity     *float32 `xml:"opacity,attr"`
}

type Circle struct {
	Cx float32 `xml:"cx,attr"`
	Cy float32 `xml:"cy,attr"`
	R  float32 `xml:"r,attr"`
}

type Statistics struct {
	VarNames        []string
	Failures        []string
	TotalFiles      int
	TotalIVGBytes   int
	TotalPNG24Bytes int
	TotalPNG48Bytes int
	TotalSVGBytes   int
}

func (s Statistics) Add(other Statistics) Statistics {
	return Statistics{
		VarNames:        append(s.VarNames, other.VarNames...),
		Failures:        append(s.Failures, other.Failures...),
		TotalFiles:      s.TotalFiles + other.TotalFiles,
		TotalIVGBytes:   s.TotalIVGBytes + other.TotalIVGBytes,
		TotalPNG24Bytes: s.TotalPNG24Bytes + other.TotalPNG24Bytes,
		TotalPNG48Bytes: s.TotalPNG48Bytes + other.TotalPNG48Bytes,
		TotalSVGBytes:   s.TotalSVGBytes + other.TotalSVGBytes}
}
