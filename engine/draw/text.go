package draw

import "tinygo.org/x/tinyfont"

type TextBuilder struct {
	x           int
	y           int
	s           string
	defaultFont tinyfont.Fonter
	font        tinyfont.Fonter
	canvas      *FontCanvas
}

func NewTextBuilder(canvas *FontCanvas, defaultFont tinyfont.Fonter) *TextBuilder {
	return &TextBuilder{
		defaultFont: defaultFont,
		canvas:      canvas,
	}
}

func (t *TextBuilder) New(x, y int, s string) *TextBuilder {
	t.x = x
	t.y = y
	t.s = s
	t.font = t.defaultFont
	return t
}

func (t *TextBuilder) Font(f tinyfont.Fonter) *TextBuilder {
	t.font = f
	return t
}

func (t *TextBuilder) Draw(on bool) {
	t.canvas.SetInk(on)
	tinyfont.WriteLine(t.canvas, t.font, int16(t.x), int16(t.y), t.s, PixelOn)
}
