package main

import (
    "image/color"
    "os"

    "gioui.org/app"
    "gioui.org/f32"
    "gioui.org/io/system"
    "gioui.org/op"
    "gioui.org/unit"

    "github.com/reactivego/ivg/icon"
)

func RenderIcon(data []byte, rect f32.Rectangle, raster icon.Rasterizer, ops *op.Ops) {
    blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
    callOp, err := icon.FromData(data, blue, rect, icon.AspectMeet, icon.Mid, icon.Mid, raster)
    if err != nil {
        panic(err)
    }
    callOp.Add(ops)
}

// ACPlayArrow was taken from "golang.org/x/exp/shiny/materialdesign/icons"
var AVPlayArrow = []byte{
    0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0, 0xb0, 0xc0, 0x70, 0x64, 0xe9, 0xb8,
    0x20, 0xac, 0x64, 0xe1,
}

func main() {
    window := app.NewWindow(
        app.Title("IVG - PlayArrow"),
        app.Size(unit.Dp(768), unit.Dp(768)),
    )
    var ops = new(op.Ops)
    go func() {
        for next := range window.Events() {
            switch e := next.(type) {
            case system.DestroyEvent:
                os.Exit(1)
            case system.FrameEvent:
                ops.Reset()
                dx, dy := float32(e.Size.X), float32(e.Size.Y)
                RenderIcon(AVPlayArrow, f32.Rect(0, 0, dx, dy), icon.GioRasterizer, ops)
                e.Frame(ops)
            }
        }
        os.Exit(0)
    }()
    app.Main()
}