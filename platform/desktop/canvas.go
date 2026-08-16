package desktop

import (
	"image"
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/draw"
)

type canvas struct {
	engine.Canvas
	image *image.Paletted
}

var _ engine.Canvas = &canvas{}

func NewCanvas() *canvas {
	return &canvas{
		Canvas: engine.NewCanvas(84, 48),
		image: image.NewPaletted(
			image.Rect(0, 0, 84, 48),
			color.Palette{
				draw.PixelOff,
				draw.PixelOn,
			},
		),
	}
}

func (c *canvas) Image() *image.Paletted {
	for i := range c.Width() {
		for j := range c.Height() {
			if c.Get(i, j) {
				c.image.Set(i, j, draw.PixelOn)
			} else {
				c.image.Set(i, j, draw.PixelOff)
			}
		}
	}
	return c.image
}
