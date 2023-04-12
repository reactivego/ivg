package mdicons

import (
	"fmt"
	"io"
	"strings"

	"github.com/reactivego/ivg/encode"
	"golang.org/x/image/math/f32"
)

func ParsePathData(enc *encode.Encoder, pathData string, adj uint8, size float32, offset f32.Vec2, outSize float32) error {
	pathData = strings.TrimSuffix(pathData, "z")
	r := strings.NewReader(pathData)

	var args [6]float32
	op, relative := byte(0), false
	for started := false; ; started = true {
		b, err := r.ReadByte()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch {
		case b == ' ':
			continue
		case 'A' <= b && b <= 'Z':
			op, relative = b, false
		case 'a' <= b && b <= 'z':
			op, relative = b, true
		default:
			r.UnreadByte()
		}

		n := 0
		switch op {
		case 'H', 'h', 'V', 'v':
			n = 1
		case 'L', 'l', 'M', 'm', 'T', 't':
			n = 2
		case 'Q', 'q', 'S', 's':
			n = 4
		case 'C', 'c':
			n = 6
		case 'Z', 'z':
		default:
			return fmt.Errorf("unknown opcode %c", b)
		}

		scan(&args, r, n)
		normalize(&args, n, op, size, offset, outSize, relative)

		switch op {
		case 'H':
			enc.AbsHLineTo(args[0])
		case 'h':
			enc.RelHLineTo(args[0])
		case 'V':
			enc.AbsVLineTo(args[0])
		case 'v':
			enc.RelVLineTo(args[0])
		case 'L':
			enc.AbsLineTo(args[0], args[1])
		case 'l':
			enc.RelLineTo(args[0], args[1])
		case 'M':
			if !started {
				enc.StartPath(adj, args[0], args[1])
			} else {
				enc.ClosePathAbsMoveTo(args[0], args[1])
			}
		case 'm':
			enc.ClosePathRelMoveTo(args[0], args[1])
		case 'T':
			enc.AbsSmoothQuadTo(args[0], args[1])
		case 't':
			enc.RelSmoothQuadTo(args[0], args[1])
		case 'Q':
			enc.AbsQuadTo(args[0], args[1], args[2], args[3])
		case 'q':
			enc.RelQuadTo(args[0], args[1], args[2], args[3])
		case 'S':
			enc.AbsSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 's':
			enc.RelSmoothCubeTo(args[0], args[1], args[2], args[3])
		case 'C':
			enc.AbsCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		case 'c':
			enc.RelCubeTo(args[0], args[1], args[2], args[3], args[4], args[5])
		}
	}
}

func scan(args *[6]float32, r *strings.Reader, n int) {
	for i := 0; i < n; i++ {
		for {
			if b, _ := r.ReadByte(); b != ' ' {
				r.UnreadByte()
				break
			}
		}
		fmt.Fscanf(r, "%f", &args[i])
	}
}

func normalize(args *[6]float32, n int, op byte, size float32, offset f32.Vec2, outSize float32, relative bool) {
	for i := 0; i < n; i++ {
		args[i] *= outSize / size
		if relative {
			continue
		}
		args[i] -= outSize / 2
		switch {
		case n != 1:
			args[i] -= offset[i&0x01]
		case op == 'H':
			args[i] -= offset[0]
		case op == 'V':
			args[i] -= offset[1]
		}
	}
}
