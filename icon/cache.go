package icon

import (
	"crypto/md5"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
)

type CachableIcon interface {
	Icon

	// Rasterize will rasterize the icon using a default or internal rasterizer.
	Rasterize(rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error)

	// Name is the unique name of the icon
	Name() string
}

// Cache is an icon cache that caches op.CallOp values returned by a call to
// the Rasterize method.
type Cache struct {
	item map[key]op.CallOp
}

type key struct {
	checksum [md5.Size]byte
	rect     f32.Rectangle
}

// NewCache returns a new icon cache.
func NewCache() *Cache {
	return &Cache{make(map[key]op.CallOp)}
}

// Rasterize returns a gio op.CallOp that paints the 'icon' inside the given
// rectangle 'rect' overiding colors with the colors 'col'.
func (c *Cache) Rasterize(icon CachableIcon, rect f32.Rectangle, col ...color.RGBA) (op.CallOp, error) {
	data := []byte(icon.Name())
	for _, c := range col {
		data = append(data, c.R, c.G, c.B, c.A)
	}
	key := key{md5.Sum(data), rect}
	if callOp, present := c.item[key]; present {
		return callOp, nil
	}
	if callOp, err := icon.Rasterize(rect, col...); err == nil {
		c.item[key] = callOp
		return callOp, nil
	} else {
		return op.CallOp{}, err
	}
}
