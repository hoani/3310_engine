package draw

import "github.com/hoani/3310_engine/engine"

func (d *draw) Sprite(x, y int, spr engine.Sprite, index int, invert bool) {
	for i := 0; i < spr.Width(); i++ {
		for j := 0; j < spr.Height(); j++ {
			shade := spr.At(i, j, index)
			show, on := DitherSprite(x+i, y+j, shade)
			if show {
				if invert {
					on = !on
				}

				d.c.Set(x+i, y+j, on)
			}
		}
	}
}
