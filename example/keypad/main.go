package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sound"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/example/assets/fonts/cink"
	"github.com/hoani/3310_engine/example/assets/fonts/mwelch"
	"github.com/hoani/3310_engine/platform"
)

type Item struct {
	update func() error
	draw   func(canvas engine.Canvas) error
	name   string
}

type Game struct {
	draw   draw.Draw
	snd    engine.SoundPlayer
	debug  engine.Debug
	keypad command.Command[engine.Key]
	info   *engine.GameInfo
	index  int
	items  []Item
}

func (g *Game) Setup(p engine.Platform) {
	g.keypad = p.Cmd()
	g.debug = p.Debug()
	g.snd = p.Snd()
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {

	return g.items[g.index].update()
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	canvas.Clear(false)
	g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Font(&mwelch.Tiny).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) itemKeypad(name string) Item {
	snds := make([]engine.Sound, 0, 16)
	for i := range 16 {
		snds = append(snds, sound.Sound(10, sound.Note(note.C3+note.Index(i), 0xFF, 8)))
	}

	text := ""
	return Item{
		name: name,
		update: func() error {
			for i := range 16 {
				if g.keypad.Pressed(engine.Key(i)) {
					g.snd.Play(snds[i])
					text = engine.KeyName(engine.Key(i))
				}
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
		g.itemKeypad("keypad"),
	)

	platform.Run(g)
}
