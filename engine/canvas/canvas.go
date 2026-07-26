package canvas

import (
	"image"
	"image/color"
)

type Canvas interface {
	Image() *image.Paletted
	Clear()
	Size() image.Point
}

type canvas struct {
	image *image.Paletted
}

type CanvasBuilder struct {
	W, H    int
	Palette color.Palette
}

func New() *CanvasBuilder {
	return &CanvasBuilder{
		W: 84, H: 48,
		Palette: color.Palette{
			color.Black,
			color.White,
		},
	}
}

func (b *CanvasBuilder) Build() Canvas {
	return &canvas{
		image: image.NewPaletted(
			image.Rect(0, 0, b.W, b.H),
			b.Palette,
		),
	}
}

func (c *canvas) Size() image.Point {
	return c.image.Bounds().Size()
}

func (c *canvas) Clear() {
	idx := uint8(c.image.Palette.Index(color.Black))
	for i := range c.image.Pix {
		c.image.Pix[i] = idx
	}
}

func (c *canvas) Image() *image.Paletted {
	return c.image
}
