package surface

import "github.com/hoani/3310_engine/engine"

type surface struct {
	w, h  int
	ink   engine.Canvas
	paper engine.Canvas
}

func roundUp8(n int) int {
	return (n + 7) &^ 7
}

func New(w, h int) engine.Surface {
	w = roundUp8(w)
	h = roundUp8(h)
	return &surface{
		w:     w,
		h:     h,
		ink:   engine.NewCanvas(w, h),
		paper: engine.NewCanvas(w, h),
	}
}

func (c *surface) Reset() {
	c.ink.Clear(false)
	c.paper.Clear(false)
}

// Probably want reset most of the time.
func (c *surface) Clear(set bool) {
	if set {
		c.ink.Clear(true)
		c.paper.Clear(false)
	} else {
		c.ink.Clear(false)
		c.paper.Clear(true)
	}
}

func (c *surface) Set(x, y int, val bool) {
	c.ink.Set(x, y, val)
	c.paper.Set(x, y, !val)
}

func (c *surface) Get(x, y int) bool {
	return c.ink.Get(x, y) // Kind of best guess here, false means transparent or paper
}

func (c *surface) Width() int {
	return int(c.w)
}

func (c *surface) Height() int {
	return int(c.h)
}

func (c *surface) Buffer() []byte {
	return c.ink.Buffer() // Again, best guess...
}

func (c *surface) Overlay(o engine.Overlay) {
	chunk := o.Get()
	cw := len(chunk)
	if cw == 0 {
		return
	}
	ch := len(chunk[0])
	if ch == 0 {
		return
	}
	buffAdd := c.ink.Buffer()
	buffSub := c.paper.Buffer()
	if !o.Ink() {
		buffAdd, buffSub = buffSub, buffAdd

	}
	for i := range buffAdd {
		x := (i % c.w) % cw
		y := (i / c.w) % ch
		buffAdd[i] |= chunk[x][y]
		buffSub[i] &^= chunk[x][y]
	}

}

func (c *surface) DrawSurface(x, y int, s engine.Surface) {
	c.ink.DrawSurface(x, y, s)
	sw := SwapSurface{s}
	c.paper.DrawSurface(x, y, &sw)
}

func (c *surface) Ink() []byte {
	return c.ink.Buffer()
}

func (c *surface) Paper() []byte {
	return c.paper.Buffer()
}

type SwapSurface struct {
	engine.Surface
}

func (c *SwapSurface) Ink() []byte {
	return c.Surface.Paper()
}

func (c *SwapSurface) Paper() []byte {
	return c.Surface.Ink()
}
