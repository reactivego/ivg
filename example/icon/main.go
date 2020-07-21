package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"golang.org/x/exp/shiny/materialdesign/colornames"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/reactivego/ivg/icon"
)

const (
	Title    = "IVG - Icon"
	WidthDp  = 768
	HeightDp = 768
	MarginDp = 12
)

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

var (
	ops      = new(op.Ops)
	backdrop Backdrop
	index    = 0
	raster   icon.Rasterizer
)

func HandleFrameEvent(event system.FrameEvent) {
	ops.Reset()

	constraints := f32.Rect(0, 0, float32(event.Size.X), float32(event.Size.Y))

	backdrop.Color = colornames.Grey800
	if backdrop.Press(constraints, event.Queue, ops) {
		switch raster {
		case icon.GioRasterizer:
			raster = icon.VectorRasterizer
		case icon.VectorRasterizer:
			raster = icon.GioRasterizer
		}
	}
	backdrop.Paint(constraints, ops)

	margin := unit.Dp(MarginDp)
	leftInset := unit.Add(event.Metric, event.Insets.Left, margin)
	topInset := unit.Add(event.Metric, event.Insets.Top, margin)
	rightInset := unit.Add(event.Metric, event.Insets.Right, margin)
	bottomInset := unit.Add(event.Metric, event.Insets.Bottom, margin)
	constraints.Min.X += float32(event.Metric.Px(leftInset))
	constraints.Min.Y += float32(event.Metric.Px(topInset))
	constraints.Max.X -= float32(event.Metric.Px(rightInset))
	constraints.Max.Y -= float32(event.Metric.Px(bottomInset))

	op.Offset(constraints.Min).Add(ops)
	constraints = constraints.Sub(constraints.Min)

	paint.ColorOp{Color: colornames.Grey300}.Add(ops)
	paint.PaintOp{Rect: constraints}.Add(ops)

	n := uint(len(Icons))
	ico := Icons[(uint(index)+n)%n]
	index++

	if callOp, err := icon.FromData(ico.data, colornames.LightBlue600, constraints, icon.AspectMeet, icon.Mid, icon.Mid, raster); err == nil {
		callOp.Add(ops)
	}

	switch raster {
	case icon.GioRasterizer:
		PrintText("Gio", constraints.Min, 0.0, 0.0, 1000, H3, ops)
	case icon.VectorRasterizer:
		PrintText("Vector", constraints.Min, 0.0, 0.0, 1000, H3, ops)
	}
	PrintText(ico.name, f32.Pt(constraints.Min.X, constraints.Max.Y), 0.0, 1.0, 1000, BodyText1, ops)

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
