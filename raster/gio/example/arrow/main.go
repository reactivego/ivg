// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	raster "github.com/reactivego/ivg/raster/gio"
)

func main() {
	go Arrow()
	app.Main()
}

func Arrow() {
	window := app.NewWindow(
		app.Title("IVG - Arrow"),
		app.Size(768, 768),
	)
	widget, err := raster.Icon(AVPlayArrow, 48, 48, raster.WithColors(Amber400))
	if err != nil {
		log.Fatal(err)
	}
	ops := new(op.Ops)
	for next := range window.Events() {
		if event, ok := next.(system.FrameEvent); ok {
			gtx := layout.NewContext(ops, event)
			widget(gtx)
			event.Frame(ops)
		}
	}
	os.Exit(0)
}

var (
	// From "golang.org/x/exp/shiny/materialdesign/icons"
	AVPlayArrow = []byte{
		0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50,
		0x50, 0xb0, 0xb0, 0xc0, 0x70, 0x64, 0xe9, 0xb8,
		0x20, 0xac, 0x64, 0xe1,
	}
	// From "golang.org/x/exp/shiny/materialdesign/colors"
	Amber400 = color.RGBA{0xff, 0xca, 0x28, 0xff} // rgb(255, 202, 40)
)
