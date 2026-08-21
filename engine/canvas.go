package engine

type Overlay interface {
	Get() [][]byte
	Ink() bool
}

type Canvas interface {
	Clear(set bool)
	Width() int
	Height() int
	Set(x, y int, val bool)
	Get(x, y int) bool
	Buffer() []byte
	Overlay(overlay Overlay)
	DrawSurface(x, y int, s Surface)
}

type Surface interface {
	Canvas
	Reset() // All pixels transparent.
	Ink() []byte
	Paper() []byte
}

type canvas struct {
	buffer []byte
	w, h   int
}

var _ Canvas = &canvas{}

func NewCanvas(w, h int) Canvas {
	return &canvas{
		buffer: make([]byte, w*h/8),
		w:      int(w),
		h:      int(h),
	}
}

func (c *canvas) Clear(set bool) {
	clear(c.buffer)
	if set {
		for i := range c.buffer {
			c.buffer[i] = 0xff
		}
	}
}

func (c *canvas) Set(x, y int, val bool) {
	if x < 0 || x >= c.w || y < 0 || y >= c.h {
		return
	}
	byteIndex := x + (y/8)*c.w
	if val {
		c.buffer[byteIndex] |= 1 << uint8(y%8)
	} else {
		c.buffer[byteIndex] &^= 1 << uint8(y%8)
	}
}

func (c *canvas) Get(x, y int) bool {
	if x < 0 || x >= c.w || y < 0 || y >= c.h {
		return false
	}
	byteIndex := x + (y/8)*c.w
	return (c.buffer[byteIndex] >> uint8(y%8) & 0x1) == 1
}

func (c *canvas) Width() int {
	return int(c.w)
}

func (c *canvas) Height() int {
	return int(c.h)
}

func (c *canvas) Buffer() []byte {
	return c.buffer
}

func (c *canvas) Overlay(o Overlay) {
	chunk := o.Get()
	cw := len(chunk)
	if cw == 0 {
		return
	}
	ch := len(chunk[0])
	if ch == 0 {
		return
	}
	if o.Ink() {
		for i := range c.buffer {
			x := (i % c.w) % cw
			y := (i / c.w) % ch
			c.buffer[i] |= chunk[x][y]
		}
	} else {
		for i := range c.buffer {
			x := (i % c.w) % cw
			y := (i / c.w) % ch
			c.buffer[i] &^= chunk[x][y]
		}
	}
}

func (c *canvas) DrawSurface(x, y int, s Surface) {
	yshift := y % 8
	if yshift < 0 {
		yshift += 8
	}
	start := x + (y/8)*c.w
	if y < 0 && yshift != 0 {
		start -= c.w
	}

	ink := s.Ink()
	paper := s.Paper()

	sw := s.Width()
	rowOffset := c.w - sw
	rowEnd := c.w - x
	rowStart := -x

	for i := range ink {
		// Avoid drawing beyond the row end.
		if (i%sw) < rowStart || (i%sw) > rowEnd {
			continue
		}
		idx := start + i + rowOffset*(i/sw)
		if idx < 0 {
			continue
		}
		if idx >= len(c.buffer) {
			break
		}
		c.buffer[idx] |= (ink[i] << yshift)
		c.buffer[idx] &^= (paper[i] << yshift)
	}
	if yshift != 0 {
		nyshift := 8 - yshift
		start += c.w
		for i := range ink {
			// Avoid drawing beyond the row end.
			if (i%sw) < rowStart || (i%sw) > rowEnd {
				continue
			}
			idx := start + i + rowOffset*(i/sw)
			if idx < 0 {
				continue
			}
			if idx >= len(c.buffer) {
				break
			}
			c.buffer[idx] |= (ink[i] >> nyshift)
			c.buffer[idx] &^= (paper[i] >> nyshift)
		}
	}
}
