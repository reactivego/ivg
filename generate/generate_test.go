package generate

import (
	"bytes"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
)

// overwriteTestdataFiles is temporarily set to true when adding new
// testdataTestCases.
const overwriteTestdataFiles = false

// TestOverwriteTestdataFilesIsFalse tests that any change to
// overwriteTestdataFiles is only temporary. Programmers are assumed to run "go
// test" before sending out for code review or committing code.
func TestOverwriteTestdataFilesIsFalse(t *testing.T) {
	if overwriteTestdataFiles {
		t.Errorf("overwriteTestdataFiles is true; do not commit code changes")
	}
}

func testEncode(t *testing.T, e *encode.Encoder, wantFilename string) {
	got, err := e.Bytes()
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}
	if overwriteTestdataFiles {
		if err := ioutil.WriteFile(filepath.FromSlash(wantFilename), got, 0666); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		return
	}
	want, err := ioutil.ReadFile(filepath.FromSlash(wantFilename))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("\ngot  %d bytes:\n% x\nwant %d bytes:\n% x", len(got), got, len(want), want)
	}
}

func encodePathData(e *encode.Encoder, d string, adj uint8, normalizeTo64X64 bool) error {
	var args [7]float32
	prevN, prevVerb := 0, byte(0)
	for first := true; d != "z"; first = false {
		n, verb, implicit := 0, d[0], false
		switch d[0] {
		case 'H', 'h', 'V', 'v':
			n = 1
		case 'L', 'M', 'l', 'm':
			n = 2
		case 'S', 's':
			n = 4
		case 'C', 'c':
			n = 6
		case 'A', 'a':
			n = 7
		case 'z':
			n = 0
		default:
			if prevVerb == '\x00' {
				panic("unrecognized verb")
			}
			n, verb, implicit = prevN, prevVerb, true
		}
		prevN, prevVerb = n, verb
		if prevVerb == 'M' {
			prevVerb = 'L'
		} else if prevVerb == 'm' {
			prevVerb = 'l'
		}
		if !implicit {
			d = d[1:]
		}

		for i := 0; i < n; i++ {
			nDots := 0
			if d[0] == '.' {
				nDots = 1
			}
			j := 1
			for ; ; j++ {
				switch d[j] {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					continue
				case '.':
					nDots++
					if nDots == 1 {
						continue
					}
				}
				break
			}
			f, err := strconv.ParseFloat(d[:j], 64)
			if err != nil {
				return err
			}
			args[i] = float32(f)
			for ; d[j] == ' ' || d[j] == ','; j++ {
			}
			d = d[j:]
		}

		if normalizeTo64X64 {
			// The original SVG is 32x32 units, with the top left being (0, 0).
			// Normalize to 64x64 units, with the center being (0, 0).
			if verb == 'A' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[2] /= 360
				args[5] = 2*args[5] - 32
				args[6] = 2*args[6] - 32
			} else if verb == 'a' {
				args[0] = 2 * args[0]
				args[1] = 2 * args[1]
				args[2] /= 360
				args[5] = 2 * args[5]
				args[6] = 2 * args[6]
			} else if first || ('A' <= verb && verb <= 'Z') {
				for i := range args {
					args[i] = 2*args[i] - 32
				}
			} else {
				for i := range args {
					args[i] = 2 * args[i]
				}
			}
		} else if verb == 'A' || verb == 'a' {
			args[2] /= 360
		}

		if first {
			first = false
			e.StartPath(adj, args[0], args[1])
			continue
		}
		switch verb {
		case 'H':
			e.AbsHLineTo(args[0])
		case 'h':
			e.RelHLineTo(args[0])
		case 'V':
			e.AbsVLineTo(args[0])
		case 'v':
			e.RelVLineTo(args[0])
		case 'L':
			e.AbsLineTo(args[0], args[1])
		case 'l':
			e.RelLineTo(args[0], args[1])
		case 'm':
			e.ClosePathRelMoveTo(args[0], args[1])
		case 'S':
			e.AbsSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 's':
			e.RelSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 'C':
			e.AbsCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'c':
			e.RelCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'A':
			e.AbsArcTo(args[0], args[1], args[2], args[3] != 0, args[4] != 0, args[5], args[6])
		case 'a':
			e.RelArcTo(args[0], args[1], args[2], args[3] != 0, args[4] != 0, args[5], args[6])
		case 'z':
			// No-op.
		default:
			panic("unrecognized verb")
		}
	}
	e.ClosePathEndPath()
	return nil
}

