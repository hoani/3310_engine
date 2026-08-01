package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/example/text/font"
	"github.com/hoani/3310_engine/platform"
)

type Game struct {
	count int
	draw  draw.Draw
	debug engine.Debug
	col   bool
	info  *engine.GameInfo
}

func (g *Game) Setup(keypad *command.Command[engine.Key], debug engine.Debug) {
	g.debug = debug
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

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
	platform.Run(&Game{info: &engine.GameInfo{Debug: true, Fps: 60}})
}
