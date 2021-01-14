// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image/color"
	"sync"

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
	once   sync.Once
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
	RobotoThin   = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: 100}
	RobotoLight  = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: 200}
	RobotoNormal = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: text.Normal /*400*/}
	RobotoMedium = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: text.Medium /*500*/}
	RobotoBold   = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: text.Bold /*600*/}
	RobotoBlack  = text.Font{Typeface: "Roboto", Variant: "", Style: text.Regular, Weight: 800}
)

type TextStyle struct {
	Font  text.Font
	Size  int
	Color color.NRGBA
}

var (
	Black     = color.NRGBA{0, 0, 0, 255}
	H1        = TextStyle{RobotoThin, 96, Black}   // w300
	H2        = TextStyle{RobotoLight, 60, Black}  // w300
	H3        = TextStyle{RobotoNormal, 48, Black} // w400
	H4        = TextStyle{RobotoNormal, 34, Black} // w400
	H5        = TextStyle{RobotoNormal, 24, Black} // w400
	H6        = TextStyle{RobotoMedium, 20, Black} // w500
	Subtitle1 = TextStyle{RobotoNormal, 16, Black} // w400
	Subtitle2 = TextStyle{RobotoMedium, 14, Black} // w500
	BodyText1 = TextStyle{RobotoNormal, 16, Black} // w400
	BodyText2 = TextStyle{RobotoNormal, 14, Black} // w400
	Button    = TextStyle{RobotoMedium, 14, Black} // w500
	Caption   = TextStyle{RobotoNormal, 12, Black} // w400
	Overline  = TextStyle{RobotoNormal, 10, Black} // w400
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
	offset := f32.Pt(pt.X-ax*dx, pt.Y-ay*dy)
	for _, line := range lines {
		state := op.Save(ops)
		offset.Y += float32(line.Ascent.Ceil())
		op.Offset(offset).Add(ops)
		offset.Y += float32(line.Descent.Ceil())
		shaper.Shape(style.Font, fixed.I(style.Size), line.Layout).Add(ops)
		paint.ColorOp{Color: style.Color}.Add(ops)
		paint.PaintOp{}.Add(ops)
		state.Load()
	}
	return
}
