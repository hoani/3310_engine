package firmware

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/pcd8544"
)

func Run() {
	// Configure SPI with a 1 MHz frequency.
	bus := machine.SPI0
	bus.Configure(machine.SPIConfig{
		Frequency: 1000000,
		SCK:       machine.Pin(18),
		SDO:       machine.Pin(19),
		SDI:       machine.Pin(16),
	})

	dcPin := machine.Pin(20)
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dcPin.High()

	rstPin := machine.Pin(21)
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin.High()

	scePin := machine.Pin(17)
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin.High() // CS inactive

	d := pcd8544.New(bus, dcPin, rstPin, scePin)

	d.Configure(pcd8544.Config{
		Width:  84,
		Height: 48,
	})

	d.ClearDisplay()

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	c := color.RGBA{255, 255, 255, 255}

	for j := range 48 {
		for i := range 84 {
			if j > 28 && i > 64 {
				c = color.RGBA{255, 255, 255, 255}
			} else {
				c = color.RGBA{0, 0, 0, 0}
			}
			d.SetPixel(int16(i), int16(j), c)
		}
	}

	for {

		// time.Sleep(100 * time.Millisecond)
		// led.Set(!led.Get())
		// d.Display()

		for j := range 48 {
			for i := range 84 {
				led.Set(!led.Get())
				d.SetPixel(int16(i), int16(j), c)
				time.Sleep(10 * time.Millisecond)
				d.Display()
			}
		}
		if c.A == 0 {
			c = color.RGBA{255, 255, 255, 255}
		} else {
			c = color.RGBA{0, 0, 0, 0}
		}
	}
}