var cowbellGradients = []struct {
	radial bool

	// Linear gradient coefficients.
	x1, y1 float32
	x2, y2 float32
	tx, ty float32

	// Radial gradient coefficients.
	cx, cy, r float32
	transform Aff3

	stops []GradientStop
}{{
	// The 0th element is unused.
}, {
	radial: true,
	cx:     -102.14,
	cy:     20.272,
	r:      18.012,
	transform: Aff3{
		.33050, -.50775, 65.204,
		.17296, .97021, 16.495,
	},
	stops: []GradientStop{
		{Offset: 0, Color: color.RGBA{0xed, 0xd4, 0x00, 0xff}},
		{Offset: 1, Color: color.RGBA{0xfc, 0xe9, 0x4f, 0xff}},
	},
}, {
	radial: true,
	cx:     -97.856,
	cy:     26.719,
	r:      18.61,
	transform: Aff3{
		.35718, -.11527, 51.072,
		.044280, .92977, 7.6124,
	},
	stops: []GradientStop{
		{Offset: 0, Color: color.RGBA{0xed, 0xd4, 0x00, 0xff}},
		{Offset: 1, Color: color.RGBA{0xfc, 0xe9, 0x4f, 0xff}},
	},
}, {
	x1: -16.183,
	y1: 35.723,
	x2: -18.75,
	y2: 29.808,
	tx: 48.438,
	ty: -.22321,
	stops: []GradientStop{
		{Offset: 0, Color: color.RGBA{0x39, 0x21, 0x00, 0xff}},
		{Offset: 1, Color: color.RGBA{0x0f, 0x08, 0x00, 0xff}},
	},
}}

var cowbellSVGData = []struct {
	rgba      color.RGBA
	gradient  int
	d         string
	transform *Aff3
}{{
	gradient: 2,
	d:        "m5.6684 17.968l.265-4.407 13.453 19.78.301 8.304-14.019-23.677z",
}, {
	gradient: 1,
	d:        "m19.299 33.482l-13.619-19.688 3.8435-2.684.0922-2.1237 4.7023-2.26 2.99 1.1274 4.56-1.4252 20.719 16.272-23.288 10.782z",
}, {
	rgba: color.RGBA{0xfd * 127 / 255, 0xee * 127 / 255, 0x74 * 127 / 255, 127},
	d:    "m19.285 32.845l-13.593-19.079 3.995-2.833.1689-2.0377 1.9171-.8635 18.829 18.965-11.317 5.848z",
}, {
	rgba: color.RGBA{0xc4, 0xa0, 0x00, 0xff},
	d:    "m19.211 40.055c-.11-.67-.203-2.301-.205-3.624l-.003-2.406-2.492-3.769c-3.334-5.044-11.448-17.211-9.6752-14.744.3211.447 1.6961 2.119 2.1874 2.656.4914.536 1.3538 1.706 1.9158 2.6 2.276 3.615 8.232 12.056 8.402 12.056.1 0 10.4-5.325 11.294-5.678.894-.354 11.25-4.542 11.45-4.342.506.506 1.27 7.466.761 8.08-.392.473-5.06 3.672-10.256 6.121-5.195 2.45-11.984 4.269-12.594 4.269-.421 0-.639-.338-.785-1.219z",
}, {
	gradient: 3,
	d:        "m19.825 33.646c.422-.68 10.105-5.353 10.991-5.753s9.881-4.123 10.468-4.009c.512.099.844 6.017.545 6.703-.23.527-8.437 4.981-9.516 5.523-1.225.616-11.642 4.705-12.145 4.369-.553-.368-.707-6.245-.343-6.833z",
}, {
	rgba: color.RGBA{0x00, 0x00, 0x00, 0xff},
	d:    "m21.982 5.8789-4.865 1.457-2.553-1.1914-5.3355 2.5743l-.015625.29688-.097656 1.8672-4.1855 2.7383.36719 4.5996.054687.0957s3.2427 5.8034 6.584 11.654c1.6707 2.9255 3.3645 5.861 4.6934 8.0938.66442 1.1164 1.2366 2.0575 1.6719 2.7363.21761.33942.40065.6121.54883.81641.07409.10215.13968.18665.20312.25976.06345.07312.07886.13374.27148.22461.27031.12752.38076.06954.54102.04883.16025-.02072.34015-.05724.55078-.10938.42126-.10427.95998-.26728 1.584-.4707 1.248-.40685 2.8317-.97791 4.3926-1.5586 3.1217-1.1614 6.1504-2.3633 6.1504-2.3633l.02539-.0098.02539-.01367s2.5368-1.3591 5.1211-2.8027c1.2922-.72182 2.5947-1.4635 3.6055-2.0723.50539-.30438.93732-.57459 1.2637-.79688.16318-.11114.29954-.21136.41211-.30273.11258-.09138.19778-.13521.30273-.32617.16048-.292.13843-.48235.1543-.78906s.01387-.68208.002-1.1094c-.02384-.8546-.09113-1.9133-.17188-2.9473-.161-2.067-.373-4.04-.373-4.04l-.021-.211-20.907-16.348zm-.209 1.1055 20.163 15.766c.01984.1875.19779 1.8625.34961 3.8066.08004 1.025.14889 2.0726.17188 2.8965.01149.41192.01156.76817-.002 1.0293-.01351.26113-.09532.47241-.0332.35938.05869-.10679.01987-.0289-.05664.0332s-.19445.14831-.34375.25c-.29859.20338-.72024.46851-1.2168.76758-.99311.59813-2.291 1.3376-3.5781 2.0566-2.5646 1.4327-5.0671 2.7731-5.0859 2.7832-.03276.01301-3.0063 1.1937-6.0977 2.3438-1.5542.5782-3.1304 1.1443-4.3535 1.543-.61154.19936-1.1356.35758-1.5137.45117-.18066.04472-.32333.07255-.41992.08594-.02937-.03686-.05396-.06744-.0957-.125-.128-.176-.305-.441-.517-.771-.424-.661-.993-1.594-1.655-2.705-1.323-2.223-3.016-5.158-4.685-8.08-3.3124-5.8-6.4774-11.465-6.5276-11.555l-.3008-3.787 4.1134-2.692.109-2.0777 4.373-2.1133 2.469 1.1523 4.734-1.4179z",
}}

