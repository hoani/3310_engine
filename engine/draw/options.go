package draw

var DefaultOpts = &Opts{}

type OutlineOpts struct {
	Apply bool
	Ink   bool
	Only  bool
}

type DitherFunc func(x, y int, value uint8) bool

type AlphaOpts struct {
	Apply  bool
	Amount uint8
	Dither DitherFunc
}

type Opts struct {
	Invert  bool
	Outline OutlineOpts
	Alpha   AlphaOpts
}

func NewOpts() *Opts {
	return &Opts{}
}

func (o *Opts) WithInvert() *Opts {
	o.Invert = true
	return o
}

func (o *Opts) WithOutline(set bool) *Opts {
	o.Outline.Apply = true
	o.Outline.Ink = set
	return o
}

func (o *Opts) WithOutlineOnly() *Opts {
	o.Outline.Only = true
	return o
}

func (o *Opts) WithAlpha(amount uint8) *Opts {
	if amount == 0xFF {
		o.Alpha.Apply = false
		return o
	}
	o.Alpha.Apply = true
	o.Alpha.Amount = amount
	o.Alpha.Dither = Dither
	return o
}

func (o *Opts) WithAlphaCustom(amount uint8, dither DitherFunc) *Opts {
	if amount == 0xFF {
		o.Alpha.Apply = false
		return o
	}
	o.Alpha.Apply = true
	o.Alpha.Amount = amount
	o.Alpha.Dither = dither
	return o
}

func (o *Opts) Show(x, y int) bool {
	if !o.Alpha.Apply {
		return true
	}
	return o.Alpha.Dither(x, y, o.Alpha.Amount)
}
