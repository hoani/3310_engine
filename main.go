package main

import (
	"math"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/sprites/pgm"
)

type Game struct {
	count int
	draw  engine.Draw
	font  engine.Sprite
}

func (g *Game) Fps() int {
	return 60
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = engine.NewDraw(canvas)
	}

	i := g.count % canvas.Width()
	j := (g.count / canvas.Width()) % canvas.Height()

	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}
	canvas.Set(i, j, c)

	if g.count%120 == 0 {

		alpha := math.Pi * float64(g.count) / float64(10*g.Fps())
		s0a := math.Sin(alpha)
		c0a := math.Cos(alpha)

		s1a := math.Sin(alpha + math.Pi*2.0/3.0)
		c1a := math.Cos(alpha + math.Pi*2.0/3.0)

		s2a := math.Sin(alpha + math.Pi*4.0/3.0)
		c2a := math.Cos(alpha + math.Pi*4.0/3.0)

		d := 20.0

		x0, y0 := canvas.Width()/2, canvas.Height()/2

		p0 := engine.P(int(d*c0a)+x0, int(d*s0a)+y0)
		p1 := engine.P(int(d*c1a)+x0, int(d*s1a)+y0)
		p2 := engine.P(int(d*c2a)+x0, int(d*s2a)+y0)

		g.draw.FillTriangle(p0, p1, p2, c)
	}

	for i := 0; i < 26; i++ {
		g.draw.Sprite(i*7, 12, g.font, i, !c)
	}

	return nil
}

func main() {
	font, err := engine.StripFromP5(pgm.ClassicLight, 7)
	if err != nil {
		panic(err)
	}
	platform.Run(&Game{font: font})
}
