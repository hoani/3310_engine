package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/platform"
)

type Game struct {
	count  int
	draw   draw.Draw
	keypad *command.Command[engine.Key]
	debug  engine.Debug
	ypos   int
	xpos   int
	col    bool
	info   *engine.GameInfo
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.keypad = keypad
	g.debug = debug
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.keypad.Update()
	g.count++

	if g.keypad.Pressed(engine.K2) {
		g.ypos -= 4
	}
	if g.keypad.Pressed(engine.K8) {
		g.ypos += 4
	}
	if g.keypad.Pressed(engine.K4) {
		g.xpos -= 4
	}
	if g.keypad.Pressed(engine.K6) {
		g.xpos += 4
	}
	if g.keypad.Pressed(engine.K5) {
		g.col = !g.col
		g.debug.Console("Switched!")
	}
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}

	canvas.Clear()

	g.draw.Text(48, 6, "Hello Tiny").Draw(true)

	g.draw.Text(g.xpos, g.ypos+32, "[0.0]").Draw(g.col)

	return nil
}

func main() {
	platform.Run(&Game{info: &engine.GameInfo{Debug: true, Fps: 60}})
}
