package desktop

import (
	"image"
	"image/color"

	"github.com/hoani/3310_engine/engine"
)

var PixelOff = color.White
var PixelOn = color.Black

type canvas struct {
	image *image.Paletted
}

var _ engine.Canvas = &canvas{}

func NewCanvas() *canvas {
	return &canvas{
		image: image.NewPaletted(
			image.Rect(0, 0, 84, 48),
			color.Palette{
				PixelOff,
				PixelOn,
			},
		),
	}
}

func (c *canvas) Clear() {
	idx := uint8(c.image.Palette.Index(PixelOff))
	for i := range c.image.Pix {
		c.image.Pix[i] = idx
	}
}

func (c *canvas) Set(x, y int, val bool) {
	if val {
		c.image.Set(x, y, PixelOn)
	} else {
		c.image.Set(x, y, PixelOff)
	}
}

func (c *canvas) Get(x, y int) bool {
	return c.image.At(x, y) == PixelOn
}

func (c *canvas) Width() int {
	return c.image.Rect.Size().X
}

func (c *canvas) Height() int {
	return c.image.Rect.Size().Y
}
