package firmware

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"tinygo.org/x/drivers/pcd8544"
)

var PixelOff = color.RGBA{0, 0, 0, 0}
var PixelOn = color.RGBA{255, 255, 255, 255}

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
		c.device.SetPixel(int16(x), int16(y), PixelOn)
	} else {
		c.device.SetPixel(int16(x), int16(y), PixelOff)
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
