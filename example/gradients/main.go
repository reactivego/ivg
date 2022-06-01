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
		Drawer gio.Drawer
	}
	backend := Backend{"Gio", gio.DrawGio}

	grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
	grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
	black := color.NRGBA{A: 255}

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
							backend = Backend{"Vec", gio.DrawVec}
						case "Vec":
							backend = Backend{"Gio", gio.DrawGio}
						}
					}
				}
			}

			// fill the whole backdrop rectangle
			paint.ColorOp{Color: grey800}.Add(ops)
			paint.PaintOp{}.Add(ops)

			// device independent content rect calculation
			margin := unit.Dp(12)
			minX := frame.Metric.Dp(margin + frame.Insets.Left)
			minY := frame.Metric.Dp(margin + frame.Insets.Top)
			maxX := frame.Size.X - frame.Metric.Dp(frame.Insets.Right+margin)
			maxY := frame.Size.Y - frame.Metric.Dp(frame.Insets.Bottom+margin)
			contentRect := image.Rect(minX, minY, maxX, maxY)

			// fill content rect
			paint.FillShape(ops, grey300, clip.Rect(contentRect).Op())

			// scale the viewbox of the icon to the content rect
			viewRect := gradients.AspectMeet(contentRect.Size(), ivg.Mid, ivg.Max).Add(contentRect.Min)

			// render actual content
			start := time.Now()
			if callOp, err := gio.Rasterize(gradients, viewRect, gio.WithDrawer(backend.Drawer)); err == nil {
				callOp.Add(ops)
			} else {
				log.Fatal(err)
			}
			msg := fmt.Sprintf("%s (%v)", backend.Name, time.Since(start).Round(time.Microsecond))
			H5.Text(ops, frame.Metric, contentRect, 0.0, 0.0, contentRect.Dx(), black, msg)

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

func (g GradientsImage) RenderOn(dst ivg.Destination, col ...color.Color) error {
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

// Text Rendering

var H5 = TextStyle{Shaper: Cache(), Font: RobotoNormal, Size: 24}

type TextStyle struct {
	Shaper text.Shaper
	Font   text.Font
	Size   unit.Sp
}

var Locale = system.Locale{Language: "en-US", Direction: system.LTR}

func (s TextStyle) Text(ops *op.Ops, metric unit.Metric, rect image.Rectangle, ax, ay float32, maxWidth int, textColor color.NRGBA, txt string) (dx, dy int) {
	size := fixed.I(metric.Sp(s.Size))
	lines := s.Shaper.LayoutString(s.Font, size, maxWidth, Locale, txt)
	for _, line := range lines {
		dy += line.Ascent.Ceil()
		if dx < line.Width.Ceil() {
			dx = line.Width.Ceil()
		}
		dy += line.Descent.Ceil()
	}
	offset := rect.Min.Add(image.Pt(int(ax*float32(rect.Dx()-dx)), int(ay*float32(rect.Dy()-dy))))
	for _, line := range lines {
		shape := clip.Outline{Path: s.Shaper.Shape(s.Font, size, line.Layout)}.Op()
		offset.Y += line.Ascent.Ceil()
		tstack := op.Offset(offset).Push(ops)
		paint.FillShape(ops, textColor, shape)
		tstack.Pop()
		offset.Y += line.Descent.Ceil()
	}
	return
}

var RobotoNormal = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: text.Normal}

func Cache() text.Shaper {
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
