package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sprite"
	"github.com/hoani/3310_engine/engine/surface"
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
	keypad command.Command[engine.Key]
	debug  engine.Debug
	sphere engine.Sprite
	info   *engine.GameInfo
	index  int
	items  []Item
}

func (g *Game) Setup(p engine.Platform) {
	g.keypad = p.Cmd()
	g.debug = p.Debug()
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

	if g.keypad.Pressed(engine.KC) {
		g.index = (g.index + 1) % len(g.items)
		g.count = 0
	}
	if g.keypad.Pressed(engine.KD) {
		g.index = (g.index - 1)
		if g.index < 0 {
			g.index += len(g.items)
		}
		g.count = 0
	}
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	canvas.Clear(false)

	g.items[g.index].draw(canvas)

	g.draw.Text(42, 0, g.items[g.index].name).HAlign(draw.FaCenter).Draw(false, draw.NewOpts().WithOutline(true))
	return nil
}

func (g *Game) drawStress(useSurface bool, spd int) func(canvas engine.Canvas) error {

	drawFunc := func(d draw.Draw) error {
		opts := draw.NewOpts()
		d.Rectangle(0, 0, 12, 16).DrawShade(0x88, opts)
		d.Rectangle(0, 0, 8, 12).Draw(true, opts)
		gradient := draw.NewLinearGradient(draw.P(16, 0), draw.P(48, 48), 0xff, 0x00)
		d.Rectangle(16, 0, 84, 48).DrawGradient(gradient, opts)

		d.Rectangle(24, 8, 76, 40).Draw(false, opts.WithAlphaCustom(55, draw.DitherOffset))
		return nil
	}

	if useSurface == false {
		return func(canvas engine.Canvas) error {
			canvas.Clear(false)
			return drawFunc(g.draw)
		}
	}

	surf := surface.New(84, 48)

	return func(canvas engine.Canvas) error {
		canvas.Clear(false)
		if useSurface {
			useSurface = false
			d := draw.New(surf)
			drawFunc(d)
		}

		xpos := (g.count/8)%84 - 42
		ypos := (g.count/8)%48 - 24
		canvas.DrawSurface(spd*xpos, spd*ypos, surf)
		return nil
	}
}

func main() {
	sphere, err := sprite.FromP5(pgm.Gradsphere)
	if err != nil {
		panic(err)
	}
	g := &Game{
		sphere: sphere,
		info:   &engine.GameInfo{Debug: true, Fps: 60},
	}
	g.items = append(
		g.items,
		Item{draw: g.drawStress(false, 0), name: "No Surface"},
		Item{draw: g.drawStress(true, 0), name: "With Surface"},
		Item{draw: g.drawStress(true, 1), name: "Moving"},
	)
	platform.Run(g)
}
