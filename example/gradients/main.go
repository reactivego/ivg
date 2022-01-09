// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sync"
	"time"

	"eliasnaur.com/font/roboto/robotoregular"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/opentype"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/math/fixed"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/raster/gio"
)

func main() {
	go Gradients()
	app.Main()
}

func Gradients() {
	window := app.NewWindow(
		app.Title("IVG - Gradients"),
		app.Size(unit.Dp(768), unit.Dp(768)),
	)
	type Backend struct {
		Name   string
		Driver gio.Driver
	}
	backend := Backend{"Gio", gio.Gio}
	Grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
	Grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
	gradients := GradientsImage{}
	ops := new(op.Ops)
	backdrop := new(int)
	for next := range window.Events() {
		if frame, ok := next.(system.FrameEvent); ok {
			ops.Reset()

			// clicking on backdrop will switch active renderer
			pointer.InputOp{Tag: backdrop, Types: pointer.Release}.Add(ops)
			for _, next := range frame.Queue.Events(backdrop) {
				if event, ok := next.(pointer.Event); ok {
					if event.Type == pointer.Release {
						switch backend.Name {
						case "Gio":
							backend = Backend{"Vec", gio.Vec}
						case "Vec":
							backend = Backend{"Gio", gio.Gio}
						}
					}
				}
			}

			// fill the whole backdrop rectangle
			paint.ColorOp{Color: Grey800}.Add(ops)
			paint.PaintOp{}.Add(ops)

			// device independent content rect calculation
			margin := unit.Dp(12)
			minX := unit.Add(frame.Metric, margin, frame.Insets.Left)
			minY := unit.Add(frame.Metric, margin, frame.Insets.Top)
			maxX := unit.Add(frame.Metric, unit.Px(float32(frame.Size.X)), frame.Insets.Right.Scale(-1), margin.Scale(-1))
			maxY := unit.Add(frame.Metric, unit.Px(float32(frame.Size.Y)), frame.Insets.Bottom.Scale(-1), margin.Scale(-1))
			contentRect := f32.Rect(
				float32(frame.Metric.Px(minX)), float32(frame.Metric.Px(minY)),
				float32(frame.Metric.Px(maxX)), float32(frame.Metric.Px(maxY)))
			contentMin := image.Pt(int(contentRect.Min.X), int(contentRect.Min.Y))
			contentSize := image.Pt(int(contentRect.Dx()), int(contentRect.Dy()))

			// fill content rect
			paint.ColorOp{Color: Grey300}.Add(ops)
			tstack := op.Offset(contentRect.Min).Push(ops)
			cstack := clip.Rect(image.Rectangle{Max: contentSize}).Push(ops)
			paint.PaintOp{}.Add(ops)
			cstack.Pop()
			tstack.Pop()

			// scale the viewbox of the icon to the content rect
			viewRect := gradients.AspectMeet(contentSize, 0.5, 0.5).Add(contentMin)

			// render actual content
			start := time.Now()
			if callOp, err := gio.Rasterize(gradients, viewRect, gio.WithDriver(backend.Driver)); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			msg := fmt.Sprintf("%s (%v)", backend.Name, time.Since(start).Round(time.Microsecond))
			H5 := Style(H5, WithMetric(frame.Metric))
			PrintText(msg, contentRect.Min, 0.0, 0.0, contentRect.Dx(), H5, ops)

			frame.Frame(ops)
		}
	}
	os.Exit(0)
}

type GradientsImage struct{}

var GradientsViewBox = ivg.ViewBox{
	MinX: -32, MinY: -32,
	MaxX: +32, MaxY: +32,
}

func (g GradientsImage) AspectMeet(size image.Point, ax, ay float32) image.Rectangle {
	minx, miny, maxx, maxy := GradientsViewBox.AspectMeet(float32(size.X), float32(size.Y), ax, ay)
	return image.Rect(int(minx), int(miny), int(maxx), int(maxy))
}

func (GradientsImage) Name() string {
	return "Gradients"
}

