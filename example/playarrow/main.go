// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/ivg/raster/gio"
)

func main() {
	go PlayArrow()
	app.Main()
}

func PlayArrow() {
	window := app.NewWindow(
		app.Title("IVG - PlayArrow"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	playArrow, err := gio.NewIcon([]byte{
		// AVPlayArrow data taken from "golang.org/x/exp/shiny/materialdesign/icons"
		0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0,
		0xb0, 0xc0, 0x70, 0x64, 0xe9, 0xb8, 0x20, 0xac, 0x64, 0xe1,
	})
	if err != nil {
		log.Fatal(err)
	}
	ops := new(op.Ops)
	for event := range window.Events() {
		if frame, ok := event.(system.FrameEvent); ok {
			ops.Reset()

			viewRect := playArrow.AspectMeet(frame.Size, 0.5, 0.5)
			blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
			callOp, err := gio.Rasterize(playArrow, viewRect, gio.WithColors(blue))

			if err != nil {
				log.Fatal(err)
			}
			callOp.Add(ops)
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
