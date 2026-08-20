package main

import (
	"fmt"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sound"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/example/fonts/cink"
	"github.com/hoani/3310_engine/example/snd/sound/music"
	"github.com/hoani/3310_engine/platform"
)

type Item struct {
	update func() error
	draw   func(canvas engine.Canvas) error
	name   string
}

type Game struct {
	count  int
	draw   draw.Draw
	snd    engine.SoundPlayer
	debug  engine.Debug
	keypad *command.Command[engine.Key]
	info   *engine.GameInfo
	index  int
	items  []Item
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.keypad = keypad
	g.debug = debug
	g.snd = snd
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

	g.keypad.Update()
	g.count++

	if g.keypad.Pressed(engine.K8) {
		g.index = (g.index + 1) % len(g.items)
		g.count = 0
	}
	if g.keypad.Pressed(engine.K7) {
		g.index = (g.index - 1)
		if g.index < 0 {
			g.index += len(g.items)
		}
		g.count = 0
	}

	return g.items[g.index].update()
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	canvas.Clear(false)
	g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) play(name string) Item {
	sfx := sound.Sound(10, sound.Note(note.C4, 0xFF, 8), sound.None(8), sound.Note(note.A4, 0xFF, 16), sound.None(8))
	tune := music.Song
	return Item{
		name: name,
		update: func() error {

			if g.keypad.Pressed(engine.K4) {
				g.snd.Play(sfx)
			}
			if g.keypad.Pressed(engine.K5) {
				g.snd.Track(tune, true)
			}
			if g.keypad.Pressed(engine.K6) {
				g.snd.Stop()
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error { return nil },
	}

}

func (g *Game) custom(name string) Item {

	freq := float32(500.0)
	tune := sound.Sound(10)
	text := fmt.Sprintf("%f", freq)
	return Item{
		name: name,
		update: func() error {

			if g.keypad.Pressed(engine.K4) {
				freq -= 10.0
				text = fmt.Sprintf("%f", freq)
			}
			if g.keypad.Pressed(engine.K5) {
				tune = sound.Sound(100, sound.Custom(freq, 0xFF, 2))
				g.snd.Play(tune)
			}
			if g.keypad.Pressed(engine.K6) {
				freq += 10.0
				text = fmt.Sprintf("%f", freq)
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			g.draw.Text(42, 24, text).Font(&cink.Frogotype).HAlign(draw.FaCenter).VAlign(draw.FaMiddle).Draw(true, nil)
			return nil
		},
	}

}

func main() {
	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		g.play("basic"),
		g.custom("custom"),
	)

	platform.Run(g)
}
