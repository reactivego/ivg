// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package ivg provides rendering of IconVG icons.

IconVG (github.com/google/iconvg) is a compact, binary format for simple
vector graphics: icons, logos, glyphs and emoji.

The code in this package does away with rendering the icon to an intermediate
bitmap image and instead directly uses a vector Rasterizer interface.
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