func inv(x *Aff3) Aff3 {
	invDet := 1 / (x[0]*x[4] - x[1]*x[3])
	return Aff3{
		+x[4] * invDet,
		-x[1] * invDet,
		(x[1]*x[5] - x[2]*x[4]) * invDet,
		-x[3] * invDet,
		+x[0] * invDet,
		(x[2]*x[3] - x[0]*x[5]) * invDet,
	}
}

func TestEncodeCowbell(t *testing.T) {
	var e encode.Encoder
	var gen Generator
	gen.SetDestination(&e)
	gen.Reset(
		ivg.ViewBox{
			MinX: 0, MinY: 0,
			MaxX: +48, MaxY: +48,
		},
		&ivg.DefaultPalette,
	)

	for _, data := range cowbellSVGData {
		if data.rgba != (color.RGBA{}) {
			gen.SetCReg(0, false, ivg.RGBAColor(data.rgba))
		} else if data.gradient != 0 {
			g := cowbellGradients[data.gradient]
			if g.radial {
				iform := inv(&g.transform)
				iform[2] -= g.cx
				iform[5] -= g.cy
				for i := range iform {
					iform[i] /= g.r
				}
				gen.SetGradient(10, 10, GradientShapeRadial, iform, GradientSpreadPad, g.stops)
			} else {
				x1 := g.x1 + g.tx
				y1 := g.y1 + g.ty
				x2 := g.x2 + g.tx
				y2 := g.y2 + g.ty
				gen.SetLinearGradient(10, 10, x1, y1, x2, y2, GradientSpreadPad, g.stops)
			}
		}

		if err := encodePathData(&e, data.d, 0, false); err != nil {
			t.Fatal(err)
		}
	}

	testEncode(t, &e, "../testdata/cowbell.ivg")
}

func TestEncodeElliptical(t *testing.T) {
	var e encode.Encoder
	var gen Generator
	gen.SetDestination(&e)

	const (
		cx, cy = -20, -10
		rx, ry = 0, 24
		sx, sy = 30, 15
	)

	gen.SetEllipticalGradient(10, 10, cx, cy, rx, ry, sx, sy, GradientSpreadReflect, []GradientStop{
		{Offset: 0, Color: color.RGBA{0xc0, 0x00, 0x00, 0xff}},
		{Offset: 1, Color: color.RGBA{0x00, 0x00, 0xc0, 0xff}},
	})
	gen.StartPath(0, -32, -32)
	gen.AbsHLineTo(+32)
	gen.AbsVLineTo(+32)
	gen.AbsHLineTo(-32)
	gen.ClosePathEndPath()

	gen.SetCReg(0, false, ivg.RGBAColor(color.RGBA{0xff, 0xff, 0xff, 0xff}))
	diamond := func(x, y float32) {
		gen.StartPath(0, x-1, y)
		gen.AbsLineTo(x, y-1)
		gen.AbsLineTo(x+1, y)
		gen.AbsLineTo(x, y+1)
		gen.ClosePathEndPath()
	}
	diamond(cx, cy)
	diamond(cx+rx, cy+ry)
	diamond(cx+sx, cy+sy)

	testEncode(t, &e, "../testdata/elliptical.ivg")
}

func TestEncodeGradient(t *testing.T) {
	rgb := []GradientStop{
		{Offset: 0.00, Color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
		{Offset: 1.00, Color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}
	cmy := []GradientStop{
		{Offset: 0.00, Color: color.RGBA{0x00, 0xff, 0xff, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0xff, 0x00, 0xff, 0xff}},
		{Offset: 0.75, Color: color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{Offset: 1.00, Color: color.RGBA{0xff, 0xff, 0x00, 0xff}},
	}

	var e encode.Encoder
	var gen Generator
	gen.SetDestination(&e)

	gen.SetLinearGradient(10, 10, -12, -30, +12, -18, GradientSpreadNone, rgb)
	gen.StartPath(0, -30, -30)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(-18)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetLinearGradient(10, 10, -12, -14, +12, -2, GradientSpreadPad, cmy)
	gen.StartPath(0, -30, -14)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(-2)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetCircularGradient(10, 10, -8, 8, 0, 16, GradientSpreadReflect, rgb)
	gen.StartPath(0, -30, +2)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(+14)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetCircularGradient(10, 10, -8, 24, 0, 16, GradientSpreadRepeat, cmy)
	gen.StartPath(0, -30, +18)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(+30)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	testEncode(t, &e, "../testdata/gradient.ivg")
}