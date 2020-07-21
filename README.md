# ivg

    import "github.com/reactivego/ivg"

[![](testdata/godev.svg)](https://pkg.go.dev/github.com/reactivego/ivg?tab=doc)
[![](testdata/godoc.svg)](https://godoc.org/github.com/reactivego/ivg)

[Gio](https://gioui.org) (immediate mode GUI in Go) uses [IconVG](https://golang.org/x/exp/shiny/iconvg) (binary format for simple vector graphic icons).
This code is a refactoring of the IconVG code. It removes the need for rendering to an intermediate RGBA image. Instead it uses Gio `clip.Path` functionality.

The name of the *IconVG* package has been changed to *ivg* so we don't confuse people about what's what.

The most important changes w.r.t. the original IconVG code are:

1. Separate code into packages with a clear purpose and responsibility for better cohesion and less coupling.
2. Split icon encoding into `encode` and `generate` package.
3. SVG gradient and path support is now part of `generate` package.
4. Rename `Rasterizer` to `Renderer` and place it in the `render` package.
5. Move `Destination` interface into root `ivg` package.
6. Make both `Encoder` and `Renderer` implement `Destination`.
7. Make both `Decoder` and `Generator` use only `Destination` interface.
8. `Generator` can now directly render by plugging in a `Renderer` (very useful).
9. `Encoder` can be plugged directly into a `Decoder` (useful for testing).
10. Abstract away rasterizing into a seperate package `raster`
    - Declare interface `Rasterizer`.
    - Declare interface `GradientConfig` implemented by `Renderer`.
11. Create a rasterizer using "golang.org/x/image/vector" in directory `raster/vector`
12. Create a rasterizer using "gioui.org/op/clip" in directory `raster/gio`
    - Special case for `GradientConfig`, selectively sample gradient only inside path bounds.
14. Create examples in the `example` folder.
    - `actioninfo` generate an icon on the fly, render it and cache the result, [see below](#example-actioninfo).
    - `cowbell` exposes rendering differences between gio and vector rasterizer (click to switch rasterizer).
    - `favicon` exposes rendering difference between gio and vector rasterizer (click to switch rasterizer).
    - `gradient` shows speed advantage of gio over vector rasterizer for gradient rendering (click to switch rasterizer).
    - `icons` renders golang.org/x/exp/shiny/materialdesign/icons (click to switch rasterizer).
    - `playarrow` simplest example of rendering an icon, [see below](#example-playarrow).

## Example PlayArrow

Simplest example of rendering an icon.

```go
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
```
## Example ActionInfo

Generating an icon on the fly and then rendering it. Rendering operations are cached in an icon cache.

```go
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
```

## Acknowledgement

This code is base on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

## License

Everything under the raster folder is Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).

All the other code is is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
