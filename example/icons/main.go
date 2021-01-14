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
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
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
	var rasterizer icon.Rasterizer = icon.GioRasterizer
	Grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
	Grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
	ops := new(op.Ops)
	backdrop := new(int)
	index := 0
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
			paint.ColorOp{Color: Grey300}.Add(ops)
			state := op.Save(ops)
			op.Offset(contentRect.Min).Add(ops)
			clip.Rect(image.Rect(0, 0, int(contentRect.Dx()), int(contentRect.Dy()))).Add(ops)
			paint.PaintOp{}.Add(ops)
			state.Load()

			// select next icon and paint
			n := uint(len(IconCollection))
			ico := IconCollection[(uint(index)+n)%n]
			index++
			start := time.Now()
			icon, err := icon.New(ico.data)
			if err != nil {
				log.Fatal(err)
			}
			viewRect := icon.AspectMeet(contentRect, 0.5, 0.5)
			if callOp, err := rasterizer.Rasterize(icon, viewRect, colornames.LightBlue600); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			msg := fmt.Sprintf("%s (%v)", rasterizer.Name(), time.Since(start).Round(time.Microsecond))
			PrintText(msg, contentRect.Min, 0.0, 0.0, contentRect.Dx(), H5, ops)

			at := time.Now().Add(500 * time.Millisecond)
			op.InvalidateOp{At: at}.Add(ops)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
