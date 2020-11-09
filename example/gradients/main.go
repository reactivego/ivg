// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/raster/gio"
	"github.com/reactivego/ivg/raster/vec"
	"github.com/reactivego/ivg/render"
)

func main() {
	go Gradient()
	app.Main()
}

func Gradient() {
	window := app.NewWindow(
		app.Title("IVG - Gradients"),
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
	MinX: -32, MinY: -32,
	MaxX: +32, MaxY: +32,
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
}
