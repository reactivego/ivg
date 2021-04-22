package ivg_test

import (
	"bytes"
	"testing"

	"github.com/reactivego/ivg"
	"github.com/reactivego/ivg/decode"
	"github.com/reactivego/ivg/encode"
	"github.com/reactivego/ivg/generate"
	"github.com/reactivego/ivg/render"
)

func TestEncodeDecode(t *testing.T) {
	var e = &encode.Encoder{}
	e.HighResolutionCoordinates = true
	e.Reset(
		ivg.ViewBox{
			MinX: -24, MinY: -24,
			MaxX: +24, MaxY: +24,
		},
		ivg.DefaultPalette,
	)

	e.StartPath(0, 0, -20)
	e.AbsCubeTo(-11.05, -20, -20, -11.05, -20, 0)
	e.RelSmoothCubeTo(8.95, 20, 20, 20)
	e.RelSmoothCubeTo(20, -8.95, 20, -20)
	e.AbsSmoothCubeTo(11.05, -20, 0, -20)
	e.ClosePathRelMoveTo(2, 30)
	e.RelHLineTo(-4)
	e.AbsVLineTo(-2)
	e.RelHLineTo(4)
	e.RelVLineTo(12)
	e.ClosePathRelMoveTo(0, -16)
	e.RelHLineTo(-4)
	e.RelVLineTo(-4)
	e.RelHLineTo(4)
	e.RelVLineTo(4)
	e.ClosePathEndPath()

	expect, err := e.Bytes()
	if err != nil {
		t.Fatalf("encoding: %v", err)
	}

	e = &encode.Encoder{}
	e.HighResolutionCoordinates = true
	if err := decode.Decode(e, expect); err != nil {
		t.Fatalf("decoding: %v", err)
	}
	actual, err := e.Bytes()
	if err != nil {
		t.Fatalf("encoding: %v", err)
	} else {
		if len(expect) != len(actual) {
			t.Fatalf("len(actual)!=len(expect): %d %d", len(actual), len(expect))
		} else {
			if !bytes.Equal(expect, actual) {
				t.Fatal("actual!=expect")
			}
		}
	}

	var r = &render.Renderer{}

	var g generate.Generator
	g.SetDestination(e)
	g.SetDestination(r)
}
