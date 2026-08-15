package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sprite"
	"github.com/hoani/3310_engine/example/sprite/sprites/pgm"
	"github.com/hoani/3310_engine/platform"
)

type Item struct {
	draw func(canvas engine.Canvas) error
	name string
}

type Game struct {
	count  int
	draw   draw.Draw
	debug  engine.Debug
	col    bool
	info   *engine.GameInfo
	keypad *command.Command[engine.Key]
	index  int
	items  []Item
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.debug = debug
	g.keypad = keypad
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.keypad.Update()
	g.count++

	if g.keypad.Pressed(engine.K8) {
		g.index = (g.index + 1) % len(g.items)
		g.count = 0
	}

	g.col = (g.count/300)%2 == 1
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}

	canvas.Clear(false)
	g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) drawSphere() func(engine.Canvas) error {

	sphere, err := sprite.FromP5(pgm.Gradsphere)
	if err != nil {
		panic(err)
	}

	return func(c engine.Canvas) error {

		g.draw.Sprite((g.count)%128-40, 0, sphere, 0, draw.NewSpriteOpts())

		return nil
	}
}

func (g *Game) drawLetters() func(engine.Canvas) error {

	letters, err := sprite.StripFromP5(pgm.ClassicLight, 7)
	if err != nil {
		panic(err)
	}

	return func(c engine.Canvas) error {
		textOps := draw.NewSpriteOpts()
		textOps.WithOutline(true)

		for i := range 12 {
			g.draw.Sprite(i*7, 12, letters, i, textOps)
		}

		textOps.WithAlpha(uint8(g.count))
		for i := range 12 {
			g.draw.Sprite(i*7, 24, letters, i, textOps)
		}

		return nil
	}
}

func (g *Game) drawOutlines() func(engine.Canvas) error {

	letters, err := sprite.StripFromP5(pgm.ClassicLight, 7)
	if err != nil {
		panic(err)
	}

	lineOps := draw.NewOpts()
	noneOps := draw.NewSpriteOpts().WithInvert()
	outlineOps := draw.NewSpriteOpts().WithOutline(true)
	outlineOnlyOps := draw.NewSpriteOpts().WithOutlineOnly().WithInvert()

	return func(c engine.Canvas) error {

		g.draw.Line(0, 16, 84, 16).Draw(true, lineOps)
		for i := range 10 {
			g.draw.Sprite(7+i*7, 12, letters, i, noneOps)
		}

		g.draw.Line(0, 28, 84, 28).Draw(true, lineOps)
		for i := range 10 {
			g.draw.Sprite(7+i*7, 24, letters, i, outlineOps)
		}

		g.draw.Line(0, 40, 84, 40).Draw(true, lineOps)
		for i := range 10 {
			g.draw.Sprite(7+i*7, 36, letters, i, outlineOnlyOps)
		}

		return nil
	}
}

func (g *Game) drawArrow() func(engine.Canvas) error {

	arrow, err := sprite.FromP5(pgm.Arrow)
	if err != nil {
		panic(err)
	}

	opts := draw.NewSpriteOpts()

	return func(c engine.Canvas) error {
		if g.keypad.Pressed(engine.K3) {
			opts.HFlip = !opts.HFlip
		}
		if g.keypad.Pressed(engine.K1) {
			opts.VFlip = !opts.VFlip
		}
		if g.keypad.Pressed(engine.K2) {
			opts.Rotation++
			if opts.Rotation > 3 {
				opts.Rotation = 0
			}
		}
		g.draw.Sprite(42-24, 0, arrow, 0, opts)

		return nil
	}
}

func main() {

	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}

	g.items = append(
		g.items,
		Item{draw: g.drawSphere(), name: "Sphere"},
		Item{draw: g.drawOutlines(), name: "Outlines"},
		Item{draw: g.drawLetters(), name: "Fading"},
		Item{draw: g.drawArrow(), name: "Transforms"},
	)
	platform.Run(g)
}
