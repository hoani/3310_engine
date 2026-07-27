package desktop

import (
	"image"
	"image/color"

	"github.com/hoani/3310_engine/engine"
)

type canvas struct {
	image *image.Paletted
}

var _ engine.Canvas = &canvas{}

func NewCanvas() *canvas {
	return &canvas{
		image: image.NewPaletted(
			image.Rect(0, 0, 84, 48),
			color.Palette{
				engine.PixelOff,
				engine.PixelOn,
			},
		),
	}
}

func (c *canvas) Clear() {
	idx := uint8(c.image.Palette.Index(engine.PixelOff))
	for i := range c.image.Pix {
		c.image.Pix[i] = idx
	}
}

func (c *canvas) Set(x, y int, val bool) {
	if val {
		c.image.Set(x, y, engine.PixelOn)
	} else {
		c.image.Set(x, y, engine.PixelOff)
	}
}

func (c *canvas) Get(x, y int) bool {
	return c.image.At(x, y) == engine.PixelOn
}

func (c *canvas) Width() int {
	return c.image.Rect.Size().X
}

func (c *canvas) Height() int {
	return c.image.Rect.Size().Y
}

func (c *canvas) Size() (x, y int16) {
	return int16(c.image.Rect.Size().X), int16(c.image.Rect.Size().Y)
}

func (c *canvas) SetPixel(x, y int16, col color.RGBA) {
	if col.A != 0 {
		if col.R == 0 && col.G == 0 && col.B == 0 {
			c.Set(int(x), int(y), false)
		} else {
			c.Set(int(x), int(y), true)
		}
	}
}

func (c *canvas) Display() error {
	return nil
}
