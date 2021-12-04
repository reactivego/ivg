// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/encode"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/raster/gio"
)

func main() {
	go ActionInfo()
	app.Main()
}

func ActionInfo() {
	window := app.NewWindow(
		app.Title("IVG - ActionInfo"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	// generate ivg data bytes on the fly for the ActionInfo icon.
	enc := &encode.Encoder{}
	gen := &generate.Generator{Destination: enc}
	gen.Reset(ivg.ViewBox{MinX: 0, MinY: 0, MaxX: 48, MaxY: 48}, ivg.DefaultPalette)
	gen.SetPathData("M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 "+
		"20-20S35.05 4 24 4zm2 30h-4V22h4v12zm0-16h-4v-4h4v4z", 0, false)
	actionInfoData, err := enc.Bytes()
	if err != nil {
		log.Fatal(err)
	}
	actionInfo, err := gio.NewIcon(actionInfoData)
	if err != nil {
		log.Fatal(err)
	}
	cache := gio.NewIconCache()
	ops := new(op.Ops)
	for next := range window.Events() {
		if frame, ok := next.(system.FrameEvent); ok {
			ops.Reset()
			contentRect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))
			viewRect := actionInfo.AspectMeet(contentRect, ivg.Mid, ivg.Mid)
			blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
			if callOp, err := cache.Rasterize(actionInfo, viewRect, gio.WithColors(blue)); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			frame.Frame(ops)
		}
	}
	os.Exit(0)
}
