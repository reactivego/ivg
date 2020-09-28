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
    "github.com/reactivego/ivg/icon"
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

    actionInfoData, err := ActionInfoData()
    if err != nil {
        log.Fatal(err)
    }

    cache := icon.NewCache(icon.GioRasterizer)

    ops := new(op.Ops)
    for next := range window.Events() {
        if frame, ok := next.(system.FrameEvent); ok {
            ops.Reset()
            rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))
            blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
            if callOp, err := cache.FromData(actionInfoData, blue, rect, icon.AspectMeet, icon.Mid, icon.Mid); err == nil {
                callOp.Add(ops)
            } else {
                log.Fatal(err)
            }
            frame.Frame(ops)
        }
    }
    os.Exit(0)
}

// ActionInfoData generates ivg data bytes on the fly for the ActionInfo icon.
func ActionInfoData() ([]byte, error) {
    e := &encode.Encoder{}
    g := &generate.Generator{Destination: e}
    g.Reset(ivg.ViewBox{0, 0, 48, 48}, &ivg.DefaultPalette)
    g.SetPathData("M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 "+
        "20-20S35.05 4 24 4zm2 30h-4V22h4v12zm0-16h-4v-4h4v4z", 0, false)
    return e.Bytes()
}
