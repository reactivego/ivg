// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/raster/gio"
	"github.com/reactivego/ivg/raster/vec"
	"github.com/reactivego/ivg/render"
)

func main() {
	go Favicon()
	app.Main()
}

func Favicon() {
	window := app.NewWindow(
		app.Title("IVG - Favicon"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	const (
		Gio = iota
		Vec
	)
	rasterizer := Gio
	ops := new(op.Ops)
	backdrop := new(int)
	for next := range window.Events() {
		if frame, ok := next.(system.FrameEvent); ok {
			ops.Reset()

			// initial window rect in pixels
			rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))

			// backdrop switch renderer on release and fill rectangle
			pointer.InputOp{Tag: backdrop, Types: pointer.Release}.Add(ops)
			for _, next := range frame.Queue.Events(backdrop) {
				if event, ok := next.(pointer.Event); ok {
					if event.Type == pointer.Release {
						switch rasterizer {
						case Gio:
							rasterizer = Vec
						case Vec:
							rasterizer = Gio
						}
					}
				}
			}
			paint.ColorOp{Color: colornames.Grey800}.Add(ops)
			paint.PaintOp{Rect: rect}.Add(ops)

			// device independent content rect calculation
			pt32 := func(x, y unit.Value) f32.Point {
				return f32.Pt(float32(frame.Metric.Px(x)), float32(frame.Metric.Px(y)))
			}
			margin := pt32(unit.Dp(12), unit.Dp(12))
			lefttop := pt32(frame.Insets.Left, frame.Insets.Top).Add(margin)
			rightbottom := pt32(frame.Insets.Right, frame.Insets.Bottom).Add(margin)
			rect = f32.Rectangle{Min: rect.Min.Add(lefttop), Max: rect.Max.Sub(rightbottom)}

			// fill content rect
			op.Offset(rect.Min).Add(ops)
			rect = f32.Rectangle{Max: rect.Size()}
			paint.ColorOp{Color: colornames.Grey300}.Add(ops)
			paint.PaintOp{Rect: rect}.Add(ops)

			// render actual content
			viewrect := ViewBox.SizeToRect(ivg.Rect(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y), ivg.AspectMeet, ivg.Mid, ivg.Mid)
			bounds := image.Rect(viewrect.IntFields())
			rect = f32.Rect(viewrect.Fields())
			renderer := &render.Renderer{}
			switch rasterizer {
			case Gio:
				start := time.Now()
				renderer.SetRasterizer(&gio.Rasterizer{Ops: ops}, bounds)
				Render(renderer, ViewBox)
				PrintText(fmt.Sprintf("Gio (%v)", time.Since(start).Round(time.Microsecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
			case Vec:
				start := time.Now()
				dst := image.NewRGBA(bounds)
				renderer.SetRasterizer(&vec.Rasterizer{Dst: dst, DrawOp: draw.Src}, bounds)
				Render(renderer, ViewBox)
				paint.NewImageOp(dst).Add(ops)
				paint.PaintOp{Rect: rect}.Add(ops)
				PrintText(fmt.Sprintf("Vec (%v)", time.Since(start).Round(time.Millisecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
			}

			frame.Frame(ops)
		}
	}
	os.Exit(0)
}

var ViewBox = ivg.ViewBox{
	MinX: 0, MinY: 0,
	MaxX: +48, MaxY: +48,
}

func Render(renderer ivg.Destination, viewbox ivg.ViewBox) {

	// Uncomment lines below to log rasterizer calls.
	// logger := &raster.RasterizerLogger{Rasterizer: rasterizer}
	// renderer.SetRasterizer(logger, bounds)

	gen := generate.Generator{}
	gen.SetDestination(renderer)

	// Uncomment lines below to log destination calls.
	// logger := &ivg.DestinationLogger{Destination: renderer}
	// gen.SetDestination(logger)

	gen.Reset(viewbox, &ivg.DefaultPalette)

	colors := []color.RGBA{
		{0x76, 0xe1, 0xfe, 0xff}, // 0
		{0x38, 0x4e, 0x54, 0xff}, // 1
		{0xff, 0xff, 0xff, 0xff}, // 2
		{0x17, 0x13, 0x11, 0xff}, // 3
		{0x00, 0x00, 0x00, 0x54}, // 4
		{0xff, 0xfc, 0xfb, 0xff}, // 5
		{0xc3, 0x8c, 0x74, 0xff}, // 6
		{0x23, 0x20, 0x1f, 0xff}, // 7
	}

	type Path struct {
		i int
		d string
	}

	paths := []Path{{
		i: 1,
		d: "m16.092 1.002c-1.1057.01-2.2107.048844-3.3164.089844-2.3441.086758-4.511.88464-6.2832 2.1758a3.8208 3.5794 29.452 0 0 -.8947 -.6856 3.8208 3.5794 29.452 0 0 -5.0879 1.2383 3.8208 3.5794 29.452 0 0 1.5664 4.9961 3.8208 3.5794 29.452 0 0 .3593 .1758c-.2784.9536-.4355 1.9598-.4355 3.0078v20h28v-20c0-1.042-.152-2.0368-.418-2.9766a3.5794 3.8208 60.548 0 0 .43359 -.20703 3.5794 3.8208 60.548 0 0 1.5684 -4.9961 3.5794 3.8208 60.548 0 0 -5.0879 -1.2383 3.5794 3.8208 60.548 0 0 -.92969 .72461c-1.727-1.257-3.843-2.0521-6.1562-2.2148-1.1058-.078-2.2126-.098844-3.3184-.089844z",
	}, {
		i: 0,
		d: "m16 3c-4.835 0-7.9248 1.0791-9.7617 2.8906-.4777-.4599-1.2937-1.0166-1.6309-1.207-.9775-.5520-2.1879-.2576-2.7051.6582-.5171.9158-.1455 2.1063.8321 2.6582.2658.1501 1.2241.5845 1.7519.7441-.3281.9946-.4863 2.0829-.4863 3.2559v20h24c-.049-7.356 0-18 0-20 0-1.209-.166-2.3308-.516-3.3496.539-.2011 1.243-.5260 1.463-.6504.978-.5519 1.351-1.7424.834-2.6582s-1.729-1.2102-2.707-.6582c-.303.1711-.978.6356-1.463 1.0625-1.854-1.724-4.906-2.7461-9.611-2.7461z",
	}, {
		i: 1,
		d: "m3.0918 5.9219c-.060217.00947-.10772.020635-.14648.033203-.019384.00628-.035462.013581-.052734.021484-.00864.00395-.019118.00825-.03125.015625-.00607.00369-.011621.00781-.021484.015625-.00493.00391-.017342.015389-.017578.015625-.0002366.0002356-.025256.031048-.025391.03125a.19867 .19867 0 0 0 .26367 .28320c.0005595-.0002168.00207-.00128.00391-.00195a.19867 .19867 0 0 0 .00391 -.00195c.015939-.00517.045148-.013113.085937-.019531.081581-.012836.20657-.020179.36719.00391.1020.0152.2237.0503.3535.0976-.3277.0694-.5656.1862-.7227.3145-.1143.0933-.1881.1903-.2343.2695-.023099.0396-.039499.074216-.050781.10547-.00564.015626-.00989.029721-.013672.046875-.00189.00858-.00458.017085-.00586.03125-.0006392.00708-.0005029.014724 0 .027344.0002516.00631.00192.023197.00195.023437.0000373.0002412.0097.036937.00977.037109a.19867 .19867 0 0 0 .38477 -.039063 .19867 .19867 0 0 0 0 -.00195c.00312-.00751.00865-.015947.017578-.03125.0230-.0395.0660-.0977.1425-.1601.1530-.1250.4406-.2702.9863-.2871a.19930 .19930 0 0 0 .082031 -.019531c.12649.089206.25979.19587.39844.32422a.19867 .19867 0 1 0 .2696 -.2911c-.6099-.5646-1.1566-.7793-1.5605-.8398-.2020-.0303-.3679-.0229-.4883-.0039z",
	}, {
		i: 1,
		d: "m28.543 5.8203c-.12043-.018949-.28631-.026379-.48828.00391-.40394.060562-.94869.27524-1.5586.83984a.19867 .19867 0 1 0 .26953 .29102c.21354-.19768.40814-.33222.59180-.44141.51624.023399.79659.16181.94531.28320.07652.062461.11952.12063.14258.16016.0094.016037.01458.025855.01758.033203a.19867 .19867 0 0 0 .38476 .039063c.000062-.0001719.0097-.036868.0098-.037109.000037-.0002412.0017-.017125.002-.023437.000505-.012624.000639-.020258 0-.027344-.0013-.01417-.004-.022671-.0059-.03125-.0038-.017158-.008-.031248-.01367-.046875-.01128-.031254-.02768-.067825-.05078-.10742-.04624-.079195-.12003-.17424-.23437-.26758-.11891-.097066-.28260-.18832-.49609-.25781.01785-.00328.03961-.011119.05664-.013672.16062-.024082.28561-.016738.36719-.00391.03883.00611.06556.012409.08203.017578.000833.0002613.0031.0017.0039.00195a.19867 .19867 0 0 0 .271 -.2793c-.000135-.0002016-.02515-.031014-.02539-.03125-.000236-.0002356-.01265-.011717-.01758-.015625-.0099-.00782-.01737-.01194-.02344-.015625-.01213-.00737-.02066-.011673-.0293-.015625-.01727-.0079-.03336-.013247-.05273-.019531-.03877-.012568-.08822-.025682-.14844-.035156z",
	}, {
		i: 2,
		d: "m15.171 9.992a4.8316 4.8316 0 0 1 -4.832 4.832 4.8316 4.8316 0 0 1 -4.8311 -4.832 4.8316 4.8316 0 0 1 4.8311 -4.8316 4.8316 4.8316 0 0 1 4.832 4.8316z",
	}, {
		i: 2,
		d: "m25.829 9.992a4.6538 4.6538 0 0 1 -4.653 4.654 4.6538 4.6538 0 0 1 -4.654 -4.654 4.6538 4.6538 0 0 1 4.654 -4.6537 4.6538 4.6538 0 0 1 4.653 4.6537z",
	}, {
		i: 3,
		d: "m14.377 9.992a1.9631 1.9631 0 0 1 -1.963 1.963 1.9631 1.9631 0 0 1 -1.963 -1.963 1.9631 1.9631 0 0 1 1.963 -1.963 1.9631 1.9631 0 0 1 1.963 1.963z",
	}, {
		i: 3,
		d: "m25.073 9.992a1.9631 1.9631 0 0 1 -1.963 1.963 1.9631 1.9631 0 0 1 -1.963 -1.963 1.9631 1.9631 0 0 1 1.963 -1.963 1.9631 1.9631 0 0 1 1.963 1.963z",
	}, {
		i: 4,
		d: "m14.842 15.555h2.2156c.40215 0 .72590.3237.72590.7259v2.6545c0 .4021-.32375.7259-.72590.7259h-2.2156c-.40215 0-.72590-.3238-.72590-.7259v-2.6545c0-.4022.32375-.7259.72590-.7259z",
	}, {
		i: 5,
		d: "m14.842 14.863h2.2156c.40215 0 .72590.3238.72590.7259v2.6546c0 .4021-.32375.7259-.72590.7259h-2.2156c-.40215 0-.72590-.3238-.72590-.7259v-2.6546c0-.4021.32375-.7259.72590-.7259z",
	}, {
		i: 4,
		d: "m20 16.167c0 .838-.87123 1.2682-2.1448 1.1659-.02366 0-.04795-.6004-.25415-.5832-.50367.042-1.0959-.02-1.686-.02-.61294 0-1.2063.1826-1.6855.017-.11023-.038-.17830.5838-.26153.5816-1.2437-.033-2.0788-.3383-2.0788-1.1618 0-1.2118 1.8156-2.1941 4.0554-2.1941 2.2397 0 4.0554.9823 4.0554 2.1941z",
	}, {
		i: 6,
		d: "m19.977 15.338c0 .5685-.43366.8554-1.1381 1.0001-.29193.06-.63037.096-1.0037.1166-.56405.032-1.2078.031-1.8912.031-.67283 0-1.3072 0-1.8649-.029-.30627-.017-.58943-.043-.84316-.084-.81383-.1318-1.325-.417-1.325-1.0344 0-1.1601 1.8056-2.1006 4.033-2.1006s4.033.9405 4.033 2.1006z",
	}, {
		i: 7,
		d: "m18.025 13.488a2.0802 1.3437 0 0 1 -2.0802 1.3437 2.0802 1.3437 0 0 1 -2.0802 -1.3437 2.0802 1.3437 0 0 1 2.0802 -1.3437 2.0802 1.3437 0 0 1 2.0802 1.3437z",
	}}

	// Set up a base color for theming the favicon, gopher blue by default.
	pal := ivg.DefaultPalette
	pal[0] = colors[0] // color.RGBA{0x76, 0xe1, 0xfe, 0xff}

	gen.Reset(ivg.DefaultViewBox, &pal)

	// The favicon graphic also uses a dark version of that base color. blend
	// is 75% dark (CReg[63]) and 25% the base color (pal[0]).
	dark := color.RGBA{0x23, 0x1d, 0x1b, 0xff}
	blend := ivg.BlendColor(0x40, 0xff, 0x80)

	// First, set CReg[63] to dark, then set CReg[63] to the blend of that dark
	// color with pal[0].
	gen.SetCReg(1, false, ivg.RGBAColor(dark))
	gen.SetCReg(1, false, blend)

	// Set aside the remaining, non-themable colors.
	remainingColors := colors[2:]

	seenFCI2 := false
	for _, path := range paths {
		adj := uint8(path.i)
		if adj >= 2 {
			if !seenFCI2 {
				seenFCI2 = true
				for i, c := range remainingColors {
					gen.SetCReg(uint8(i), false, ivg.RGBAColor(c))
				}
			}
			adj -= 2
		}
		gen.SetPathData(path.d, adj, true)
	}
}
