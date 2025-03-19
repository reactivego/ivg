// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"

	"github.com/reactivego/gio"
	"github.com/reactivego/gio/style"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
	"github.com/reactivego/ivg/generate"
	raster "github.com/reactivego/ivg/raster/gio"
)

func main() {
	go Gradients()
	app.Main()
}

func Gradients() {
	window := app.NewWindow(
		app.Title("IVG - Gradients"),
		app.Size(768, 768),
	)

	grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
	grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
	black := color.NRGBA{A: 255}

	data, err := GradientsIVG()
	if err != nil {
		log.Fatal(err)
	}

	ops := new(op.Ops)
	shaper := text.NewShaper(style.FontFaces())
	backend := "Gio"
	for next := range window.Events() {
		if event, ok := next.(system.FrameEvent); ok {
			gtx := layout.NewContext(ops, event)

			// backdrop
			pointer.InputOp{Tag: backend, Types: pointer.Release}.Add(gtx.Ops)
			for _, next := range event.Queue.Events(backend) {
				if event, ok := next.(pointer.Event); ok {
					if event.Type == pointer.Release {
						backend = map[string]string{"Gio": "Img", "Img": "Gio"}[backend]
					}
				}
			}
			paint.Fill(gtx.Ops, grey800)

			layout.UniformInset(12).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				size := gtx.Constraints.Max
				paint.FillShape(ops, grey300, clip.Rect(image.Rectangle{Max: size}).Op())
				start := time.Now()
				var icon layout.Widget
				switch backend {
				case "Gio":
					icon, _ = raster.Icon(data, 48, 48)
				case "Img":
					icon, _ = raster.Icon(data, 48, 48, raster.WithImageBackend())
				}
				icon(gtx)
				msg := fmt.Sprintf("%s (%v)", backend, time.Since(start).Round(time.Microsecond))
				text := gio.Text(shaper, style.H5, 0.0, 0.0, black, msg)
				text(gtx)
				return layout.Dimensions{Size: size}
			})

			event.Frame(ops)
		}
	}
	os.Exit(0)
}

func GradientsIVG() ([]byte, error) {
	enc := &encode.Encoder{}
	gen := generate.Generator{}
	gen.SetDestination(enc)

	viewbox := ivg.ViewBox{
		MinX: -32, MinY: -32,
		MaxX: +32, MaxY: +32}
	gen.Reset(viewbox, ivg.DefaultPalette)

	rgb := []generate.GradientStop{
		{Offset: 0.00, Color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
		{Offset: 1.00, Color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}
	cmy := []generate.GradientStop{
		{Offset: 0.00, Color: color.RGBA{0x00, 0xff, 0xff, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0xff, 0x00, 0xff, 0xff}},
		{Offset: 0.75, Color: color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{Offset: 1.00, Color: color.RGBA{0xff, 0xff, 0x00, 0xff}},
	}

	x1, y1 := float32(-12), float32(-30)
	x2, y2 := float32(+12), float32(-18)
	minX, minY := float32(-30), float32(-30)
	maxX, maxY := float32(+30), float32(-18)

	gen.SetLinearGradient(x1, y1, x2, y2, generate.GradientSpreadNone, rgb)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	x1, y1 = -12, -14
	x2, y2 = +12, -2
	minY = -14
	maxY = -2

	gen.SetLinearGradient(x1, y1, x2, y2, generate.GradientSpreadPad, cmy)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	cx, cy := float32(-8), float32(+8)
	rx, ry := float32(0), float32(+16)
	minY = +2
	maxY = +14

	gen.SetCircularGradient(cx, cy, rx, ry, generate.GradientSpreadReflect, rgb)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	cx, cy = -8, +24
	rx, ry = 0, 16
	minY = +18
	maxY = +30

	gen.SetCircularGradient(cx, cy, rx, ry, generate.GradientSpreadRepeat, cmy)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	return enc.Bytes()
}
