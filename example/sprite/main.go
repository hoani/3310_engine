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
	debug  engine.Debug
	col    bool
	config *engine.Config
	keypad command.Command[engine.Key]
	index  int
	items  []Item
}

func (g *Game) Setup(p engine.Platform) {
	p.Config(*g.config)
	g.keypad = p.Cmd()
	g.debug = p.Debug()
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

func (g *Game) drawIso() func(engine.Canvas) error {

	iso, err := sprite.FromP5(pgm.Iso)
	if err != nil {
		panic(err)
	}

	opts := draw.NewSpriteOpts().WithAlign(draw.SaCenter)

	return func(c engine.Canvas) error {

		g.draw.Sprite(42, 24, iso, 0, opts)
		g.draw.Sprite(42+iso.W, 24, iso, 0, opts)
		g.draw.Sprite(42+iso.W/2+1, 24+iso.H/4, iso, 0, opts)

		return nil
	}
}

func (g *Game) drawLetters() func(engine.Canvas) error {

	letters, err := sprite.StripFromP4(pgm.ClassicLight, pgm.ClassicLightMask, 7)
	if err != nil {
		panic(err)
	}

	return func(c engine.Canvas) error {
		textOps := draw.NewSpriteOpts()
		textOps.WithOutline(true).WithInvert()

		for i := range 12 {
			g.draw.Sprite(i*7, 12, letters, i, textOps)
			g.draw.Sprite(i*7, 24, letters, 12+i, textOps)
		}

		textOps.WithAlpha(uint8(g.count))
		for i := range 12 {
			g.draw.Sprite(i*7, 40, letters, i, textOps)
		}

		return nil
	}
}

func (g *Game) drawOutlines() func(engine.Canvas) error {

	letters, err := sprite.StripFromP4(pgm.ClassicLight, pgm.ClassicLightMask, 7)
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

func (g *Game) drawWindowed() func(engine.Canvas) error {

	scene, err := sprite.FromP4(pgm.Beeg, "")
	if err != nil {
		panic(err)
	}

	opts := draw.NewSpriteOpts().WithWindow(80, 44)

	// Using surface cut the cpu from ~78% to 38% - does go up a bit when panning etc.
	s := surface.New(80, 44)
	ds := draw.New(s)
	pending := true

	return func(c engine.Canvas) error {
		if g.keypad.Check(engine.K5) {
			opts.Window.Y--
			pending = true
		}
		if g.keypad.Check(engine.K7) {
			opts.Window.X--
			pending = true
		}
		if g.keypad.Check(engine.K0) {
			opts.Window.Y++
			pending = true
		}
		if g.keypad.Check(engine.K9) {
			opts.Window.X++
			pending = true
		}

		if pending {
			ds.Sprite(0, 0, scene, 0, opts)
			pending = false
		}
		c.DrawSurface(2, 2, s)

		return nil
	}
}

func (g *Game) drawAlignment() func(engine.Canvas) error {

	letter, err := sprite.StripFromP4(pgm.ClassicLight, pgm.ClassicLightMask, 7)
	if err != nil {
		panic(err)
	}

	opts := draw.NewSpriteOpts().WithOutline(true)

	names := []string{
		"Top Left",
		"Top",
		"Top Right",
		"Left",
		"Center",
		"Right",
		"Bottom Left",
		"Bottom",
		"Bottom Right",
	}

	lopts := draw.NewOpts()

	return func(c engine.Canvas) error {
		if g.keypad.Pressed(engine.K5) {
			opts.Align = (opts.Align + 1) % draw.SpriteAlignNum
		}
		g.draw.Line(42, 16, 42, 32).Draw(true, lopts)
		g.draw.Line(16, 24, 68, 24).Draw(true, lopts)

		g.draw.Sprite(42, 24, letter, 0, opts)

		g.draw.Text(42, 46, names[opts.Align]).HAlign(draw.FaCenter).VAlign(draw.FaBottom).Draw(true, nil)
		return nil
	}
}

func (g *Game) drawTiled() func(engine.Canvas) error {

	tile, err := sprite.FromP4NoMask(pgm.Tile)
	if err != nil {
		panic(err)
	}

	opts := draw.NewSpriteOpts().WithRepeat(1, 1)

	names := []string{
		"Top Left",
		"Top",
		"Top Right",
		"Left",
		"Center",
		"Right",
		"Bottom Left",
		"Bottom",
		"Bottom Right",
	}

	textOps := draw.NewOpts().WithOutline(false)

	return func(c engine.Canvas) error {
		if g.keypad.Pressed(engine.K5) {
			opts.Align = (opts.Align + 1) % draw.SpriteAlignNum
		}

		if g.keypad.Pressed(engine.K1) {
			opts.Repeat.X++
		}
		if g.keypad.Pressed(engine.K4) {
			opts.Repeat.X--
		}
		if g.keypad.Pressed(engine.K3) {
			opts.Repeat.Y++
		}
		if g.keypad.Pressed(engine.K6) {
			opts.Repeat.Y--
		}

		g.draw.Sprite(42, 24, tile, 0, opts)

		g.draw.Text(42, 46, names[opts.Align]).HAlign(draw.FaCenter).VAlign(draw.FaBottom).Draw(true, textOps)
		return nil
	}
}

func main() {

	g := &Game{config: &engine.Config{Debug: true, Fps: 60}}

	g.items = append(
		g.items,
		Item{draw: g.drawSphere(), name: "Sphere"},
		Item{draw: g.drawIso(), name: "Iso"},
		Item{draw: g.drawWindowed(), name: "Pan"},
		Item{draw: g.drawAlignment(), name: "Alignment"},
		Item{draw: g.drawOutlines(), name: "Outlines"},
		Item{draw: g.drawLetters(), name: "Fading"},
		Item{draw: g.drawArrow(), name: "Transforms"},
		Item{draw: g.drawTiled(), name: "Tiled"},
	)
	platform.Run(g)
}
