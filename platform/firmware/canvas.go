//go:build tinygo

package firmware

import (
	"github.com/hoani/3310_engine/engine"
	"tinygo.org/x/drivers/pcd8544"
)

type canvas struct {
	engine.Canvas
	device *pcd8544.Device
}

var _ engine.Canvas = &canvas{}

func NewCanvas(device *pcd8544.Device) *canvas {
	w, h := device.Size()
	return &canvas{
		Canvas: engine.NewCanvas(int(w), int(h)),
		device: device,
	}
}

func (c *canvas) Display() error {
	if err := c.device.SetBuffer(c.Buffer()); err != nil {
		return err
	}
	return c.device.Display()
}
