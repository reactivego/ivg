// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/vector"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func main() {
	go Blend()
	app.Main()
}

func Blend() {
	window := app.NewWindow(
		app.Title("IVG - Blend"),
		app.Size(480, 480),
	)
	ops := new(op.Ops)
	for next := range window.Events() {
		if event, ok := next.(system.FrameEvent); ok {
			ops.Reset()
			dx, dy := event.Size.X, event.Size.Y

			// This app demonstrates blending of a translucent highlight color on top of an opaque
			// background color. The upper half of the window is filled using "gioui.org/op/paint"
			// operations, while the lower half of the window is filled using "golang.org/x/image/vector".
			// Both halves should look identical.

			yellow := color.NRGBA{0xfc, 0xe9, 0x4f, 0xff}
			highlight := color.NRGBA{0xfd, 0xee, 0x74, 0x7f}

			// Upper half of the window is painted using "gioui.org/op/paint"
			// to blend a translucent highlight color over an opaque yellow
			// background color.
			paint.ColorOp{Color: yellow}.Add(ops)
			cstack := clip.Rect(image.Rect(0, 0, dx, dy/2)).Push(ops)
			paint.PaintOp{}.Add(ops)
			paint.ColorOp{Color: highlight}.Add(ops)
			paint.PaintOp{}.Add(ops)
			cstack.Pop()

			// Lower half of the window is painted using
			// "golang.org/x/image/vector".
			RGBA := func(c color.Color) color.RGBA {
				return color.RGBAModel.Convert(c).(color.RGBA)
			}
			z := vector.NewRasterizer(dx, dy/2)
			z.MoveTo(0, 0)
			z.LineTo(float32(dx), 0)
			z.LineTo(float32(dx), float32(dy/2))
			z.LineTo(0, float32(dy/2))
			z.ClosePath()
			dst := image.NewRGBA(z.Bounds())
			src := image.NewUniform(RGBA(yellow))
			z.DrawOp = draw.Src
			z.Draw(dst, dst.Bounds(), src, src.Bounds().Min)
			src = image.NewUniform(RGBA(highlight))
			z.DrawOp = draw.Over
			z.Draw(dst, dst.Bounds(), src, src.Bounds().Min)
			paint.NewImageOp(dst).Add(ops)
			lower := image.Rect(0, dy/2, dx, dy)
			tstack := op.Offset(lower.Min).Push(ops)
			cstack = clip.Rect(lower.Sub(lower.Min)).Push(ops)
			paint.PaintOp{}.Add(ops)
			cstack.Pop()
			tstack.Pop()

			event.Frame(ops)
		}
	}
	os.Exit(0)
}
