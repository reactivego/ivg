package raster

import "image/color"

// GradientConfig interface could be used in the future to extract the gradient
// configuration of a source image and have it generated on the GPU.
type GradientConfig interface {
	// GradientShape returns 0 for a linear gradient and 1 for a radial
	// gradient.
	GradientShape() int
	// SpreadMethod returns 0 for 'none', 1 for 'pad', 2 for 'reflect', 3 for
	// 'repeat'.
	SpreadMethod() int
	// StopColors returns the colors of the gradient stops.
	StopColors() []color.RGBA
	// StopOffsets returns the offsets of the gradient stops.
	StopOffsets() []float64
	// Transform is the pixel space to gradient space affine transformation
	// matrix.
	// | a b c |
	// | d e f |
	Transform() (a, b, c, d, e, f float64)
}
