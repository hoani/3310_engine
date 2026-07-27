package engine

import (
	"image/color"
)

// Annoying generated fonts ignore colors when writing, so added this inverter

type FontCanvas struct {
	c   Canvas
	col color.RGBA
}

func (fc *FontCanvas) SetColor(col color.RGBA) {
	fc.col = col
}

func (fc *FontCanvas) Size() (x, y int16) {
	return fc.c.Size()
}

func (fc *FontCanvas) SetPixel(x, y int16, col color.RGBA) {
	if col.A == 0x00 {
		return
	}
	fc.c.SetPixel(x, y, fc.col)
}

func (fc *FontCanvas) Display() error {
	return nil
}
