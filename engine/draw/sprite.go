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

func (o *SpriteOpts) WithInvert() *SpriteOpts {
	o.Opts.WithInvert()
	return o
}

func (o *SpriteOpts) WithOutline(set bool) *SpriteOpts {
	o.Opts.WithOutline(set)
	return o
}

func (o *SpriteOpts) WithOutlineOnly() *SpriteOpts {
	o.Opts.WithOutlineOnly()
	return o
}

func (o *SpriteOpts) WithAlpha(amount uint8) *SpriteOpts {
	o.Opts.WithAlpha(amount)
	return o
}

func (o *SpriteOpts) WithAlphaCustom(amount uint8, dither DitherFunc) *SpriteOpts {
	o.WithAlphaCustom(amount, dither)
	return o
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

	applyOutline := opts.Outline.Apply || opts.Outline.Only
	outlineInk := opts.Outline.Ink
	if applyOutline {
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
				if !opts.Outline.Only {
					d.c.Set(x+i, y+j, on)
				}
				continue
			}

			if applyOutline {
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
