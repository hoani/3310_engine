package main

import (
	"fmt"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sound"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/example/assets/fonts/cink"
)

type Hal interface {
	VoltageMv() uint16
	Backlight(r, g, b uint8, enable bool)
}

type Item struct {
	update func() error
	draw   func(canvas engine.Canvas) error
	name   string
}

type Game struct {
	draw     draw.Draw
	snd      engine.SoundPlayer
	debug    engine.Debug
	keypad   command.Command[engine.Key]
	config   *engine.Config
	platform engine.Platform
	index    int
	items    []Item
	Hal      Hal
}

func (g *Game) Setup(p engine.Platform) {
	p.Config(*g.config)
	g.platform = p
	g.keypad = p.Cmd()
	g.debug = p.Debug()
	g.snd = p.Snd()
}

func (g *Game) Update() error {
	if g.keypad.Pressed(engine.KC) {
		g.index = (g.index + 1) % len(g.items)
	}
	if g.keypad.Pressed(engine.KD) {
		g.index = (g.index - 1)
		if g.index < 0 {
			g.index += len(g.items)
		}
	}

	return g.items[g.index].update()
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}

	return g.items[g.index].draw(canvas)
}

func mVToPercent(mV uint16) float32 {
	if mV > 4200 {
		return 100.0
	}
	if mV > 4000 {
		delta := mV - 4000
		return 80.0 + 20.0*float32(delta)/200.0
	}
	if mV > 3700 {
		delta := mV - 3700
		return 30.0 + 50.0*float32(delta)/300.0
	}
	if mV > 3200 {
		delta := mV - 3200
		return 0.0 + 30.0*float32(delta)/500.0
	}
	return 0.0
}

func (g *Game) itemBattery() Item {

	snds := make([]engine.Sound, 0, 16)
	for i := range 16 {
		snds = append(snds, sound.Sound(10, sound.Note(note.C3+note.Index(i), 0xFF, 8)))
	}

	lightOff := 5 * 50
	sleep := 20 * 60

	timeout := 0
	step := 0
	batteryMv := uint16(0)
	percent := float32(0.0)
	return Item{
		name: "battery",
		update: func() error {
			for i := range 16 {
				if g.keypad.Pressed(engine.Key(i)) {
					g.snd.Play(snds[i])
					timeout = 0
					g.platform.Display().Enable(true)
					g.config.Fps = 60
					g.platform.Config(*g.config)
					hue := uint32(0xFF) << i
					g.Hal.Backlight(uint8(hue>>16), uint8(hue>>8), uint8(hue), true)
				}
			}

			timeout++
			step++
			if (step % 15) == 0 {
				if batteryMv == 0 {
					batteryMv = g.Hal.VoltageMv()
				} else {
					batteryMv = uint16((uint32(g.Hal.VoltageMv()) + uint32(15*batteryMv)) / 16)
				}
				percent = mVToPercent(batteryMv)
			}

			if timeout > lightOff {
				g.Hal.Backlight(0, 0, 0, true)
				g.Hal.Backlight(0, 0, 0, false)
			}

			if timeout > sleep {
				// Go into a deep sleep, only updating step twice a second to check for button presses
				if g.config.Fps != 2 {
					g.config.Fps = 2
					g.platform.Config(*g.config)
					g.platform.Display().Enable(false)
				}
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			canvas.Clear(false)

			g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Draw(true, nil)

			g.draw.Text(42, 12, fmt.Sprintf("%d mV", batteryMv)).Font(&cink.Frogotype).HAlign(draw.FaCenter).VAlign(draw.FaMiddle).Draw(true, nil)
			g.draw.Text(42, 32, fmt.Sprintf("%.1f %%", percent)).Font(&cink.Frogotype).HAlign(draw.FaCenter).VAlign(draw.FaMiddle).Draw(true, nil)

			return nil
		},
	}
}

func main() {
	g := &Game{config: &engine.Config{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		g.itemBattery(),
	)

	Launch(g)
}
