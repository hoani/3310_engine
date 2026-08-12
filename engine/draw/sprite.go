package draw

import "github.com/hoani/3310_engine/engine"

type SpriteOpts struct {
	Opts
	HFlip    bool
	VFlip    bool
	Rotation Rotation
}

func NewSpriteOpts() *SpriteOpts {
	return &SpriteOpts{Opts: *NewOpts(), HFlip: false, VFlip: false, Rotation: Rot0}
}

func (o *SpriteOpts) Transform(xi, yi, w, h int) (xo, yo int) {
	xr := 2*xi - (w - 1)
	yr := 2*yi - (h - 1)
	if o.HFlip {
		xr = -xr
	}
	if o.VFlip {
		yr = -yr
	}

	switch o.Rotation {
	case Rot0:
		break
	case Rot90:
		xr, yr = yr, -xr
	case Rot180:
		xr, yr = -xr, -yr
	case Rot270:
		xr, yr = -yr, xr
	}
	return (xr + (w - 1)) / 2, (yr + (h - 1)) / 2
}

func (d *draw) Sprite(x, y int, spr engine.Sprite, index int, opts *SpriteOpts) {

	i0 := 0
	i1 := spr.Width()
	j0 := 0
	j1 := spr.Height()

	outlineInk := opts.Outline.Ink
	if opts.Outline.Apply {
		i0 -= 1
		i1 += 1
		j0 -= 1
		j1 += 1
		if opts.Invert {
			outlineInk = !opts.Outline.Ink
		}
	}

	for i := i0; i < i1; i++ {
		for j := j0; j < j1; j++ {
			if !opts.Show(x+i, y+j) {
				continue
			}

			sampleX, sampleY := opts.Transform(i, j, spr.Width(), spr.Height())

			shade := spr.At(sampleX, sampleY, index)
			show, on := DitherSprite(x+i, y+j, shade)

			if show {
				if opts.Invert {
					on = !on
				}
				d.c.Set(x+i, y+j, on)
				continue
			}

			if opts.Outline.Apply {
				if spr.At(sampleX-1, sampleY, index) != engine.SpriteTransparent ||
					spr.At(sampleX+1, sampleY, index) != engine.SpriteTransparent ||
					spr.At(sampleX, sampleY-1, index) != engine.SpriteTransparent ||
					spr.At(sampleX, sampleY+1, index) != engine.SpriteTransparent {
					d.c.Set(x+i, y+j, outlineInk)
				}
			}
		}
	}
}
