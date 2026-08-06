package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sprite"
	"github.com/hoani/3310_engine/example/sprite/sprites/pgm"
	"github.com/hoani/3310_engine/platform"
)

type Game struct {
	count   int
	draw    draw.Draw
	debug   engine.Debug
	col     bool
	sphere  engine.Sprite
	letters engine.Sprite
	info    *engine.GameInfo
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.debug = debug
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++
	g.col = (g.count/300)%2 == 1
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}

	canvas.Clear()

	g.draw.Sprite((g.count)%128-40, 0, g.sphere, 0, draw.NewOpts())

	textOps := draw.NewOpts().WithOutline(true)

	for i := 0; i < 12; i++ {
		g.draw.Sprite(i*7, 12, g.letters, i, textOps)
	}

	for i := 0; i < 12; i++ {
		g.draw.Sprite(i*7, 32, g.letters, i, textOps.WithAlpha(uint8(g.count)))
	}

	return nil
}

func main() {
	sphere, err := sprite.FromP5(pgm.Gradsphere)
	if err != nil {
		panic(err)
	}
	letters, err := sprite.StripFromP5(pgm.ClassicLight, 7)
	if err != nil {
		panic(err)
	}
	platform.Run(&Game{sphere: sphere, letters: letters, info: &engine.GameInfo{Debug: true, Fps: 60}})
}
