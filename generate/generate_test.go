package generate

import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
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
		if err := os.WriteFile(filepath.FromSlash(wantFilename), got, 0666); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		return
	}
	want, err := os.ReadFile(filepath.FromSlash(wantFilename))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("\ngot  %d bytes:\n% x\nwant %d bytes:\n% x", len(got), got, len(want), want)
	}
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
		ivg.DefaultPalette,
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
				gen.SetGradient(GradientShapeRadial, GradientSpreadPad, g.stops, iform)
			} else {
				x1 := g.x1 + g.tx
				y1 := g.y1 + g.ty
				x2 := g.x2 + g.tx
				y2 := g.y2 + g.ty
				gen.SetLinearGradient(x1, y1, x2, y2, GradientSpreadPad, g.stops)
			}
		}

		if err := gen.SetPathData(data.d, 0); err != nil {
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

	gen.SetEllipticalGradient(cx, cy, rx, ry, sx, sy, GradientSpreadReflect, []GradientStop{
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

	gen.SetLinearGradient(-12, -30, +12, -18, GradientSpreadNone, rgb)
	gen.StartPath(0, -30, -30)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(-18)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetLinearGradient(-12, -14, +12, -2, GradientSpreadPad, cmy)
	gen.StartPath(0, -30, -14)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(-2)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetCircularGradient(-8, 8, 0, 16, GradientSpreadReflect, rgb)
	gen.StartPath(0, -30, +2)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(+14)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	gen.SetCircularGradient(-8, 24, 0, 16, GradientSpreadRepeat, cmy)
	gen.StartPath(0, -30, +18)
	gen.AbsHLineTo(+30)
	gen.AbsVLineTo(+30)
	gen.AbsHLineTo(-30)
	gen.ClosePathEndPath()

	testEncode(t, &e, "../testdata/gradient.ivg")
}

var faviconColors = []color.RGBA{
	{0x76, 0xe1, 0xfe, 0xff},
	{0x38, 0x4e, 0x54, 0xff},
	{0xff, 0xff, 0xff, 0xff},
	{0x17, 0x13, 0x11, 0xff},
	{0x00, 0x00, 0x00, 0x54},
	{0xff, 0xfc, 0xfb, 0xff},
	{0xc3, 0x8c, 0x74, 0xff},
	{0x23, 0x20, 0x1f, 0xff},
}

var faviconSVGData = []struct {
	faviconColorsIndex int
	d                  string
}{{
	faviconColorsIndex: 1,
	d:                  "m16.092 1.002c-1.1057.01-2.2107.048844-3.3164.089844-2.3441.086758-4.511.88464-6.2832 2.1758a3.8208 3.5794 29.452 0 0 -.8947 -.6856 3.8208 3.5794 29.452 0 0 -5.0879 1.2383 3.8208 3.5794 29.452 0 0 1.5664 4.9961 3.8208 3.5794 29.452 0 0 .3593 .1758c-.2784.9536-.4355 1.9598-.4355 3.0078v20h28v-20c0-1.042-.152-2.0368-.418-2.9766a3.5794 3.8208 60.548 0 0 .43359 -.20703 3.5794 3.8208 60.548 0 0 1.5684 -4.9961 3.5794 3.8208 60.548 0 0 -5.0879 -1.2383 3.5794 3.8208 60.548 0 0 -.92969 .72461c-1.727-1.257-3.843-2.0521-6.1562-2.2148-1.1058-.078-2.2126-.098844-3.3184-.089844z",
}, {
	faviconColorsIndex: 0,
	d:                  "m16 3c-4.835 0-7.9248 1.0791-9.7617 2.8906-.4777-.4599-1.2937-1.0166-1.6309-1.207-.9775-.5520-2.1879-.2576-2.7051.6582-.5171.9158-.1455 2.1063.8321 2.6582.2658.1501 1.2241.5845 1.7519.7441-.3281.9946-.4863 2.0829-.4863 3.2559v20h24c-.049-7.356 0-18 0-20 0-1.209-.166-2.3308-.516-3.3496.539-.2011 1.243-.5260 1.463-.6504.978-.5519 1.351-1.7424.834-2.6582s-1.729-1.2102-2.707-.6582c-.303.1711-.978.6356-1.463 1.0625-1.854-1.724-4.906-2.7461-9.611-2.7461z",
}, {
	faviconColorsIndex: 1,
	d:                  "m3.0918 5.9219c-.060217.00947-.10772.020635-.14648.033203-.019384.00628-.035462.013581-.052734.021484-.00864.00395-.019118.00825-.03125.015625-.00607.00369-.011621.00781-.021484.015625-.00493.00391-.017342.015389-.017578.015625-.0002366.0002356-.025256.031048-.025391.03125a.19867 .19867 0 0 0 .26367 .28320c.0005595-.0002168.00207-.00128.00391-.00195a.19867 .19867 0 0 0 .00391 -.00195c.015939-.00517.045148-.013113.085937-.019531.081581-.012836.20657-.020179.36719.00391.1020.0152.2237.0503.3535.0976-.3277.0694-.5656.1862-.7227.3145-.1143.0933-.1881.1903-.2343.2695-.023099.0396-.039499.074216-.050781.10547-.00564.015626-.00989.029721-.013672.046875-.00189.00858-.00458.017085-.00586.03125-.0006392.00708-.0005029.014724 0 .027344.0002516.00631.00192.023197.00195.023437.0000373.0002412.0097.036937.00977.037109a.19867 .19867 0 0 0 .38477 -.039063 .19867 .19867 0 0 0 0 -.00195c.00312-.00751.00865-.015947.017578-.03125.0230-.0395.0660-.0977.1425-.1601.1530-.1250.4406-.2702.9863-.2871a.19930 .19930 0 0 0 .082031 -.019531c.12649.089206.25979.19587.39844.32422a.19867 .19867 0 1 0 .2696 -.2911c-.6099-.5646-1.1566-.7793-1.5605-.8398-.2020-.0303-.3679-.0229-.4883-.0039z",
}, {
	faviconColorsIndex: 1,
	d:                  "m28.543 5.8203c-.12043-.018949-.28631-.026379-.48828.00391-.40394.060562-.94869.27524-1.5586.83984a.19867 .19867 0 1 0 .26953 .29102c.21354-.19768.40814-.33222.59180-.44141.51624.023399.79659.16181.94531.28320.07652.062461.11952.12063.14258.16016.0094.016037.01458.025855.01758.033203a.19867 .19867 0 0 0 .38476 .039063c.000062-.0001719.0097-.036868.0098-.037109.000037-.0002412.0017-.017125.002-.023437.000505-.012624.000639-.020258 0-.027344-.0013-.01417-.004-.022671-.0059-.03125-.0038-.017158-.008-.031248-.01367-.046875-.01128-.031254-.02768-.067825-.05078-.10742-.04624-.079195-.12003-.17424-.23437-.26758-.11891-.097066-.28260-.18832-.49609-.25781.01785-.00328.03961-.011119.05664-.013672.16062-.024082.28561-.016738.36719-.00391.03883.00611.06556.012409.08203.017578.000833.0002613.0031.0017.0039.00195a.19867 .19867 0 0 0 .271 -.2793c-.000135-.0002016-.02515-.031014-.02539-.03125-.000236-.0002356-.01265-.011717-.01758-.015625-.0099-.00782-.01737-.01194-.02344-.015625-.01213-.00737-.02066-.011673-.0293-.015625-.01727-.0079-.03336-.013247-.05273-.019531-.03877-.012568-.08822-.025682-.14844-.035156z",
}, {
	faviconColorsIndex: 2,
	d:                  "m15.171 9.992a4.8316 4.8316 0 0 1 -4.832 4.832 4.8316 4.8316 0 0 1 -4.8311 -4.832 4.8316 4.8316 0 0 1 4.8311 -4.8316 4.8316 4.8316 0 0 1 4.832 4.8316z",
}, {
	faviconColorsIndex: 2,
	d:                  "m25.829 9.992a4.6538 4.6538 0 0 1 -4.653 4.654 4.6538 4.6538 0 0 1 -4.654 -4.654 4.6538 4.6538 0 0 1 4.654 -4.6537 4.6538 4.6538 0 0 1 4.653 4.6537z",
}, {
	faviconColorsIndex: 3,
	d:                  "m14.377 9.992a1.9631 1.9631 0 0 1 -1.963 1.963 1.9631 1.9631 0 0 1 -1.963 -1.963 1.9631 1.9631 0 0 1 1.963 -1.963 1.9631 1.9631 0 0 1 1.963 1.963z",
}, {
	faviconColorsIndex: 3,
	d:                  "m25.073 9.992a1.9631 1.9631 0 0 1 -1.963 1.963 1.9631 1.9631 0 0 1 -1.963 -1.963 1.9631 1.9631 0 0 1 1.963 -1.963 1.9631 1.9631 0 0 1 1.963 1.963z",
}, {
	faviconColorsIndex: 4,
	d:                  "m14.842 15.555h2.2156c.40215 0 .72590.3237.72590.7259v2.6545c0 .4021-.32375.7259-.72590.7259h-2.2156c-.40215 0-.72590-.3238-.72590-.7259v-2.6545c0-.4022.32375-.7259.72590-.7259z",
}, {
	faviconColorsIndex: 5,
	d:                  "m14.842 14.863h2.2156c.40215 0 .72590.3238.72590.7259v2.6546c0 .4021-.32375.7259-.72590.7259h-2.2156c-.40215 0-.72590-.3238-.72590-.7259v-2.6546c0-.4021.32375-.7259.72590-.7259z",
}, {
	faviconColorsIndex: 4,
	d:                  "m20 16.167c0 .838-.87123 1.2682-2.1448 1.1659-.02366 0-.04795-.6004-.25415-.5832-.50367.042-1.0959-.02-1.686-.02-.61294 0-1.2063.1826-1.6855.017-.11023-.038-.17830.5838-.26153.5816-1.2437-.033-2.0788-.3383-2.0788-1.1618 0-1.2118 1.8156-2.1941 4.0554-2.1941 2.2397 0 4.0554.9823 4.0554 2.1941z",
}, {
	faviconColorsIndex: 6,
	d:                  "m19.977 15.338c0 .5685-.43366.8554-1.1381 1.0001-.29193.06-.63037.096-1.0037.1166-.56405.032-1.2078.031-1.8912.031-.67283 0-1.3072 0-1.8649-.029-.30627-.017-.58943-.043-.84316-.084-.81383-.1318-1.325-.417-1.325-1.0344 0-1.1601 1.8056-2.1006 4.033-2.1006s4.033.9405 4.033 2.1006z",
}, {
	faviconColorsIndex: 7,
	d:                  "m18.025 13.488a2.0802 1.3437 0 0 1 -2.0802 1.3437 2.0802 1.3437 0 0 1 -2.0802 -1.3437 2.0802 1.3437 0 0 1 2.0802 -1.3437 2.0802 1.3437 0 0 1 2.0802 1.3437z",
}}

func TestEncodeFavicon(t *testing.T) {
	var e encode.Encoder
	var gen Generator
	gen.SetDestination(&e)

	// Set up a base color for theming the favicon, gopher blue by default.
	pal := ivg.DefaultPalette
	pal[0] = faviconColors[0] // color.RGBA{0x76, 0xe1, 0xfe, 0xff}

	gen.Reset(ivg.DefaultViewBox, pal)

	// The favicon graphic also uses a dark version of that base color. blend
	// is 75% dark (CReg[63]) and 25% the base color (pal[0]).
	dark := color.RGBA{0x23, 0x1d, 0x1b, 0xff}
	blend := ivg.BlendColor(0x40, 0xff, 0x80)

	// First, set CReg[63] to dark, then set CReg[63] to the blend of that dark
	// color with pal[0].
	gen.SetCReg(1, false, ivg.RGBAColor(dark))
	gen.SetCReg(1, false, blend)

	// Check that, for the suggested palette, blend resolves to the
	// (non-themable) SVG file's faviconColors[1].
	got := blend.Resolve(&pal, &[64]color.RGBA{
		63: dark,
	})
	want := faviconColors[1]
	if got != want {
		t.Fatalf("Blend:\ngot  %#02x\nwant %#02x", got, want)
	}

	// Set aside the remaining, non-themable colors.
	remainingColors := faviconColors[2:]

	seenFCI2 := false
	for _, data := range faviconSVGData {
		adj := uint8(data.faviconColorsIndex)
		if adj >= 2 {
			if !seenFCI2 {
				seenFCI2 = true
				for i, c := range remainingColors {
					gen.SetCReg(uint8(i), false, ivg.RGBAColor(c))
				}
			}
			adj -= 2
		}

		fromVB := Concat(Translate(0, 0), Scale(1.0/32.0, 1.0/32.0))
		toVB := Concat(Scale(64, 64), Translate(-32, -32))
		gen.SetTransform(fromVB, toVB)
		if err := gen.SetPathData(data.d, adj); err != nil {
			t.Fatal(err)
		}
	}

	testEncode(t, &e, "../testdata/favicon.ivg")
}
