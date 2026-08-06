package draw

import "github.com/hoani/3310_engine/engine"

func (d *draw) Sprite(x, y int, spr engine.Sprite, index int, opts *Opts) {

	i0 := 0
	i1 := spr.Width()
	j0 := 0
	j1 := spr.Height()

	if opts.HasOutline {
		i0 -= 1
		i1 += 1
		j0 -= 1
		j1 += 1
	}

	for i := i0; i < i1; i++ {
		for j := j0; j < j1; j++ {
			shade := spr.At(i, j, index)
			show, on := DitherSprite(x+i, y+j, shade)

			if show {
				if opts.Invert {
					on = !on
				}
				d.c.Set(x+i, y+j, on)
				continue
			}

			if opts.HasOutline {
				if spr.At(i-1, j, index) != engine.SpriteTransparent ||
					spr.At(i+1, j, index) != engine.SpriteTransparent ||
					spr.At(i, j-1, index) != engine.SpriteTransparent ||
					spr.At(i, j+1, index) != engine.SpriteTransparent {
					d.c.Set(x+i, y+j, opts.OutlineInk)
				}
			}
		}
	}
}
