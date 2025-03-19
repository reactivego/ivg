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

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
	"github.com/reactivego/ivg/generate"
	raster "github.com/reactivego/ivg/raster/gio"
)

func main() {
	go Info()
	app.Main()
}

func Info() {
	window := app.NewWindow(
		app.Title("IVG - Info"),
		app.Size(768, 768),
	)

	data, err := InfoIVG()
	if err != nil {
		log.Fatal(err)
	}

	blue := color.NRGBA{0x21, 0x96, 0xf3, 0xff}

	widget, err := raster.Icon(data, 48, 48, raster.WithColors(blue))
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

// InfoIVG generates ivg data bytes on the fly for the Info icon.
func InfoIVG() ([]byte, error) {
	enc := &encode.Encoder{}
	gen := &generate.Generator{Destination: enc}
	gen.Reset(ivg.ViewBox{MinX: 0, MinY: 0, MaxX: 48, MaxY: 48}, ivg.DefaultPalette)
	gen.SetPathData("M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 "+
		"20-20S35.05 4 24 4zm2 30h-4V22h4v12zm0-16h-4v-4h4v4z", 0)
	return enc.Bytes()
}
