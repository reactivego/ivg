package mdicons

import (
	"github.com/reactivego/ivg"
	"golang.org/x/image/math/f32"
)

func ParsePath(enc ivg.Destination, p *Path, adjs map[float32]uint8, size float32, offset f32.Vec2, outSize float32, circles []Circle) error {
	adj := uint8(0)
	opacity := float32(1)
	if p.Opacity != nil {
		opacity = *p.Opacity
	} else if p.FillOpacity != nil {
		opacity = *p.FillOpacity
	}
	if opacity != 1 {
		var ok bool
		if adj, ok = adjs[opacity]; !ok {
			adj = uint8(len(adjs) + 1)
			adjs[opacity] = adj
			// Set CREG[0-adj] to be a blend of transparent (0x7f) and the
			// first custom palette color (0x80).
			enc.SetCReg(adj, false, ivg.BlendColor(uint8(opacity*0xff), 0x7f, 0x80))
		}
	}

	needStartPath := true
	if p.D != "" {
		needStartPath = false
		if err := ParsePathData(enc, p.D, adj, size, offset, outSize); err != nil {
			return err
		}
	}

	for _, c := range circles {
		// Normalize.
		cx := c.Cx * outSize / size
		cx -= outSize/2 + offset[0]
		cy := c.Cy * outSize / size
		cy -= outSize/2 + offset[1]
		r := c.R * outSize / size

		if needStartPath {
			needStartPath = false
			enc.StartPath(adj, cx-r, cy)
		} else {
			enc.ClosePathAbsMoveTo(cx-r, cy)
		}

		// Convert a circle to two relative arcTo ops, each of 180 degrees.
		// We can't use one 360 degree arcTo as the start and end point
		// would be coincident and the computation is degenerate.
		enc.RelArcTo(r, r, 0, false, true, +2*r, 0)
		enc.RelArcTo(r, r, 0, false, true, -2*r, 0)
	}

	enc.ClosePathEndPath()
	return nil
}
