// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/vector"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/reactivego/ivg/raster/gio"
)

func main() {
	go Blend()
	app.Main()
}

func Blend() {
	window := app.NewWindow(
		app.Title("IVG - Blend"),
		app.Size(unit.Dp(480), unit.Dp(480)),
	)
	ops := new(op.Ops)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()

			dx, dy := float32(frame.Size.X), float32(frame.Size.Y)

			yellow := color.RGBA{0xfc, 0xe9, 0x4f, 0xff}
			highlight := color.NRGBA{0xfd, 0xee, 0x74, 0x7f}

			// Use special gio color model to pre-multiply color for correctly blending
			// highlight color over opaque yellow background color using gio.
			RGBA := func(c color.Color) color.RGBA {
				return gio.RGBAModel.Convert(c).(color.RGBA)
			}
			upper := f32.Rect(0, 0, dx, dy/2)
			paint.ColorOp{Color: RGBA(yellow)}.Add(ops)
			paint.PaintOp{Rect: upper}.Add(ops)
			paint.ColorOp{Color: RGBA(highlight)}.Add(ops)
			paint.PaintOp{Rect: upper}.Add(ops)

			// Using image/vector rasterizer to blend highlight color over opaque
			// yellow background color.
			RGBA = func(c color.Color) color.RGBA {
				return color.RGBAModel.Convert(c).(color.RGBA)
			}
			lower := f32.Rect(0, dy/2, dx, dy)
			z := vector.NewRasterizer(int(dx), int(dy/2))
			z.MoveTo(0, 0)
			z.LineTo(dx, 0)
			z.LineTo(dx, dy/2)
			z.LineTo(0, dy/2)
			z.ClosePath()
			dst := image.NewRGBA(z.Bounds())
			src := image.NewUniform(RGBA(yellow))
			z.DrawOp = draw.Src
			z.Draw(dst, dst.Bounds(), src, src.Bounds().Min)
			src = image.NewUniform(RGBA(highlight))
			z.DrawOp = draw.Over
			z.Draw(dst, dst.Bounds(), src, src.Bounds().Min)
			paint.NewImageOp(dst).Add(ops)
			paint.PaintOp{Rect: lower}.Add(ops)

			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
