package desktop

import (
	"image"
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/draw"
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
				draw.PixelOff,
				draw.PixelOn,
			},
		),
	}
}

func (c *canvas) Clear(set bool) {
	px := draw.PixelOff
	if set {
		px = draw.PixelOn
	}
	idx := uint8(c.image.Palette.Index(px))
	for i := range c.image.Pix {
		c.image.Pix[i] = idx
	}
}

func (c *canvas) Set(x, y int, val bool) {
	if val {
		c.image.Set(x, y, draw.PixelOn)
	} else {
		c.image.Set(x, y, draw.PixelOff)
	}
}

func (c *canvas) Get(x, y int) bool {
	return c.image.At(x, y) == draw.PixelOn
}

func (c *canvas) Width() int {
	return c.image.Rect.Size().X
}

func (c *canvas) Height() int {
	return c.image.Rect.Size().Y
}
