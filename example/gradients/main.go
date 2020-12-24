// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/icon"
)

func main() {
	go Gradients()
	app.Main()
}

func Gradients() {
	window := app.NewWindow(
		app.Title("IVG - Gradients"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	rasterizer := icon.Rasterizer(icon.GioRasterizer)
	Grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
	Grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
	gradients := GradientsImage{}
	ops := new(op.Ops)
	backdrop := new(int)
	for next := range window.Events() {
		if frame, ok := next.(system.FrameEvent); ok {
			ops.Reset()

			// clicking on backdrop will switch active renderer
			pointer.InputOp{Tag: backdrop, Types: pointer.Release}.Add(ops)
			for _, next := range frame.Queue.Events(backdrop) {
				if event, ok := next.(pointer.Event); ok {
					if event.Type == pointer.Release {
						switch rasterizer {
						case icon.GioRasterizer:
							rasterizer = icon.VecRasterizer
						case icon.VecRasterizer:
							rasterizer = icon.GioRasterizer
						}
					}
				}
			}

			// fill the whole backdrop rectangle
			paint.ColorOp{Color: Grey800}.Add(ops)
			paint.PaintOp{}.Add(ops)

			// device independent content rect calculation
			margin := unit.Dp(12)
			minX := unit.Add(frame.Metric, margin, frame.Insets.Left)
			minY := unit.Add(frame.Metric, margin, frame.Insets.Top)
			maxX := unit.Add(frame.Metric, unit.Px(float32(frame.Size.X)), frame.Insets.Right.Scale(-1), margin.Scale(-1))
			maxY := unit.Add(frame.Metric, unit.Px(float32(frame.Size.Y)), frame.Insets.Bottom.Scale(-1), margin.Scale(-1))
			contentRect := f32.Rect(
				float32(frame.Metric.Px(minX)), float32(frame.Metric.Px(minY)),
				float32(frame.Metric.Px(maxX)), float32(frame.Metric.Px(maxY)))

			// fill content rect
			stack := op.Push(ops)
			paint.ColorOp{Color: Grey300}.Add(ops)
			op.Offset(contentRect.Min).Add(ops)
			clip.Rect(image.Rect(0, 0, int(contentRect.Dx()), int(contentRect.Dy()))).Add(ops)
			paint.PaintOp{}.Add(ops)
			stack.Pop()

			// scale the viewbox of the icon to the content rect
			viewRect := gradients.AspectMeet(contentRect, 0.5, 0.5)

			// render actual content
			start := time.Now()
			if callOp, err := rasterizer.Rasterize(gradients, viewRect); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			msg := fmt.Sprintf("%s (%v)", rasterizer.Name(), time.Since(start).Round(time.Microsecond))
			PrintText(msg, contentRect.Min, 0.0, 0.0, contentRect.Dx(), H5, ops)

			frame.Frame(ops)
		}
	}
	os.Exit(0)
}

type GradientsImage struct{}

var GradientsViewBox = ivg.ViewBox{
	MinX: -32, MinY: -32,
	MaxX: +32, MaxY: +32,
}

func (g GradientsImage) AspectMeet(rect f32.Rectangle, ax, ay float32) f32.Rectangle {
	return f32.Rect(GradientsViewBox.AspectMeet(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y, ax, ay))
}

func (g GradientsImage) RenderOn(dst ivg.Destination, col ...color.RGBA) error {

	// Uncomment lines below to log rasterizer calls.
	// logger := &raster.RasterizerLogger{Rasterizer: rasterizer}
	// renderer.SetRasterizer(logger, bounds)

	gen := generate.Generator{}
	gen.SetDestination(dst)

	// Uncomment lines below to log destination calls.
	// logger := &ivg.DestinationLogger{Destination: renderer}
	// gen.SetDestination(logger)

	gen.Reset(GradientsViewBox, &ivg.DefaultPalette)

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

	return nil
}
