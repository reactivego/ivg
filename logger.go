package ivg

import (
	"fmt"
	"image/color"
)

type DestinationLogger struct {
	Destination
	Alt bool
}

func (d *DestinationLogger) Reset(viewbox ViewBox, palette [64]color.RGBA) {
	if !d.Alt {
		fmt.Printf("Reset(viewbox:%#v, colors:%#v)\n", viewbox, palette)
	} else {
		fmt.Printf("dst.Reset(%#v, %#v)\n", viewbox, palette)
	}
	if d.Destination != nil {
		d.Destination.Reset(viewbox, palette)
	}
}

func (d *DestinationLogger) SetCSel(cSel uint8) {
	if !d.Alt {
		fmt.Printf("SetCSel(cSel:%d)\n", cSel)
	} else {
		fmt.Printf("dst.SetCSel(%d)\n", cSel)
	}
	if d.Destination != nil {
		d.Destination.SetCSel(cSel)
	}
}

func (d *DestinationLogger) SetNSel(nSel uint8) {
	if !d.Alt {
		fmt.Printf("SetNSel(nSel:%d)\n", nSel)
	} else {
		fmt.Printf("dst.SetNSel(%d)\n", nSel)
	}
	if d.Destination != nil {
		d.Destination.SetNSel(nSel)
	}
}

func (d *DestinationLogger) SetCReg(adj uint8, incr bool, c Color) {
	if !d.Alt {
		fmt.Printf("SetCReg(adj:%d, incr:%t, c:%#v)\n", adj, incr, c)
	} else {
		fmt.Printf("dst.SetCReg(%d, %t, %#v)\n", adj, incr, c)
	}
	if d.Destination != nil {
		d.Destination.SetCReg(adj, incr, c)
	}
}

func (d *DestinationLogger) SetNReg(adj uint8, incr bool, f float32) {
	if !d.Alt {
		fmt.Printf("SetNReg(adj:%d, incr:%t, f:%.2f)\n", adj, incr, f)
	} else {
		fmt.Printf("dst.SetNReg(%d, %t, %.2f)\n", adj, incr, f)
	}
	if d.Destination != nil {
		d.Destination.SetNReg(adj, incr, f)
	}
}

func (d *DestinationLogger) SetLOD(lod0, lod1 float32) {
	if !d.Alt {
		fmt.Printf("SetLOD(lod0:%.2f, lod1:%.2f)\n", lod0, lod1)
	} else {
		fmt.Printf("dst.SetLOD(%.2f, %.2f)\n", lod0, lod1)
	}
	if d.Destination != nil {
		d.Destination.SetLOD(lod0, lod1)
	}
}

func (d *DestinationLogger) StartPath(adj uint8, x, y float32) {
	if !d.Alt {
		fmt.Printf("StartPath(adj:%d, x:%.2f, y:%.2f)\n", adj, x, y)
	} else {
		fmt.Printf("dst.StartPath(%d, %.2f, %.2f)\n", adj, x, y)
	}
	if d.Destination != nil {
		d.Destination.StartPath(adj, x, y)
	}
}

func (d *DestinationLogger) ClosePathEndPath() {
	if !d.Alt {
		fmt.Println("ClosePathEndPath()")
	} else {
		fmt.Println("dst.ClosePathEndPath()")
	}
	if d.Destination != nil {
		d.Destination.ClosePathEndPath()
	}
}

func (d *DestinationLogger) ClosePathAbsMoveTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("ClosePathAbsMoveTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.ClosePathAbsMoveTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.ClosePathAbsMoveTo(x, y)
	}
}

func (d *DestinationLogger) ClosePathRelMoveTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("ClosePathRelMoveTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.ClosePathRelMoveTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.ClosePathRelMoveTo(x, y)
	}
}

func (d *DestinationLogger) AbsHLineTo(x float32) {
	if !d.Alt {
		fmt.Printf("AbsHLineTo(x:%.2f)\n", x)
	} else {
		fmt.Printf("dst.AbsHLineTo(%.2f)\n", x)
	}
	if d.Destination != nil {
		d.Destination.AbsHLineTo(x)
	}
}

func (d *DestinationLogger) RelHLineTo(x float32) {
	if !d.Alt {
		fmt.Printf("RelHLineTo(x:%.2f)\n", x)
	} else {
		fmt.Printf("dst.RelHLineTo(%.2f)\n", x)
	}
	if d.Destination != nil {
		d.Destination.RelHLineTo(x)
	}
}

func (d *DestinationLogger) AbsVLineTo(y float32) {
	if !d.Alt {
		fmt.Printf("AbsVLineTo(y:%.2f)\n", y)
	} else {
		fmt.Printf("dst.AbsVLineTo(%.2f)\n", y)
	}
	if d.Destination != nil {
		d.Destination.AbsVLineTo(y)
	}
}

func (d *DestinationLogger) RelVLineTo(y float32) {
	if !d.Alt {
		fmt.Printf("RelVLineTo(y:%.2f)\n", y)
	} else {
		fmt.Printf("dst.RelVLineTo(%.2f)\n", y)
	}
	if d.Destination != nil {
		d.Destination.RelVLineTo(y)
	}
}

func (d *DestinationLogger) AbsLineTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsLineTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.AbsLineTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsLineTo(x, y)
	}
}

