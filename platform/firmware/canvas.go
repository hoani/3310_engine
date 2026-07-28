//go:build tinygo

package firmware

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/draw"
	"tinygo.org/x/drivers/pcd8544"
)

type canvas struct {
	device *pcd8544.Device
}

func NewCanvas(device *pcd8544.Device) engine.Canvas {
	return &canvas{
		device: device,
	}
}

func (c *canvas) Clear() {
	c.device.ClearBuffer()
}

func (c *canvas) Set(x, y int, val bool) {
	if val {
		c.device.SetPixel(int16(x), int16(y), draw.PixelOn)
	} else {
		c.device.SetPixel(int16(x), int16(y), draw.PixelOff)
	}
}

func (c *canvas) Get(x, y int) bool {
	return c.device.GetPixel(int16(x), int16(y))
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
	if col.R == 0 && col.G == 0 && col.B == 0 {
		col = draw.PixelOff
	} else {
		col = draw.PixelOn
	}
	c.device.SetPixel(x, y, col)
}

func (c *canvas) Display() error {
	return c.Display()
}
