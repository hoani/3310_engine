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
	keypad *command.Command[engine.Key]
	debug  engine.Debug
	sphere engine.Sprite
	info   *engine.GameInfo
	index  int
	items  []Item
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

	if g.keypad.Pressed(engine.K8) {
		g.index = (g.index + 1) % len(g.items)
		g.count = 0
	}
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	canvas.Clear()

	g.draw.Text(8, 0, g.items[g.index].name)

	return g.items[g.index].draw(canvas)
}

func (g *Game) drawTriangle(canvas engine.Canvas) error {
	{
		count := g.count / 128 * 128

		alpha := math.Pi * float64(count) / float64(10*g.Info().Fps)
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

		if (g.count/128)%2 == 0 {
			g.draw.Triangle(p0, p1, p2).Draw(true)
		} else {
			shade := uint8((g.count % 128) << 1)
			g.draw.Triangle(p0, p1, p2).DrawShade(shade)
		}
	}
	return nil
}

func (g *Game) drawCircle(canvas engine.Canvas) error {
	canvas.Clear()

	g.draw.Circle(draw.P(84-16, 16), 13).Draw(true)
	g.draw.Circle(draw.P(84-16, 16), 5).Draw(false)

	gradient := draw.NewRadialGradient(draw.P(24, 16), 4, 32, 0xff, 0x00)
	g.draw.Circle(draw.P(32, 24), 46).DrawGradient(gradient)

	return nil
}

func (g *Game) drawRectangle(canvas engine.Canvas) error {
	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}
	canvas.Clear()

	g.draw.Rectangle(0, 0, 12, 16).DrawShade(0x88)
	g.draw.Rectangle(0, 0, 8, 12).Draw(c)
	gradient := draw.NewLinearGradient(draw.P(16, 0), draw.P(48, 48), 0xff, 0x00)
	g.draw.Rectangle(16, 0, 84, 48).DrawGradient(gradient)
	return nil
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
		Item{draw: g.drawTriangle, name: "Triangle"},
		Item{draw: g.drawCircle, name: "Circles"},
		Item{draw: g.drawRectangle, name: "Rectangle"},
	)
	platform.Run(g)
}
