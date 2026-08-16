package engine

type Canvas interface {
	Clear(set bool)
	Width() int
	Height() int
	Set(x, y int, val bool)
	Get(x, y int) bool
	Buffer() []byte
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
