package main

import (
	"math"

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
	keypad command.Command[engine.Key]
	debug  engine.Debug
	sphere engine.Sprite
	config *engine.Config
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
	}
	if g.keypad.Pressed(engine.KD) {
		g.index = (g.index - 1)
		if g.index < 0 {
			g.index += len(g.items)
		}
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

func (g *Game) drawTriangle(canvas engine.Canvas) error {
	{
		count := g.count / 128 * 128

		alpha := math.Pi * float64(count) / float64(10*g.config.Fps)
		s0a := math.Sin(alpha)
		c0a := math.Cos(alpha)

		s1a := math.Sin(alpha + math.Pi*2.0/3.0)
		c1a := math.Cos(alpha + math.Pi*2.0/3.0)

		s2a := math.Sin(alpha + math.Pi*4.0/3.0)
		c2a := math.Cos(alpha + math.Pi*4.0/3.0)

		d := 20.0

		x0, y0 := canvas.Width()/2, canvas.Height()/2

		p0 := draw.P(int(d*c0a)+x0, int(d*s0a)+y0)
		p1 := draw.P(int(d*c1a)+x0, int(d*s1a)+y0)
		p2 := draw.P(int(d*c2a)+x0, int(d*s2a)+y0)

		opts := draw.NewOpts()

		if (g.count/128)%2 == 0 {
			g.draw.Triangle(p0, p1, p2).Draw(true, opts)
		} else {
			shade := uint8((g.count % 128) << 1)
			g.draw.Triangle(p0, p1, p2).DrawShade(shade, opts)
		}
	}
	return nil
}

func (g *Game) drawCircle(canvas engine.Canvas) error {
	opts := draw.NewOpts()

	g.draw.Circle(draw.P(84-16, 16), 13).Draw(true, opts)
	g.draw.Circle(draw.P(84-16, 16), 5).Draw(false, opts)

	gradient := draw.NewRadialGradient(draw.P(24, 16), 4, 32, 0xff, 0x00)
	g.draw.Circle(draw.P(32, 24), 46).DrawGradient(gradient, opts)

	return nil
}

func (g *Game) drawOval(canvas engine.Canvas) error {
	opts := draw.NewOpts()

	g.draw.Oval(draw.P(42, 24), 25, 20).Draw(true, opts)
	g.draw.Oval(draw.P(42, 24), 23, 5).Draw(false, opts)

	g.draw.Oval(draw.P(12, 12), 10, 8).Draw(true, opts)
	g.draw.Oval(draw.P(12, 24), 15, 12).Draw(true, opts)

	opts.WithOutlineOnly()
	g.draw.Oval(draw.P(72, 12), 10, 10).Draw(true, opts)
	g.draw.Circle(draw.P(72, 24), 10).Draw(true, opts)

	return nil
}

func (g *Game) drawRectangle(canvas engine.Canvas) error {
	c := true
	if g.count%300 > 150 {
		c = false
	}
	opts := draw.NewOpts()
	canvas.Clear(!c)

	g.draw.Rectangle(0, 0, 12, 16).DrawShade(0x88, opts)
	g.draw.Rectangle(0, 0, 8, 12).Draw(c, opts)
	gradient := draw.NewLinearGradient(draw.P(16, 0), draw.P(48, 48), 0xff, 0x00)
	g.draw.Rectangle(16, 0, 84, 48).DrawGradient(gradient, opts)

	g.draw.Rectangle(24, 8, 76, 40).Draw(false, opts.WithAlphaCustom(uint8(g.count), draw.DitherOffset))

	return nil
}

func (g *Game) drawFlashing(canvas engine.Canvas) error {
	c := true
	if g.count%2 == 0 {
		c = false
	}
	opts := draw.NewOpts()
	canvas.Clear(c)

	g.draw.Rectangle(0, 0, 16, 16).Draw(false, opts)
	g.draw.Rectangle(84-16, 0, 84, 16).Draw(true, opts)

	return nil
}

func (g *Game) drawLine(canvas engine.Canvas) error {
	opts := draw.NewOpts()

	for i := 20; i < 600; i += 300 {
		count := i + g.count

		alpha := math.Pi * float64(count) / float64(10*g.config.Fps)
		s0a := math.Sin(alpha)
		c0a := math.Cos(alpha)

		d := 17.0

		x0, y0 := canvas.Width()/2, canvas.Height()/2+3

		p0 := draw.P(int(d*c0a)+x0, int(d*s0a)+y0)
		p1 := draw.P(int(-d*c0a)+x0, int(-d*s0a)+y0)

		g.draw.Line(p0.X, p0.Y, p1.X, p1.Y).Draw(true, opts)
	}

	g.draw.Line(4, 16, 4, 32).Draw(true, opts)
	g.draw.Line(4, 32, 16, 44).Draw(true, opts)
	g.draw.Line(16, 44, 32, 44).Draw(true, opts)

	g.draw.Line(82, 4, 60, 44).Draw(true, opts)
	g.draw.Line(60-2, 44, 82-2, 4).Draw(true, opts)
	return nil
}

func (g *Game) drawShapes(opts *draw.Opts) func(canvas engine.Canvas) error {

	return func(canvas engine.Canvas) error {

		c := true
		if g.count%3 == 0 {
			c = false
		}
		canvas.Clear(c)

		g.draw.Rectangle(16, 16, 32, 32).Draw(true, opts)
		g.draw.Circle(draw.P(32, 24), 19).Draw(true, opts)
		g.draw.Triangle(draw.P(58, 8), draw.P(72, 44), draw.P(44, 44)).Draw(true, opts)
		g.draw.Oval(draw.P(8, 32), 9, 13).Draw(true, opts)
		g.draw.Oval(draw.P(8, 12), 13, 5).Draw(true, opts)

		return nil
	}
}

func (g *Game) drawTriangleFan(fill bool, opts *draw.Opts) func(canvas engine.Canvas) error {
	ts := draw.NewTriangleFan(
		draw.P(16, 24),
		draw.P(16, 16),
		draw.P(24, 24),
	)

	tso := draw.NewTriangleFan(
		draw.P(24+16, 24),
		draw.P(24+16, 16),
		draw.P(24+24, 24),
		draw.P(24+16, 32),
		draw.P(24+8, 24),
	)

	ts1 := draw.NewTriangleFan(
		draw.P(48+16, 24),
		draw.P(48+16, 16),
		draw.P(48+24, 24),
		draw.P(48+16, 32),
		draw.P(48+8, 24),
		draw.P(48+16, 16),
	)

	return func(canvas engine.Canvas) error {
		canvas.Clear(false)

		g.draw.Shape(ts).Draw(fill, opts)
		g.draw.Shape(tso).Draw(fill, opts)
		g.draw.Shape(ts1).Draw(fill, opts)

		return nil
	}
}

func (g *Game) drawTriangleStrip(fill bool, opts *draw.Opts) func(canvas engine.Canvas) error {
	ts := draw.NewTriangleStrip(
		draw.P(4, 34),
		draw.P(16, 44),
		draw.P(26, 36),
		draw.P(44, 44),
		draw.P(82, 34),
	)

	tso := draw.NewTriangleStrip(
		draw.P(32, 14),
		draw.P(32, 18),
		draw.P(48, 28),
		draw.P(40, 24),
		draw.P(16, 28),
		draw.P(24, 24),
		draw.P(32, 14),
		draw.P(32, 18),
	)

	ts1 := draw.NewTriangleStrip(draw.P(72, 28), draw.P(80, 24), draw.P(64, 16))

	return func(canvas engine.Canvas) error {
		canvas.Clear(false)

		g.draw.Shape(ts).Draw(fill, opts)
		g.draw.Shape(tso).Draw(fill, opts)
		g.draw.Shape(ts1).Draw(fill, opts)

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
		config: &engine.Config{Debug: true, Fps: 60},
	}
	g.items = append(
		g.items,
		Item{draw: g.drawOval, name: "Oval"},
		Item{draw: g.drawLine, name: "Line"},
		Item{draw: g.drawTriangleFan(true, draw.NewOpts()), name: "TriangleFan"},
		Item{draw: g.drawTriangleFan(true, draw.NewOpts().WithOutlineOnly()), name: "TriangleFan Outline"},
		Item{draw: g.drawTriangleStrip(true, draw.NewOpts()), name: "TriangleStrip"},
		Item{draw: g.drawTriangleStrip(true, draw.NewOpts().WithOutlineOnly()), name: "TriangleStrip Outline"},
		Item{draw: g.drawShapes(draw.NewOpts().WithOutline(false)), name: "With Outline"},
		Item{draw: g.drawShapes(draw.NewOpts().WithOutlineOnly()), name: "Outline Only"},
		Item{draw: g.drawShapes(draw.NewOpts()), name: "No Outlines"},
		Item{draw: g.drawTriangle, name: "Triangle"},
		Item{draw: g.drawCircle, name: "Circles"},
		Item{draw: g.drawRectangle, name: "Rectangle"},
		Item{draw: g.drawFlashing, name: "Flashing"},
	)
	platform.Run(g)
}
