package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/text"
	"github.com/hoani/3310_engine/example/assets/fonts/mwelch"
	"github.com/hoani/3310_engine/example/assets/fonts/somepx"
	"github.com/hoani/3310_engine/platform"
	"tinygo.org/x/tinyfont"
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
	keypad command.Command[engine.Key]
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
	g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) drawEffects(canvas engine.Canvas) error {
	c := true
	if g.count/(canvas.Width()*canvas.Height())%2 == 1 {
		c = false
	}

	g.draw.Text(52, 18, "Hello\nWorld").Font(&somepx.EffortsPro).Draw(c, draw.NewOpts())
	count := (4 * g.count) % 0x1FF
	if count > 0xFF {
		count = 0xFF
	}
	g.draw.Text(12, 18, "FADE").Font(&somepx.EffortsPro).Draw(c, draw.NewOpts().WithAlpha(uint8(count)))
	g.draw.Text(12, 32, "APPEAR").Font(&somepx.EffortsPro).Draw(c, draw.NewOpts().WithAlpha(uint8(0xFF-count)))

	g.draw.Text(12, 8, "OUTLINE").Font(&somepx.EffortsPro).Draw(!c, draw.NewOpts().WithOutline(c))

	g.draw.Text(48, 36, "Hello Tiny").Font(&mwelch.Tiny).Draw(true, draw.NewOpts())

	return nil
}

func (g *Game) drawAlignment(f tinyfont.Fonter, opts *draw.Opts) func(canvas engine.Canvas) error {

	return func(canvas engine.Canvas) error {
		str := ""
		index := (g.count / 120 % 9)
		switch index / 3 {
		case 0:
			str += "Left"
		case 1:
			str += "Center"
		case 2:
			str += "Right"
		}
		str += "\n"
		switch index % 3 {
		case 0:
			str += "Top"
		case 1:
			str += "Middle"
		case 2:
			str += "Bottom"
		}

		g.draw.Text(42, 24, str).Font(f).HAlign(draw.FontAlign(index/3)).VAlign(draw.FontAlign(index%3)).Draw(true, opts)

		return nil
	}
}

func (g *Game) drawNarrate() func(canvas engine.Canvas) error {
	f := &somepx.EffortsPro

	strings := []string{
		"I'll be your dream, I'll be your wish, I'll be your fantasy",
		"I'll be your hope, I'll be your love, be everything that you need",
		"I love you more with every breath truly, madly, deeply do",
		"I will be strong, I will be faithful 'cause I'm counting on",
	}

	index := 0

	n := text.New(4, 8, 84-8, f, strings[index])

	return func(canvas engine.Canvas) error {
		n.Update(1, 0)

		if n.Done() && g.keypad.Pressed(engine.K5) {
			index = (index + 1) % len(strings)
			n.Reset(strings[index])
		}
		n.Draw(g.draw, true)
		return nil
	}

}

func main() {
	g := &Game{config: &engine.Config{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		Item{draw: g.drawEffects, name: "Effects"},
		Item{draw: g.drawAlignment(&somepx.EffortsPro, draw.NewOpts()), name: "Align efforts"},
		Item{draw: g.drawAlignment(&somepx.EffortsPro, draw.NewOpts().WithOutline(false).WithInvert()), name: "Align efforts outlined"},
		Item{draw: g.drawAlignment(&mwelch.Tiny, draw.NewOpts()), name: "Align tiny"},
		Item{draw: g.drawAlignment(&mwelch.Tiny, draw.NewOpts().WithOutline(false).WithInvert()), name: "Align tiny outlined"},
		Item{draw: g.drawNarrate(), name: "Draw Narrate"},
	)

	platform.Run(g)
}
