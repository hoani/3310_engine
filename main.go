package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/font"
	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/sprites/pgm"
)

type Game struct {
	count  int
	draw   draw.Draw
	font   engine.Sprite
	keypad *command.Command[engine.Key]
	ypos   int
	xpos   int
	col    bool
}

func (g *Game) Setup(keypad *command.Command[engine.Key]) {
	g.keypad = keypad
}

func (g *Game) Fps() int {
	return 60
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
	}
	return nil
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}

	i := g.count % canvas.Width()
	j := (g.count / canvas.Width()) % canvas.Height()

	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}
	canvas.Set(i, j, c)

	// if g.count%120 == 0 {

	// 	alpha := math.Pi * float64(g.count) / float64(10*g.Fps())
	// 	s0a := math.Sin(alpha)
	// 	c0a := math.Cos(alpha)

	// 	s1a := math.Sin(alpha + math.Pi*2.0/3.0)
	// 	c1a := math.Cos(alpha + math.Pi*2.0/3.0)

	// 	s2a := math.Sin(alpha + math.Pi*4.0/3.0)
	// 	c2a := math.Cos(alpha + math.Pi*4.0/3.0)

	// 	d := 20.0

	// 	x0, y0 := canvas.Width()/2, canvas.Height()/2

	// 	p0 := engine.P(int(d*c0a)+x0, int(d*s0a)+y0)
	// 	p1 := engine.P(int(d*c1a)+x0, int(d*s1a)+y0)
	// 	p2 := engine.P(int(d*c2a)+x0, int(d*s2a)+y0)

	// 	g.draw.FillTriangle(p0, p1, p2, c)
	// }

	canvas.Clear()

	for i := 0; i < 26; i++ {
		g.draw.Sprite(i*7, 12, g.font, i, !c)
	}

	g.draw.Text(4, 28, &font.EffortsPro, "Hello World", !c)
	g.draw.Text(4, 6, &font.Tiny, "Hello World", true)
	g.draw.Text(g.xpos, g.ypos+32, &font.Tiny, "[0.0]", g.col)

	return nil
}

func main() {
	font, err := engine.StripFromP5(pgm.ClassicLight, 7)
	if err != nil {
		panic(err)
	}
	platform.Run(&Game{font: font})
}
