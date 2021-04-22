// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package ivg provides rendering of IconVG icons in Gio.

IconVG (github.com/google/iconvg) is a compact, binary format for simple
vector graphics: icons, logos, glyphs and emoji.

Gio (gioui.org) implements portable immediate mode GUI programs in Go.

The code in this package does away with rendering the icon to an intermediate
bitmap image and instead directly uses Gio's vector API.
*/
package ivg

// TODO: shapes (circles, rects) and strokes? Or can we assume that authoring
// tools will convert shapes and strokes to paths?

// TODO: mark somehow that a graphic (such as a back arrow) should be flipped
// horizontally or its paths otherwise varied when presented in a Right-To-Left
// context, such as among Arabic and Hebrew text? Or should that be the
// responsibility of higher layers, selecting different IconVG graphics based
// on context, the way they would select different PNG graphics.

// TODO: hinting?