func (g GradientsImage) RenderOn(dst ivg.Destination, col ...color.RGBA) error {
	gen := generate.Generator{}
	gen.SetDestination(dst)

	// Uncomment lines below to log destination calls.
	// logger := &ivg.DestinationLogger{Destination: renderer}
	// gen.SetDestination(logger)

	gen.Reset(GradientsViewBox, ivg.DefaultPalette)

	rgb := []generate.GradientStop{
		{Offset: 0.00, Color: color.RGBA{0xff, 0x00, 0x00, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0x00, 0xff, 0x00, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0x00, 0x00, 0xff, 0xff}},
		{Offset: 1.00, Color: color.RGBA{0x00, 0x00, 0x00, 0xff}},
	}
	cmy := []generate.GradientStop{
		{Offset: 0.00, Color: color.RGBA{0x00, 0xff, 0xff, 0xff}},
		{Offset: 0.25, Color: color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{Offset: 0.50, Color: color.RGBA{0xff, 0x00, 0xff, 0xff}},
		{Offset: 0.75, Color: color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{Offset: 1.00, Color: color.RGBA{0xff, 0xff, 0x00, 0xff}},
	}

	x1, y1 := float32(-12), float32(-30)
	x2, y2 := float32(+12), float32(-18)
	minX, minY := float32(-30), float32(-30)
	maxX, maxY := float32(+30), float32(-18)

	gen.SetLinearGradient(x1, y1, x2, y2, generate.GradientSpreadNone, rgb)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	x1, y1 = -12, -14
	x2, y2 = +12, -2
	minY = -14
	maxY = -2

	gen.SetLinearGradient(x1, y1, x2, y2, generate.GradientSpreadPad, cmy)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	cx, cy := float32(-8), float32(+8)
	rx, ry := float32(0), float32(+16)
	minY = +2
	maxY = +14

	gen.SetCircularGradient(cx, cy, rx, ry, generate.GradientSpreadReflect, rgb)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	cx, cy = -8, +24
	rx, ry = 0, 16
	minY = +18
	maxY = +30

	gen.SetCircularGradient(cx, cy, rx, ry, generate.GradientSpreadRepeat, cmy)
	gen.StartPath(0, minX, minY)
	gen.AbsHLineTo(maxX)
	gen.AbsVLineTo(maxY)
	gen.AbsHLineTo(minX)
	gen.ClosePathEndPath()

	return nil
}

func PrintText(txt string, pt f32.Point, ax, ay, width float32, style TextStyle, ops *op.Ops) (dx, dy float32) {
	size := fixed.I(style.Size)
	lines := style.Cache.LayoutString(style.Font, size, int(width), txt)
	for _, line := range lines {
		dy += float32(line.Ascent.Ceil() + line.Descent.Ceil())
		lineWidth := float32(line.Width.Ceil())
		if dx < lineWidth {
			dx = lineWidth
		}
	}
	offset := f32.Pt(pt.X-ax*dx, pt.Y-ay*dy)
	for _, line := range lines {
		offset.Y += float32(line.Ascent.Ceil())
		tstack := op.Offset(offset).Push(ops)
		offset.Y += float32(line.Descent.Ceil())
		cstack := style.Cache.Shape(style.Font, size, line.Layout).Push(ops)
		paint.ColorOp{Color: style.Color}.Add(ops)
		paint.PaintOp{}.Add(ops)
		cstack.Pop()
		tstack.Pop()
	}
	return
}

// Text Styles

var H5 = TextStyle{Font: RobotoNormal, Size: 24, Color: color.NRGBA{0, 0, 0, 255}, Cache: Cache()}

type TextStyle struct {
	Font  text.Font
	Size  int
	Color color.NRGBA
	Cache *text.Cache
}

type StyleOption func(*TextStyle)

func WithMetric(m unit.Metric) StyleOption {
	return func(s *TextStyle) {
		s.Size = m.Px(unit.Sp(float32(s.Size)))
	}
}

func WithColor(c color.Color) StyleOption {
	return func(s *TextStyle) {
		s.Color = color.NRGBAModel.Convert(c).(color.NRGBA)
	}
}

func Style(s TextStyle, options ...StyleOption) TextStyle {
	for _, o := range options {
		o(&s)
	}
	return s
}

// Fonts & Cache

var RobotoNormal = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: text.Normal}

func Cache() *text.Cache {
	cache.once.Do(func() {
		face := func(ttf []byte) text.Face {
			if face, err := opentype.Parse(ttf); err == nil {
				return face
			}
			panic("failed to parse font")
		}
		cache.ptr = text.NewCache([]text.FontFace{
			{Font: RobotoNormal, Face: face(robotoregular.TTF)},
		})
	})
	return cache.ptr
}

var cache struct {
	once sync.Once
	ptr  *text.Cache
}
