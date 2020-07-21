package main

import (
    "image/color"
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

func RenderIcon(cache *icon.Cache, data []byte, rect f32.Rectangle, ops *op.Ops) {
    blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
    callOp, err := cache.FromData(data, blue, rect, icon.AspectMeet, icon.Mid, icon.Mid)
    if err != nil {
        panic(err)
    }
    callOp.Add(ops)
}

func ActionInfo() (data []byte, err error) {
    e := &encode.Encoder{}
    g := &generate.Generator{Destination: e}
    g.Reset(ivg.ViewBox{0, 0, 48, 48}, &ivg.DefaultPalette)
    g.SetPathData("M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 "+
        "20-20S35.05 4 24 4zm2 30h-4V22h4v12zm0-16h-4v-4h4v4z", 0, false)
    return e.Bytes()
}

func main() {
    window := app.NewWindow(
        app.Title("IVG - ActionInfo"),
        app.Size(unit.Dp(768), unit.Dp(768)),
    )
    var ops = new(op.Ops)
    cache := icon.NewCache(icon.GioRasterizer)
    data, err := ActionInfo()
    if err != nil {
        panic(err)
    }
    go func() {
        for next := range window.Events() {
            switch e := next.(type) {
            case system.DestroyEvent:
                os.Exit(1)
            case system.FrameEvent:
                ops.Reset()
                dx, dy := float32(e.Size.X), float32(e.Size.Y)
                RenderIcon(cache, data, f32.Rect(0, 0, dx, dy), ops)
                e.Frame(ops)
            }
        }
        os.Exit(0)
    }()
    app.Main()
}
