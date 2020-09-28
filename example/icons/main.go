// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"log"
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

	"github.com/reactivego/ivg/icon"
)

func main() {
	go Icons()
	app.Main()
}

func Icons() {
	window := app.NewWindow(
		app.Title("IVG - Icons"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	rasterizer := icon.GioRasterizer
	ops := new(op.Ops)
	backdrop := new(int)
	index := 0
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
						case icon.GioRasterizer:
							rasterizer = icon.VecRasterizer
						case icon.VecRasterizer:
							rasterizer = icon.GioRasterizer
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

			// select next icon and paint
			n := uint(len(IconCollection))
			ico := IconCollection[(uint(index)+n)%n]
			index++
			start := time.Now()
			if callOp, err := icon.FromData(ico.data, colornames.LightBlue600, rect, icon.AspectMeet, icon.Mid, icon.Mid, rasterizer); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			switch rasterizer {
			case icon.GioRasterizer:
				PrintText(fmt.Sprintf("Gio (%v)", time.Since(start).Round(time.Microsecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
			case icon.VecRasterizer:
				PrintText(fmt.Sprintf("Vec (%v)", time.Since(start).Round(time.Millisecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
			}

			at := time.Now().Add(500 * time.Millisecond)
			op.InvalidateOp{At:at}.Add(ops)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
