package draw

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
)

// Annoying generated fonts ignore colors when writing, so added this inverter

type FontCanvas struct {
	c   engine.Canvas
	ink bool
}

func (fc *FontCanvas) SetInk(ink bool) {
	fc.ink = ink
}

func (fc *FontCanvas) Size() (x, y int16) {
	return int16(fc.c.Width()), int16(fc.c.Height())
}

func (fc *FontCanvas) SetPixel(x, y int16, col color.RGBA) {
	if col.A != 0 {
		fc.c.Set(int(x), int(y), fc.ink)
	}
}

func (fc *FontCanvas) Display() error {
	return nil
}
