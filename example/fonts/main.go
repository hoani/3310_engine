package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/text"
	"github.com/hoani/3310_engine/example/assets/fonts/cink"
	"github.com/hoani/3310_engine/example/assets/fonts/mwelch"
	"github.com/hoani/3310_engine/example/assets/fonts/somepx"
	"github.com/hoani/3310_engine/platform"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/notosans"
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
	info   *engine.GameInfo
	index  int
	items  []Item
}

func (g *Game) Setup(keypad command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.debug = debug
	g.keypad = keypad
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

	if g.keypad.Pressed(engine.K8) {
		g.index = (g.index + 1) % len(g.items)
		g.count = 0
	}
	if g.keypad.Pressed(engine.K7) {
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
	g.draw.Text(42, 1, g.items[g.index].name).Font(&mwelch.Tiny).VAlign(draw.FaTop).HAlign(draw.FaCenter).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) drawFont(f tinyfont.Fonter) func(canvas engine.Canvas) error {
	strings := []string{
		"The quick brown fox jumps over the lazy dog",
		"0123456789 0.0 123 456 789",
		"!@#$ %^& *()_ +-= {}[] ;:' \",. /<>? `~",
		"ABCD EFG HIJK LMN OPQR STU VWX YZ",
		"abcd efg hijk lmn opqr stu vwx yz",
		"She \"sells\" sea shells (at) the sea shore",
	}

	index := 0

	n := text.New(0, 8, 84, f, strings[index])

	return func(canvas engine.Canvas) error {
		n.Update(1, 0)

		if n.Done() && g.keypad.Pressed(engine.K5) {
			index = (index + 1) % len(strings)
			n.Reset(strings[index])
		}
		if n.Done() && g.keypad.Pressed(engine.K4) {
			index = (index - 1)
			if index < 0 {
				index += len(strings)
			}
			n.Reset(strings[index])
		}
		n.Draw(g.draw, true)
		return nil
	}

}

func main() {
	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		// Chequered Ink - see https://chequered.ink/ for license details
		Item{draw: g.drawFont(&cink.BittypixMonospace), name: "CI: BittypixMonospace"},
		Item{draw: g.drawFont(&cink.CodersCrux), name: "CI: CodersCrux"},
		Item{draw: g.drawFont(&cink.DiaryOfAn8bitMage), name: "CI: DiaryOfAn8BitMage"},
		Item{draw: g.drawFont(&cink.Frogotype), name: "CI: Frogotype"},
		Item{draw: g.drawFont(&cink.MesseDuesseldorf), name: "CI: MesseDuesseldorf"},
		Item{draw: g.drawFont(&cink.NineteenEightySeven), name: "CI: NineteenEightySeven"},
		Item{draw: g.drawFont(&cink.NineteenNinetySeven), name: "CI: NineteenNinetySeven"},
		Item{draw: g.drawFont(&cink.NineteenNinetySix), name: "CI: NineteenNinetySix"},
		Item{draw: g.drawFont(&cink.NineteenNinetyThree), name: "CI: NineteenNinetyThree"},
		Item{draw: g.drawFont(&cink.Notalot18), name: "CI: Notalot18"},
		Item{draw: g.drawFont(&cink.Notalot25), name: "CI: Notalot25"},
		Item{draw: g.drawFont(&cink.PixelOrGTFO), name: "CI: Pixel or GTFO"},
		Item{draw: g.drawFont(&cink.PrincessSavesYou), name: "CI: PrincessSavesYou"},
		Item{draw: g.drawFont(&cink.Rygarde), name: "CInk: Rygarde"},
		Item{draw: g.drawFont(&cink.SuperLegendBoy), name: "CInk: Super Legend Boy"},
		Item{draw: g.drawFont(&cink.TinyAndChunky), name: "CInk: Tiny And Chunky"},
		Item{draw: g.drawFont(&cink.TeenyTinyPixls), name: "CInk: Teeny Tiny Pixls"},
		// Nokia 3310 Jam fonts - see individual licesnses in examples/assets/fonts
		Item{draw: g.drawFont(&somepx.EffortsPro), name: "EffortsPro"},
		Item{draw: g.drawFont(&mwelch.Tiny), name: "Tiny"},
		// Tiny font options - note licenses
		Item{draw: g.drawFont(&notosans.Notosans12pt), name: "tinyfont: Notosans"},       //SIL Open Font License.
		Item{draw: g.drawFont(&tinyfont.Org01), name: "tinyfont: Org01"},                 //BSD 3-Clause License.
		Item{draw: g.drawFont(&tinyfont.Picopixel), name: "tinyfont: Picopixel"},         //BSD 3-Clause License.
		Item{draw: g.drawFont(&tinyfont.Tiny3x3a2pt7b), name: "tinyfont: Tiny3x3a2pt7b"}, //Licensed under CC BY-NC-SA 3.0.
		Item{draw: g.drawFont(&tinyfont.TomThumb), name: "tinyfont: TomThumb"},           //BSD 3-Clause License.
	)

	platform.Run(g)
}
