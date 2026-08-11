package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/example/text/font"
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
	keypad *command.Command[engine.Key]
	info   *engine.GameInfo
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

	g.draw.Text(52, 18, "Hello\nWorld").Font(&font.EffortsPro).Draw(c, draw.NewOpts())
	count := (4 * g.count) % 0x1FF
	if count > 0xFF {
		count = 0xFF
	}
	g.draw.Text(12, 18, "FADE").Font(&font.EffortsPro).Draw(c, draw.NewOpts().WithAlpha(uint8(count)))
	g.draw.Text(12, 32, "APPEAR").Font(&font.EffortsPro).Draw(c, draw.NewOpts().WithAlpha(uint8(0xFF-count)))

	g.draw.Text(12, 8, "OUTLINE").Font(&font.EffortsPro).Draw(!c, draw.NewOpts().WithOutline(c))

	g.draw.Text(48, 32, "Hello Tiny").Draw(true, draw.NewOpts())

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

func main() {
	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		Item{draw: g.drawEffects, name: "Effects"},
		Item{draw: g.drawAlignment(&font.EffortsPro, draw.NewOpts()), name: "Align efforts"},
		Item{draw: g.drawAlignment(&font.EffortsPro, draw.NewOpts().WithOutline(false).WithInvert()), name: "Align efforts outlined"},
		Item{draw: g.drawAlignment(&font.Tiny, draw.NewOpts()), name: "Align tiny"},
		Item{draw: g.drawAlignment(&font.Tiny, draw.NewOpts().WithOutline(false).WithInvert()), name: "Align tiny outlined"},
	)

	platform.Run(g)
}
