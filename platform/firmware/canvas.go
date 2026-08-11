//go:build tinygo

package firmware

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"tinygo.org/x/drivers/pcd8544"
)

type canvas struct {
	device *pcd8544.Device
	buffer []byte
	w, h   int
}

var _ engine.Canvas = &canvas{}

func NewCanvas(device *pcd8544.Device) *canvas {
	w, h := device.Size()
	return &canvas{
		device: device,
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
	w, _ := c.device.Size()
	return int(w)
}

func (c *canvas) Height() int {
	_, h := c.device.Size()
	return int(h)
}

func (c *canvas) Size() (x, y int16) {
	return c.device.Size()
}

func (c *canvas) SetPixel(x, y int16, col color.RGBA) {
	set := true
	if col.R != 0 || col.G != 0 || col.B != 0 {
		set = false
	}
	c.Set(int(x), int(y), set)
}

func (c *canvas) Display() error {
	if err := c.device.SetBuffer(c.buffer); err != nil {
		return err
	}
	return c.device.Display()
}