func (d *DestinationLogger) RelLineTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("RelLineTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.RelLineTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.RelLineTo(x, y)
	}
}

func (d *DestinationLogger) AbsSmoothQuadTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsSmoothQuadTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.AbsSmoothQuadTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsSmoothQuadTo(x, y)
	}
}

func (d *DestinationLogger) RelSmoothQuadTo(x, y float32) {
	if !d.Alt {
		fmt.Printf("RelSmoothQuadTo(x:%.2f, y:%.2f)\n", x, y)
	} else {
		fmt.Printf("dst.RelSmoothQuadTo(%.2f, %.2f)\n", x, y)
	}
	if d.Destination != nil {
		d.Destination.RelSmoothQuadTo(x, y)
	}
}

func (d *DestinationLogger) AbsQuadTo(x1, y1, x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsQuadTo(x1:%.2f, y1:%.2f, x:%.2f, y:%.2f)\n", x1, y1, x, y)
	} else {
		fmt.Printf("dst.AbsQuadTo(%.2f, %.2f, %.2f, %.2f)\n", x1, y1, x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsQuadTo(x1, y1, x, y)
	}
}

func (d *DestinationLogger) RelQuadTo(x1, y1, x, y float32) {
	if !d.Alt {
		fmt.Printf("RelQuadTo(x1:%.2f, y1:%.2f, x:%.2f, y:%.2f)\n", x1, y1, x, y)
	} else {
		fmt.Printf("dst.RelQuadTo(%.2f, %.2f, %.2f, %.2f)\n", x1, y1, x, y)
	}
	if d.Destination != nil {
		d.Destination.RelQuadTo(x1, y1, x, y)
	}
}

func (d *DestinationLogger) AbsSmoothCubeTo(x2, y2, x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsSmoothCubeTo(x2:%.2f, y2:%.2f, x:%.2f, y:%.2f)\n", x2, y2, x, y)
	} else {
		fmt.Printf("dst.AbsSmoothCubeTo(%.2f, %.2f, %.2f, %.2f)\n", x2, y2, x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsSmoothCubeTo(x2, y2, x, y)
	}
}

func (d *DestinationLogger) RelSmoothCubeTo(x2, y2, x, y float32) {
	if !d.Alt {
		fmt.Printf("RelSmoothCubeTo(x2:%.2f, y2:%.2f, x:%.2f, y:%.2f)\n", x2, y2, x, y)
	} else {
		fmt.Printf("dst.RelSmoothCubeTo(%.2f, %.2f, %.2f, %.2f)\n", x2, y2, x, y)
	}
	if d.Destination != nil {
		d.Destination.RelSmoothCubeTo(x2, y2, x, y)
	}
}

func (d *DestinationLogger) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsCubeTo(x1:%.2f, y1:%.2f, x2:%.2f, y2:%.2f, x:%.2f, y:%.2f)\n", x1, y1, x2, y2, x, y)
	} else {
		fmt.Printf("dst.AbsCubeTo(%.2f, %.2f, %.2f, %.2f, %.2f, %.2f)\n", x1, y1, x2, y2, x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsCubeTo(x1, y1, x2, y2, x, y)
	}
}

func (d *DestinationLogger) RelCubeTo(x1, y1, x2, y2, x, y float32) {
	if !d.Alt {
		fmt.Printf("RelCubeTo(x1:%.2f, y1:%.2f, x2:%.2f, y2:%.2f, x:%.2f, y:%.2f)\n", x1, y1, x2, y2, x, y)
	} else {
		fmt.Printf("dst.RelCubeTo(%.2f, %.2f, %.2f, %.2f, %.2f, %.2f)\n", x1, y1, x2, y2, x, y)
	}
	if d.Destination != nil {
		d.Destination.RelCubeTo(x1, y1, x2, y2, x, y)
	}
}

func (d *DestinationLogger) AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	if !d.Alt {
		fmt.Printf("AbsArcTo(rx:%.2f, ry:%.2f, xAxisRotation:%.2f, largeArc:%t, sweep:%t, x:%.2f, y:%.2f)\n", rx, ry, xAxisRotation, largeArc, sweep, x, y)
	} else {
		fmt.Printf("dst.AbsArcTo(%.2f, %.2f, %.2f, %t, %t, %.2f, %.2f)\n", rx, ry, xAxisRotation, largeArc, sweep, x, y)
	}
	if d.Destination != nil {
		d.Destination.AbsArcTo(rx, ry, xAxisRotation, largeArc, sweep, x, y)
	}
}

func (d *DestinationLogger) RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	if !d.Alt {
		fmt.Printf("RelArcTo(rx:%.2f, ry:%.2f, xAxisRotation:%.2f, largeArc:%t, sweep:%t, x:%.2f, y:%.2f)\n", rx, ry, xAxisRotation, largeArc, sweep, x, y)
	} else {
		fmt.Printf("dst.RelArcTo(%.2f, %.2f, %.2f, %t, %t, %.2f, %.2f)\n", rx, ry, xAxisRotation, largeArc, sweep, x, y)
	}
	if d.Destination != nil {
		d.Destination.RelArcTo(rx, ry, xAxisRotation, largeArc, sweep, x, y)
	}
}
