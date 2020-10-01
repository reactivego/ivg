# ivg

    import "github.com/reactivego/ivg"

[![](.github/godev.svg)](https://pkg.go.dev/github.com/reactivego/ivg?tab=doc)
[![](.github/godoc.svg)](https://godoc.org/github.com/reactivego/ivg)

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
11. Create a rasterizer using "golang.org/x/image/vector" in directory `raster/vec`
12. Create a rasterizer using "gioui.org/op/clip" in directory `raster/gio`
    - Special case for `GradientConfig`, selectively sample gradient only inside path bounds.
14. Create examples in the `example` folder.
    - `playarrow` simplest example of rendering an icon, [see below](#example-playarrow).
    - `actioninfo` generate an icon on the fly, render it and cache the result, [see below](#example-actioninfo).
    - The following examples allow you to see rendering and speed differences between rasterizers by clicking on the image to switch rasterizer. 
        - `icons` renders golang.org/x/exp/shiny/materialdesign/icons. [see below](#example-icons).
        - `favicon` vector image with several blended layers. [see below](#example-favicon).
        - `cowbell` vector image with several blended layers including gradients. [see below](#example-cowbell).
        - `gradients` vector image with lots of gradients. [see below](#example-gradients).

## Example PlayArrow

![PlayButton Gio](testdata/playbutton-gio.png)

Simplest example of rendering an icon from an .ivg file stored in a slice of bytes.

```go
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

    "github.com/reactivego/ivg/icon"
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
    ops := new(op.Ops)
    for event := range window.Events() {
        if frame, ok := event.(system.FrameEvent); ok {
            ops.Reset()
            rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))
            blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
            rasterizer := icon.GioRasterizer
            if callOp, err := icon.FromData(AVPlayArrow, blue, rect, icon.AspectMeet, icon.Mid, icon.Mid, rasterizer); err == nil {
                callOp.Add(ops)
            } else {
                log.Fatal(err)
            }
            frame.Frame(ops)
        }
    }
    os.Exit(0)
}

// ACPlayArrow was taken from "golang.org/x/exp/shiny/materialdesign/icons"
var AVPlayArrow = []byte{
    0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0, 0xb0, 0xc0, 0x70, 0x64, 0xe9, 0xb8,
    0x20, 0xac, 0x64, 0xe1,
}

```
## Example ActionInfo

![ActionInfo Gio](testdata/actioninfo-gio.png)

Generating the .ivg bytes for an icon on the fly and then rendering it. Rendering operations are cached in an icon cache. The function `ActionInfoData()` is called once to programatically generate an .ivg byte slice using the following pipeline:
```
Generator -> Encoder
```
The resulting bytes are stored for later rendering during a system.FrameEvent. When the icon needs to be rendered, call the icon.Cache FromData method with the .ivg data bytes and additional arguments.
The icon cache uses the following pipeline to render the icon.
 
```
Decoder -> Renderer -> Rasterizer
```
The cache stores the resulting op.CallOp keyed on the icon data and parameters used for rendering.

```go
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
```


# Example Icons

| ![Icons Gio](testdata/icons-gio.png) | ![Icons Vector](testdata/icons-vec.png) |
|:---:|:---:|
| Gio | Vec |

The Icons example program takes icons from the package `"golang.org/x/exp/shiny/materialdesign/icons"` and renders them. These icons are just a few layers filled with a single color. The example uses function `FromData` from package `"github.com/reactivego/ivg/icon"` for rendering. `FromData` uses the following pipeline:
```
Decoder -> Renderer -> Rasterizer
```
By clicking the window you can switch between rendering using the Gio or Vec (`"golang.org/x/image/vector"`) rasterizer. 

The rendering using Gio is extremely quick as it only generates a few clipping & blending operations that are put in an operator queue. The actual rendering by Gio takes place on the GPU. For the Vec rasterizer all the pixels of the image need to be pre-generated, which takes relatively a long time.

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "golang.org/x/exp/shiny/materialdesign/colornames"

    "gioui.org/app"
    "gioui.org/f32"
    "gioui.org/io/pointer"
    "gioui.org/io/system"
    "gioui.org/op"
    "gioui.org/op/paint"
    "gioui.org/unit"

    "github.com/reactivego/ivg/icon"
)

func main() {
    go Icons()
    app.Main()
}

func Icons() {
    window := app.NewWindow(
        app.Title("IVG - Icons"),
        app.Size(unit.Dp(768), unit.Dp(768)),
    )
    rasterizer := icon.GioRasterizer
    ops := new(op.Ops)
    backdrop := new(int)
    index := 0
    for next := range window.Events() {
        if frame, ok := next.(system.FrameEvent); ok {
            ops.Reset()

            // initial window rect in pixels
            rect := f32.Rect(0, 0, float32(frame.Size.X), float32(frame.Size.Y))

            // backdrop switch renderer on release and fill rectangle
            pointer.InputOp{Tag: backdrop, Types: pointer.Release}.Add(ops)
            for _, next := range frame.Queue.Events(backdrop) {
                if event, ok := next.(pointer.Event); ok {
                    if event.Type == pointer.Release {
                        switch rasterizer {
                        case icon.GioRasterizer:
                            rasterizer = icon.VecRasterizer
                        case icon.VecRasterizer:
                            rasterizer = icon.GioRasterizer
                        }
                    }
                }
            }
            paint.ColorOp{Color: colornames.Grey800}.Add(ops)
            paint.PaintOp{Rect: rect}.Add(ops)

            // device independent content rect calculation
            pt32 := func(x, y unit.Value) f32.Point {
                return f32.Pt(float32(frame.Metric.Px(x)), float32(frame.Metric.Px(y)))
            }
            margin := pt32(unit.Dp(12), unit.Dp(12))
            lefttop := pt32(frame.Insets.Left, frame.Insets.Top).Add(margin)
            rightbottom := pt32(frame.Insets.Right, frame.Insets.Bottom).Add(margin)
            rect = f32.Rectangle{Min: rect.Min.Add(lefttop), Max: rect.Max.Sub(rightbottom)}

            // fill content rect
            op.Offset(rect.Min).Add(ops)
            rect = f32.Rectangle{Max: rect.Size()}
            paint.ColorOp{Color: colornames.Grey300}.Add(ops)
            paint.PaintOp{Rect: rect}.Add(ops)

            // select next icon and paint
            n := uint(len(IconCollection))
            ico := IconCollection[(uint(index)+n)%n]
            index++
            start := time.Now()
            if callOp, err := icon.FromData(ico.data, colornames.LightBlue600, rect, icon.AspectMeet, icon.Mid, icon.Mid, rasterizer); err == nil {
                callOp.Add(ops)
            } else {
                log.Fatal(err)
            }
            switch rasterizer {
            case icon.GioRasterizer:
                PrintText(fmt.Sprintf("Gio (%v)", time.Since(start).Round(time.Microsecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
            case icon.VecRasterizer:
                PrintText(fmt.Sprintf("Vec (%v)", time.Since(start).Round(time.Millisecond)), rect.Min, 0.0, 0.0, rect.Dx(), H5, ops)
            }

            at := time.Now().Add(500 * time.Millisecond)
            op.InvalidateOp{At:at}.Add(ops)
            frame.Frame(ops)
        }
    }
    os.Exit(0)
}
```

# Example Favicon

| ![Favicon Gio](testdata/favicon-gio.png) | ![Favicon Vector](testdata/favicon-vec.png) |
|:---:|:---:|
| Gio | Vec |

Favicon programatically renders a vector image of a Gopher using multiple layers with translucency. It uses the following pipeline:
```
Generator -> Renderer -> Rasterizer
```
This example hooks up the generator directly to the renderer and forgoes the `Encoder -> Decoder` stages. The rendering using Gio is extremely quick as it only generates a few clipping & blending operations that are put in an operator queue. The actual rendering by Gio takes place on the GPU. For the Vec (`"golang.org/x/image/vector"`) rasterizer all the pixels of the image need to be pre-generated, which takes relatively a long time.

The resulting images are a little bit different. Look under the nose of the Gopher. Gio produces lighter results than Vec. The reason for this is that Vec performs blending in sRGB space whereas Gio performs the blending in Linear space.

# Example Cowbell

| ![Cowbell Gio](testdata/cowbell-gio.png) | ![Cowbell Vector](testdata/cowbell-vec.png) |
|:---:|:---:|
| Gio | Vec |

Cowbell programatically renders a vector image of a Cowbell using multiple layers with gradients and translucency. It uses the following pipeline:
```
Generator -> Renderer -> Rasterizer
```
The rendering takes relatively long because the gradients need to be pre-generated even for the Gio renderer.

# Example Gradients

| ![Gradients Gio](testdata/gradients-gio.png) | ![Gradients Vector](testdata/gradients-vec.png) |
|:---:|:---:|
| Gio | Vec |

Gradients uses the following pipeline to programatically render a vector image consisting of multiple different gradients:
```
Generator -> Renderer -> Rasterizer
```
The rendering takes relatively long because the gradients need to be pre-generated even for the Gio renderer. But even considering that, Gio is approximately 8 times faster than Vec.

## Acknowledgement

This code is base on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

## License

Everything under the raster folder is Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).

All the other code is is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
