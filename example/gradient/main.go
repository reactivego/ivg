// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/raster"
	"github.com/reactivego/ivg/raster/gio"
	"github.com/reactivego/ivg/raster/vector"
	"github.com/reactivego/ivg/render"
)

const (
	Title    = "IVG - Gradient"
	WidthDp  = 768
	HeightDp = 768
	MarginDp = 12
)

const (
	Gio = iota
	Vector
)

var (
	SelectedRasterizer = Gio

	ops = new(op.Ops)
)

const (
	// AspectNone stretches or squashes the ViewBox to meet the contraints rect.
	AspectNone = iota
	// AspectMeet fits the ViewBox inside the constraints rect maintaining its
	// aspect ratio.
	AspectMeet
	// AspectSlice fills the constraints rect maintaining the ViewBox's aspect
	// ratio.
	ASpectSlice
)

// PreserveAspectRatio determines how the ViewBox is positioned in the
// constraints rectangle. We always use xMidYMid to position the viewbox in
// the center of the constraints rect.
const PreserveAspectRatio = AspectMeet

func MultipleGradients(constraints f32.Rectangle, ops *op.Ops) {
	viewbox := ivg.ViewBox{
		MinX: -32, MinY: -32,
		MaxX: +32, MaxY: +32,
	}
	dx, dy := constraints.Dx(), constraints.Dy()
	vbdx, vbdy := viewbox.AspectRatio()
	vbAR := vbdx / vbdy
	switch PreserveAspectRatio {
	case AspectMeet:
		if dx/dy < vbAR {
			dy = dx / vbAR
		} else {
			dx = dy * vbAR
		}
	case ASpectSlice:
		if dx/dy < vbAR {
			dx = dy * vbAR
		} else {
			dy = dx / vbAR
		}
	}
	midX := (constraints.Min.X + constraints.Max.X) / 2
	midY := (constraints.Min.Y + constraints.Max.Y) / 2
	rect := f32.Rect(midX-dx/2, midY-dy/2, midX+dx/2, midY+dy/2)

	bounds := image.Rect(int(rect.Min.X), int(rect.Min.Y), int(rect.Max.X), int(rect.Max.Y))

	var rasterizer raster.Rasterizer
	var dst *image.RGBA
	switch SelectedRasterizer {
	case Gio:
		rasterizer = &gio.Rasterizer{Ops: ops}
	case Vector:
		dst = image.NewRGBA(bounds)
		rasterizer = &vector.Rasterizer{Dst: dst, DrawOp: draw.Src}
	}

	renderer := &render.Renderer{}
	renderer.SetRasterizer(rasterizer, bounds)

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

	if SelectedRasterizer == Vector {
		paint.NewImageOp(dst).Add(ops)
		paint.PaintOp{Rect: rect}.Add(ops)
	}
}

type Backdrop struct{ Color color.RGBA }

func (backdrop *Backdrop) Paint(constraints f32.Rectangle, ops *op.Ops) {
	paint.ColorOp{Color: backdrop.Color}.Add(ops)
	paint.PaintOp{Rect: constraints}.Add(ops)
}

func (backdrop *Backdrop) Press(constraints f32.Rectangle, queue event.Queue, ops *op.Ops) bool {
	stack := op.Push(ops)
	defer stack.Pop()
	rect := image.Rect(int(constraints.Min.X), int(constraints.Min.Y), int(constraints.Max.X), int(constraints.Max.Y))
	pointer.Rect(rect).Add(ops)
	pointer.InputOp{Tag: backdrop, Types: pointer.Press}.Add(ops)
	for _, next := range queue.Events(backdrop) {
		if event, ok := next.(pointer.Event); ok {
			if event.Type == pointer.Press {
				return true
			}
		}
	}
	return false
}

var backdrop Backdrop

func HandleFrameEvent(event system.FrameEvent) {
	ops.Reset()

	// initial contraints in pixels
	minX := float32(0)
	minY := float32(0)
	maxX := float32(event.Size.X)
	maxY := float32(event.Size.Y)
	constraints := f32.Rect(minX, minY, maxX, maxY)

	// fill backdrop and check for button press
	backdrop.Color = colornames.Grey800
	if backdrop.Press(constraints, event.Queue, ops) {
		switch SelectedRasterizer {
		case Gio:
			SelectedRasterizer = Vector
		case Vector:
			SelectedRasterizer = Gio
		}
		// Flash backdrop color
		backdrop.Color = colornames.Grey400
	}
	backdrop.Paint(constraints, ops)

	// device independent inset + margin calculation
	margin := unit.Dp(MarginDp)
	leftInset := unit.Add(event.Metric, event.Insets.Left, margin)
	topInset := unit.Add(event.Metric, event.Insets.Top, margin)
	rightInset := unit.Add(event.Metric, event.Insets.Right, margin)
	bottomInset := unit.Add(event.Metric, event.Insets.Bottom, margin)

	// apply insets + margins to pixel constraints
	minX += float32(event.Metric.Px(leftInset))
	minY += float32(event.Metric.Px(topInset))
	maxX -= float32(event.Metric.Px(rightInset))
	maxY -= float32(event.Metric.Px(bottomInset))

	constraints = f32.Rect(minX, minY, maxX, maxY)
	op.Offset(constraints.Min).Add(ops)
	constraints = f32.Rect(0, 0, constraints.Dx(), constraints.Dy())
	paint.ColorOp{Color: colornames.Grey300}.Add(ops)
	paint.PaintOp{Rect: constraints}.Add(ops)

	MultipleGradients(constraints, ops)

	switch SelectedRasterizer {
	case Gio:
		PrintText("Gio", constraints.Min, 0.0, 0.0, 1000, H3, ops)
	case Vector:
		PrintText("Vector", constraints.Min, 0.0, 0.0, 1000, H3, ops)
	}

	event.Frame(ops)
}

func observe(next event.Event, err error, done bool) {
	switch {
	case !done:
		if event, ok := next.(system.FrameEvent); ok {
			HandleFrameEvent(event)
		}
	case err != nil:
		fmt.Printf("error %+v\n", err)
	}
}

func main() {
	window := app.NewWindow(
		app.Title(Title),
		app.Size(unit.Dp(WidthDp), unit.Dp(HeightDp)),
	)
	go func() {
		var err error
		for next := range window.Events() {
			if e, ok := next.(system.DestroyEvent); ok {
				err = e.Err
				break
			} else {
				observe(next, nil, false)
			}
		}
		observe(nil, err, true)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}
