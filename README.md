# ivg

    import "github.com/reactivego/ivg"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/ivg.svg)](https://pkg.go.dev/github.com/reactivego/ivg#section-documentation)

Package `ivg` provides a powerful implementation for rendering [IconVG](https://github.com/google/iconvg) icons through a flexible Rasterizer interface. IconVG is an efficient binary format designed specifically for vector graphic icons.

This package enhances the original [IconVG](https://golang.org/x/exp/shiny/iconvg) implementation by introducing a modular vector graphics Rasterizer interface, replacing the original bitmap-only rendering system. This architectural change enables diverse rendering capabilities through different rasterizer implementations. Users can now choose between various output targets, such as bitmap images or gioui.org contexts, by implementing the appropriate rasterizer for their needs.

To maintain clarity and avoid namespace confusion with the original implementation, this package has been renamed from `iconvg` to `ivg`.

## File Format Versions

In order for the IconVG format to support animation in future incarnations. The format was simplified and updated to version 1 (FFV1), renaming the original format to FFV0 retroactively.

FFV1 targets representing static vector graphics icons, while the future FFV2 will target representing animated vector graphics icons.

The rationale for this was dicussed in a github proposal: [File Format Versions 1, 2 and Beyond](https://github.com/google/iconvg/issues/4#issue-905297018)

Below are links to the different File Format Versions of the spec:
- [IconVG FFV0](spec/iconvg-spec-v0.md)
- [IconVG FFV1](https://github.com/google/iconvg/blob/97b0c08e6e298f5f3606f79f3fb38cc0d64d3198/spec/iconvg-spec.md)

> NOTE: This package implements the [FFV0](spec/iconvg-spec-v0.md) version of the IconVG format.

## Code Organization

The original purpose of IconVG was to convert a material design icon in SVG format to a binary data blob that could be embedded in a Go program.

The code is organized in several packages that can be combined in different ways to create different IconVG render pipelines. The `Destination` interface is implemented both by the `Encoder` in package `encode` and by the `Renderer` in package `render`. The `Generator` type in the `generator` package just uses a `Destination` and doesn't care whether calls are generating a data blob or render directly to a `Rasterizer` via the `Renderer`.

```go
// Destination handles the actions decoded from an IconVG graphic's opcodes.
//
// When passed to Decode, the first method called (if any) will be Reset. No
// methods will be called at all if an error is encountered in the encoded form
// before the metadata is fully decoded.
type Destination interface {
	Reset(viewbox ViewBox, palette [64]color.RGBA)
	CSel() uint8
	SetCSel(cSel uint8)
	NSel() uint8
	SetNSel(nSel uint8)
	SetCReg(adj uint8, incr bool, c Color)
	SetNReg(adj uint8, incr bool, f float32)
	SetLOD(lod0, lod1 float32)

	StartPath(adj uint8, x, y float32)
	ClosePathEndPath()
	ClosePathAbsMoveTo(x, y float32)
	ClosePathRelMoveTo(x, y float32)

	AbsHLineTo(x float32)
	RelHLineTo(x float32)
	AbsVLineTo(y float32)
	RelVLineTo(y float32)
	AbsLineTo(x, y float32)
	RelLineTo(x, y float32)
	AbsSmoothQuadTo(x, y float32)
	RelSmoothQuadTo(x, y float32)
	AbsQuadTo(x1, y1, x, y float32)
	RelQuadTo(x1, y1, x, y float32)
	AbsSmoothCubeTo(x2, y2, x, y float32)
	RelSmoothCubeTo(x2, y2, x, y float32)
	AbsCubeTo(x1, y1, x2, y2, x, y float32)
	RelCubeTo(x1, y1, x2, y2, x, y float32)
	AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
	RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
}
```

A parser of SVG files reads the SVG and then calls methods on a `Destination` to produce a binary data blob.

For Material Design icons:

```
mdicons/Parse -> [Destination]encode/Encoder -> []byte
```

For more complex SVGs a Generator supports handling of e.g. gradients and transforms. The Generator is hooked up to a `Destination` to produce the binary data blob.

```
svgicon/Parse -> generate/Generator -> [Destination]encode/Encoder -> []byte
```

To actually render the icon, the binary data blob would be passed to a `Decoder` that would call methods on a `Renderer` hooked up to a `Rasterizer` to render the icon.

```go
// Rasterizer is a 2-D vector graphics rasterizer.
type Rasterizer interface {
	// Reset resets a Rasterizer as if it was just returned by NewRasterizer.
	// This includes setting z.DrawOp to draw.Over.
	Reset(w, h int)
	// Size returns the width and height passed to NewRasterizer or Reset.
	Size() image.Point
	// Bounds returns the rectangle from (0, 0) to the width and height passed to
	// Reset.
	Bounds() image.Rectangle
	// Pen returns the location of the path-drawing pen: the last argument to the
	// most recent XxxTo call.
	Pen() (x, y float32)
	// MoveTo starts a new path and moves the pen to (ax, ay). The coordinates
	// are allowed to be out of the Rasterizer's bounds.
	MoveTo(ax, ay float32)
	// LineTo adds a line segment, from the pen to (bx, by), and moves the pen to
	// (bx, by). The coordinates are allowed to be out of the Rasterizer's
	// bounds.
	LineTo(bx, by float32)
	// QuadTo adds a quadratic Bézier segment, from the pen via (bx, by) to (cx,
	// cy), and moves the pen to (cx, cy). The coordinates are allowed to be out
	// of the Rasterizer's bounds.
	QuadTo(bx, by, cx, cy float32)
	// CubeTo adds a cubic Bézier segment, from the pen via (bx, by) and (cx, cy)
	// to (dx, dy), and moves the pen to (dx, dy). The coordinates are allowed to
	// be out of the Rasterizer's bounds.
	CubeTo(bx, by, cx, cy, dx, dy float32)
	// ClosePath closes the current path.
	ClosePath()
	// Draw aligns r.Min in z with sp in src and then replaces the rectangle r in
	// z with the result of a Porter-Duff composition. The vector paths
	// previously added via the XxxTo calls become the mask for drawing src onto
	// z.
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}
```

Decoding a blob and rendering it to a `Rasterizer`:

```
[]byte -> decode/Decoder -> [Destination]render/Renderer -> [Rasterizer]raster/img/Rasterizer
```

To render and icon from SVG, the Generator can also be hooked up to the `Renderer` directly and the `Encoder`/`Decoder` phase would be skipped.

```
svgicon/Parse -> generate/Generator -> [Destination]render/Renderer -> [Rasterizer]raster/img/Rasterizer
```

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
12. Create examples in directory `raster/gio/example`.
    - `playarrow` simplest example of rendering an icon.
    - `actioninfo` generate an icon on the fly, render it and cache the result.
    - The following examples allow you to see rendering and speed differences between rasterizers by clicking on the image to switch rasterizer.
        - `icons` renders golang.org/x/exp/shiny/materialdesign/icons.
        - `favicon` vector image with several blended layers.
        - `cowbell` vector image with several blended layers including gradients.
        - `gradients` vector image with lots of gradients.

## Acknowledgement

The code in this package is based on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

The specification of the IconVG format has recently been moved to a separate repository [github.com/google/iconvg](https://github.com/google/iconvg).


## License

Everything under the raster folder is Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).

All the other code is is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
