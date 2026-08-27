package draw

import "github.com/hoani/3310_engine/engine"

var DefaultSpriteOpts = &SpriteOpts{}

type SpriteAlign uint8

const (
	SaTopLeft SpriteAlign = iota
	SaTop
	SaTopRight
	SaLeft
	SaCenter
	SaRight
	SaBottomLeft
	SaBottom
	SaBottomRight
)

const SpriteAlignNum = SaBottomRight + 1

type Window struct {
	Apply      bool
	W, H, X, Y int
}

type Repeat struct {
	X uint8
	Y uint8
}

type SpriteOpts struct {
	Opts
	HFlip    bool
	VFlip    bool
	Rotation Rotation
	Window   Window
	Align    SpriteAlign
	Repeat   Repeat
}

func NewSpriteOpts() *SpriteOpts {
	return &SpriteOpts{Opts: *NewOpts(), HFlip: false, VFlip: false, Rotation: Rot0, Window: Window{}}
}

func (o *SpriteOpts) WithFlip(horizontal, vertical bool) *SpriteOpts {
	o.HFlip = horizontal
	o.VFlip = vertical
	return o
}

func (o *SpriteOpts) WithRotation(r Rotation) *SpriteOpts {
	o.Rotation = r
	return o
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

func (o *SpriteOpts) WithWindow(width, height int) *SpriteOpts {
	o.Window.Apply = true
	o.Window.W = width
	o.Window.H = height

	return o
}

func (o *SpriteOpts) WithAlign(align SpriteAlign) *SpriteOpts {
	o.Align = align
	return o
}

func (o *SpriteOpts) WithAlphaCustom(amount uint8, dither DitherFunc) *SpriteOpts {
	o.WithAlphaCustom(amount, dither)
	return o
}

func (o *SpriteOpts) WithRepeat(x, y uint8) *SpriteOpts {
	o.Repeat.X = x
	o.Repeat.Y = y
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

func (a SpriteAlign) Apply(x, y, w, h int) (int, int) {
	w += w % 2
	switch a {
	case SaTopLeft:
		return x, y
	case SaTop:
		return x - w/2, y
	case SaTopRight:
		return x - w, y
	case SaLeft:
		return x, y - h/2
	case SaCenter:
		return x - w/2, y - h/2
	case SaRight:
		return x - w, y - h/2
	case SaBottomLeft:
		return x, y - h
	case SaBottom:
		return x - w/2, y - h
	case SaBottomRight:
		return x - w, y - h
	}
	return x, y
}

func (d *draw) Sprite(x, y int, spr engine.Sprite, index int, opts *SpriteOpts) {
	if opts == nil {
		opts = DefaultSpriteOpts
	}

	w, h := spr.Width(), spr.Height()
	i0 := 0
	j0 := 0

	x, y = opts.Align.Apply(x, y, w*int(1+opts.Repeat.X), h*int(1+opts.Repeat.Y))

	if opts.Window.Apply {
		w, h = opts.Window.W, opts.Window.H
		i0 = opts.Window.X
		j0 = opts.Window.Y
	}

	i1 := i0 + w
	j1 := j0 + h

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

			sampleX, sampleY := opts.Transform(i, j, w, h)

			shade := spr.At(sampleX, sampleY, index)
			show, on := DitherSprite(x+i-i0, y+j-j0, shade)

			if show {
				if opts.Invert {
					on = !on
				}
				if !opts.Outline.Only {
					for n := range opts.Repeat.X + 1 {
						xpos := x + int(n)*w
						for m := range opts.Repeat.Y + 1 {
							ypos := y + int(m)*h
							d.c.Set(xpos+i-i0, ypos+j-j0, on)

						}
					}
				}
				continue
			}

			if applyOutline {
				if spr.At(sampleX-1, sampleY, index) != engine.SpriteTransparent ||
					spr.At(sampleX+1, sampleY, index) != engine.SpriteTransparent ||
					spr.At(sampleX, sampleY-1, index) != engine.SpriteTransparent ||
					spr.At(sampleX, sampleY+1, index) != engine.SpriteTransparent {
					d.c.Set(x+i-i0, y+j-j0, outlineInk)
				}
			}
		}
	}
}
