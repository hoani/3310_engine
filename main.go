package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform"
)

type Game struct {
	count int
}

func (g *Game) Fps() int {
	return 60
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	i := g.count % canvas.Width()
	j := (g.count / canvas.Width()) % canvas.Height()

	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}
	canvas.Set(i, j, c)
	return nil
}

func main() {
	platform.Run(&Game{})
}
