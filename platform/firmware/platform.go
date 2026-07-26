package firmware

import (
	"image/color"
	"machine"
	"time"

	"github.com/ebitengine/gomobile/geom"
	"github.com/hoani/3310_engine/engine/canvas"
	"tinygo.org/x/drivers/pcd8544"
)

type Game struct {
	count  int
	canvas canvas.Canvas
	device *pcd8544.Device
	color  color.RGBA
	pos    geom.Point
	led    machine.Pin
}

func NewGame(lcd *pcd8544.Device, led machine.Pin) *Game {
	return &Game{
		count:  0,
		canvas: canvas.New().Build(),
		color:  color.RGBA{255, 255, 255, 255},
		device: lcd,
		led:    led,
	}
}

func (g *Game) Update() error {
	g.count++
	i := g.count % g.canvas.Size().X
	j := (g.count / g.canvas.Size().X) % g.canvas.Size().Y

	c := color.RGBA{255, 255, 255, 255}
	if g.count/(g.canvas.Size().X*g.canvas.Size().Y)%2 == 1 {
		c = color.RGBA{0, 0, 0, 255}
	}
	g.canvas.Image().Set(i, j, c)

	return nil
}

func ColorToRgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func (g *Game) Draw() error {
	for i := 0; i < g.canvas.Size().X; i++ {
		for j := 0; j < g.canvas.Size().Y; j++ {
			g.device.SetPixel(int16(i), int16(j), ColorToRgba(g.canvas.Image().At(i, j)))
		}
	}
	g.device.Display()
	return nil
}

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

	g := NewGame(d, led)

	period := time.Second / time.Duration(60)

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
		start := time.Now()
		g.Update()
		g.Draw()
		delta := time.Since(start)
		rem := period - delta
		time.Sleep(rem)
	}
}
