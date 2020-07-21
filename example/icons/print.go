// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image/color"
	"sync"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/math/fixed"

	"eliasnaur.com/font/roboto/robotoblack"
	"eliasnaur.com/font/roboto/robotoblackitalic"
	"eliasnaur.com/font/roboto/robotobold"
	"eliasnaur.com/font/roboto/robotobolditalic"
	"eliasnaur.com/font/roboto/robotoitalic"
	"eliasnaur.com/font/roboto/robotolight"
	"eliasnaur.com/font/roboto/robotolightitalic"
	"eliasnaur.com/font/roboto/robotomedium"
	"eliasnaur.com/font/roboto/robotomediumitalic"
	"eliasnaur.com/font/roboto/robotoregular"
	"eliasnaur.com/font/roboto/robotothin"
	"eliasnaur.com/font/roboto/robotothinitalic"

	"gioui.org/f32"
	"gioui.org/font/opentype"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
)

var (
	once       sync.Once
	roboto []text.FontFace
)

func RobotoFontFaces() []text.FontFace {
	register := func(fnt text.Font, ttf []byte) {
		face, err := opentype.Parse(ttf)
		if err != nil {
			panic(fmt.Sprintf("failed to parse font: %v", err))
		}
		fnt.Typeface = "Roboto"
		roboto = append(roboto, text.FontFace{Font: fnt, Face: face})
	}
	once.Do(func() {
		// Weight: Normal (400)
		register(text.Font{}, robotoregular.TTF)
		register(text.Font{Style: text.Italic}, robotoitalic.TTF)

		// Weight: Thin (100)
		register(text.Font{Weight: 100}, robotothin.TTF)
		register(text.Font{Style: text.Italic, Weight: 100}, robotothinitalic.TTF)

		// Weight: Light (200)
		register(text.Font{Weight: 200}, robotolight.TTF)
		register(text.Font{Style: text.Italic, Weight: 200}, robotolightitalic.TTF)

		// Weight: Medium (500)
		register(text.Font{Weight: text.Medium}, robotomedium.TTF)
		register(text.Font{Weight: text.Medium, Style: text.Italic}, robotomediumitalic.TTF)

		// Weight: Bold (600)
		register(text.Font{Weight: text.Bold}, robotobold.TTF)
		register(text.Font{Style: text.Italic, Weight: text.Bold}, robotobolditalic.TTF)

		// Weight: Black (800)
		register(text.Font{Weight: 800}, robotoblack.TTF)
		register(text.Font{Style: text.Italic, Weight: 800}, robotoblackitalic.TTF)
	})
	return roboto
}

var (
	RobotoThin   = text.Font{"Roboto", "", text.Regular, 100}
	RobotoLight  = text.Font{"Roboto", "", text.Regular, 200}
	RobotoNormal = text.Font{"Roboto", "", text.Regular, text.Normal /*400*/}
	RobotoMedium = text.Font{"Roboto", "", text.Regular, text.Medium /*500*/}
	RobotoBold   = text.Font{"Roboto", "", text.Regular, text.Bold /*600*/}
	RobotoBlack  = text.Font{"Roboto", "", text.Regular, 800}
)

type TextStyle struct {
	Font  text.Font
	Size  int
	Color color.RGBA
}

var (
	H1        = TextStyle{RobotoThin, 96, colornames.Black}   // w300
	H2        = TextStyle{RobotoLight, 60, colornames.Black}  // w300
	H3        = TextStyle{RobotoNormal, 48, colornames.Black} // w400
	H4        = TextStyle{RobotoNormal, 34, colornames.Black} // w400
	H5        = TextStyle{RobotoNormal, 24, colornames.Black} // w400
	H6        = TextStyle{RobotoMedium, 20, colornames.Black} // w500
	Subtitle1 = TextStyle{RobotoNormal, 16, colornames.Black} // w400
	Subtitle2 = TextStyle{RobotoMedium, 14, colornames.Black} // w500
	BodyText1 = TextStyle{RobotoNormal, 16, colornames.Black} // w400
	BodyText2 = TextStyle{RobotoNormal, 14, colornames.Black} // w400
	Button    = TextStyle{RobotoMedium, 14, colornames.Black} // w500
	Caption   = TextStyle{RobotoNormal, 12, colornames.Black} // w400
	Overline  = TextStyle{RobotoNormal, 10, colornames.Black} // w400
)

var shaper = text.NewCache(RobotoFontFaces())

func PrintText(txt string, pt f32.Point, ax, ay, width float32, style TextStyle, ops *op.Ops) (dx, dy float32) {
	lines := shaper.LayoutString(style.Font, fixed.I(style.Size), int(width), txt)
	for _, line := range lines {
		dy += float32(line.Ascent.Ceil() + line.Descent.Ceil())
		lineWidth := float32(line.Width.Ceil())
		if dx < lineWidth {
			dx = lineWidth
		}
	}
	txtPos := 0
	offset := f32.Pt(pt.X-ax*dx, pt.Y-ay*dy)
	for _, line := range lines {
		stack := op.Push(ops)
		bounds := f32.Rect(0, float32(-line.Ascent.Ceil()), float32(line.Width.Ceil()), float32(line.Descent.Ceil()))
		offset.Y -= bounds.Min.Y
		op.Offset(offset).Add(ops)
		offset.Y += bounds.Max.Y
		shaper.ShapeString(style.Font, fixed.I(style.Size), txt[txtPos:txtPos+line.Len], line.Layout).Add(ops)
		txtPos += line.Len
		paint.ColorOp{Color: style.Color}.Add(ops)
		paint.PaintOp{Rect: bounds}.Add(ops)
		stack.Pop()
	}
	return
}