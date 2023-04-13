# ivg

    import "github.com/reactivego/ivg"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/ivg.svg)](https://pkg.go.dev/github.com/reactivego/ivg#section-documentation)

Package `ivg` provides rendering of [IconVG](https://github.com/google/iconvg) icons using a Rasterizer interface.
IconVG is a binary format for simple vector graphic icons.

The original [IconVG](https://golang.org/x/exp/shiny/iconvg) code was changed to render an IconVG graphic using a vector graphics Rasterizer interface. The original code rendered to a bitmap image.
The use of the rasterizer allows implementing different rasterizers for different purposes. For example, a rasterizer that renders to a bitmap image, or a rasterizer that renders to a gioui.org context.

The name of the `iconvg` package has been changed to `ivg` so we don't confuse people about what's what.

## File Format Versions

In order for the IconVG format to support animation in future incarnations. The format was simplified and updated to version 1 (FFV1), renaming the original format to FFV0 retroactively.

FFV1 targets representing static vector graphics icons, while the future FFV2 will target representing animated vector graphics icons.

The rationale for this was dicussed in a github proposal: [File Format Versions 1, 2 and Beyond](https://github.com/google/iconvg/issues/4#issue-905297018)

Below are links to the different File Format Versions of the spec:
- [IconVG FFV0](spec/iconvg-spec-v0.md)
- [IconVG FFV1](https://github.com/google/iconvg/blob/97b0c08e6e298f5f3606f79f3fb38cc0d64d3198/spec/iconvg-spec.md)

> NOTE: This package implements the [FFV0](spec/iconvg-spec-v0.md) version of the IconVG format.

## Changes

This package changes the original IconVG code in several ways.
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
12. Create examples in the `example` folder.
    - `playarrow` simplest example of rendering an icon, [see below](#example-playarrow).
    - `actioninfo` generate an icon on the fly, render it and cache the result, [see below](#example-actioninfo).
    - The following examples allow you to see rendering and speed differences between rasterizers by clicking on the image to switch rasterizer. 
        - `icons` renders golang.org/x/exp/shiny/materialdesign/icons. [see below](#example-icons).
        - `favicon` vector image with several blended layers. [see below](#example-favicon).
        - `cowbell` vector image with several blended layers including gradients. [see below](#example-cowbell).
        - `gradients` vector image with lots of gradients. [see below](#example-gradients).

## Acknowledgement

The code in this package is based on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

The specification of the IconVG format has recently been moved to a separate repository [github.com/google/iconvg](https://github.com/google/iconvg).


## License

Everything under the raster folder is Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).

All the other code is is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
