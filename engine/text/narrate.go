package text

import (
	"github.com/hoani/3310_engine/engine/draw"
	"tinygo.org/x/tinyfont"
)

type Narrate struct {
	x, y, w  int
	f        tinyfont.Fonter
	s        string
	subindex uint16
	index    uint16
	halign   draw.FontAlign
	valign   draw.FontAlign
	opts     *draw.Opts
}

func New(x, y, w int, f tinyfont.Fonter, s string) *Narrate {
	return &Narrate{
		x:        x,
		y:        y,
		w:        w,
		f:        f,
		s:        Fit(f, s, w),
		subindex: 0,
		index:    0,
		halign:   draw.FaLeft,
		valign:   draw.FaTop,
		opts:     draw.NewOpts(),
	}
}

func (n *Narrate) Reset(s string) *Narrate {
	n.s = Fit(n.f, s, n.w)
	n.index = 0
	n.subindex = 0
	return n
}

func (n *Narrate) WithHalign(align draw.FontAlign) *Narrate {
	n.halign = align
	return n
}

func (n *Narrate) WithValign(align draw.FontAlign) *Narrate {
	n.valign = align
	return n
}

func (n *Narrate) Opts() *draw.Opts {
	return n.opts
}

func (n *Narrate) Done() bool {
	return int(n.index) >= len(n.s)
}

func (n *Narrate) Update(incr uint8, frac uint8) {
	if n.Done() {
		return
	}

	n.subindex += uint16(frac)
	if n.subindex > 0xff {
		n.subindex -= 0xff
		n.index++
	}
	n.index += uint16(incr)
	if n.Done() {
		n.index = uint16(len(n.s))
	}
}

func (n *Narrate) Draw(d draw.Draw, ink bool) {
	d.Text(n.x, n.y, n.s[:n.index]).Font(n.f).HAlign(n.halign).VAlign(n.valign).Draw(ink, n.opts)
}
