package draw

import (
	"strings"

	"tinygo.org/x/tinyfont"
)

type FontAlign int

const (
	FaTop FontAlign = iota
	FaMiddle
	FaBottom
	FaLeft   = FaTop
	FaCenter = FaMiddle
	FaRight  = FaBottom
)

type TextBuilder struct {
	x           int
	y           int
	s           string
	defaultFont tinyfont.Fonter
	font        tinyfont.Fonter
	canvas      *FontCanvas
	halign      FontAlign
	valign      FontAlign
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
	t.halign = FaLeft
	t.valign = FaTop
	return t
}

func (t *TextBuilder) Font(f tinyfont.Fonter) *TextBuilder {
	t.font = f
	return t
}

func (t *TextBuilder) HAlign(fa FontAlign) *TextBuilder {
	t.halign = fa
	return t
}

func (t *TextBuilder) VAlign(fa FontAlign) *TextBuilder {
	t.valign = fa
	return t
}

func (t *TextBuilder) Draw(on bool, opts *Opts) {
	if opts == nil {
		opts = &DefaultOpts
	}
	t.canvas.SetOpts(opts)

	if opts.Outline.Apply {
		t.canvas.SetInk(opts.Outline.Ink)
		t.WriteLines(t.x-1, t.y)
		t.WriteLines(t.x+1, t.y)
		t.WriteLines(t.x, t.y-1)
		t.WriteLines(t.x, t.y+1)
		t.WriteLines(t.x-1, t.y-1)
		t.WriteLines(t.x+1, t.y+1)
		t.WriteLines(t.x+1, t.y-1)
		t.WriteLines(t.x-1, t.y+1)
	}
	t.canvas.SetInk(on)
	t.WriteLines(t.x, t.y)
}

func (t *TextBuilder) WriteLines(x, y int) {

	parts := strings.Split(t.s, "\n")

	hline := int(t.font.GetYAdvance())
	h := hline * len(parts)

	y += hline // Tiny font draws lines at bottom alignment

	switch t.valign {
	case FaTop:
		break
	case FaMiddle:
		y -= h / 2
	case FaBottom:
		y -= h
	}

	for _, part := range parts {
		wline, _ := tinyfont.LineWidth(t.font, part)
		xline := x
		switch t.halign {
		case FaLeft:
			break
		case FaCenter:
			xline -= int(wline / 2)
		case FaRight:
			xline -= int(wline)
		}

		tinyfont.WriteLine(t.canvas, t.font, int16(xline), int16(y), part, PixelOn)
		y += int(t.font.GetYAdvance())
	}

}

func LineWidth(f tinyfont.Fonter, s string) int {
	parts := strings.Split(s, "\n")
	w := 0
	for _, part := range parts {
		pw, _ := tinyfont.LineWidth(f, part)
		if int(pw) > w {
			w = int(pw)
		}
	}
	return w
}
