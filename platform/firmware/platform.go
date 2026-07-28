//go:build tinygo

package firmware

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"github.com/hoani/3310_engine/engine"
	"tinygo.org/x/drivers/pcd8544"
)

type Platform struct {
	game   engine.Game
	canvas engine.Canvas
	led    machine.Pin
	lcd    *pcd8544.Device
	keypad *Keypad
}

func New(game engine.Game, lcd *pcd8544.Device, led machine.Pin, keypad *Keypad) *Platform {
	return &Platform{
		game:   game,
		canvas: NewCanvas(lcd),
		led:    led,
		lcd:    lcd,
		keypad: keypad,
	}
}

func (p *Platform) Run() error {
	period := time.Second / time.Duration(p.game.Fps())
	for {
		start := time.Now()
		p.keypad.Update()
		if err := p.game.Update(); err != nil {
			return err
		}
		dStart := time.Now()
		if err := p.Draw(); err != nil {
			return err
		}
		dDur := time.Since(dStart)
		dur := time.Since(start)
		rem := period - dur
		fmt.Printf("cpu %d%%, draw %d%% \n", 100*rem/period, 100*dDur/dur)
		// fmt.Printf("cpu %d%%, draw %d%% \n%t %t %t\n%t %t %t\n%t %t %t\n%t %t %t\n\n", 100*rem/period, 100*dDur/dur,
		// 	p.keypad.Get(engine.K1), p.keypad.Get(engine.K2), p.keypad.Get(engine.K3),
		// 	p.keypad.Get(engine.K4), p.keypad.Get(engine.K5), p.keypad.Get(engine.K6),
		// 	p.keypad.Get(engine.K7), p.keypad.Get(engine.K8), p.keypad.Get(engine.K9),
		// 	p.keypad.Get(engine.KStar), p.keypad.Get(engine.K0), p.keypad.Get(engine.KHash),
		// )
		rem = period - time.Since(start)
		time.Sleep(rem)
	}
}

func (p *Platform) Draw() error {
	if err := p.game.Draw(p.canvas); err != nil {
		return err
	}
	return p.lcd.Display()
}

func Run(game engine.Game) {
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

	keypad, cmd := NewKeypad([3]machine.Pin{machine.GP3, machine.GP4, machine.GP5},
		[4]machine.Pin{machine.GP6, machine.GP7, machine.GP8, machine.GP9})

	game.Setup(cmd)

	p := New(game, d, led, keypad)

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

	err := p.Run()
	if err != nil {
		fmt.Printf("Game crash %e", err)
	}
}
