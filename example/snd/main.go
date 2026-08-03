package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sound"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/example/text/font"
	"github.com/hoani/3310_engine/platform"
)

type Game struct {
	count int
	draw  draw.Draw
	snd   engine.SoundPlayer
	debug engine.Debug
	col   bool
	info  *engine.GameInfo
	tune  engine.Sound
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.debug = debug
	g.snd = snd
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

	if (g.count % 300) == 60 {
		g.snd.Play(g.tune)
	}

	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}
	canvas.Clear()

	g.draw.Text(52, 28, "Hello\nWorld").Font(&font.EffortsPro).Draw(c)
	g.draw.Text(48, 6, "Hello Tiny").Draw(true)

	return nil
}

func main() {

	tune := sound.Sound(10, sound.Note(note.C1, 0xFF, 8), sound.None(8), sound.Note(note.C4, 0xFF, 16), sound.None(8))

	platform.Run(&Game{info: &engine.GameInfo{Debug: true, Fps: 60}, tune: tune})
}
