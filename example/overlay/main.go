package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sprite"
	"github.com/hoani/3310_engine/example/sprite/sprites/pgm"
	"github.com/hoani/3310_engine/platform"
)

type Item struct {
	step func() error
	draw func(canvas engine.Canvas) error
	name string
}

type Game struct {
	draw   draw.Draw
	debug  engine.Debug
	info   *engine.GameInfo
	keypad *command.Command[engine.Key]
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

	if g.keypad.Pressed(engine.KB) {
		g.info.Illuminated = !g.info.Illuminated
	}
	if g.keypad.Pressed(engine.KC) {
		g.index = (g.index + 1) % len(g.items)
	}
	if g.keypad.Pressed(engine.KD) {
		g.index = (g.index - 1)
		if g.index < 0 {
			g.index += len(g.items)
		}
	}
	return g.items[g.index].step()
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

	sopts := draw.NewSpriteOpts().WithAlign(draw.SaCenter)

	return func(c engine.Canvas) error {

		g.draw.Sprite(42, 24, sphere, 0, sopts)

		return nil
	}
}

func (g *Game) itemOverlay(name string, o Overlay) Item {
	back := g.drawSphere()
	return Item{
		name: name,
		step: func() error {
			o.Update()
			return nil
		},
		draw: func(canvas engine.Canvas) error {
			back(canvas)
			canvas.Overlay(o)
			return nil
		},
	}
}

func main() {

	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}

	g.items = append(
		g.items,
		g.itemOverlay("dissolve", NewDissolve()),
		g.itemOverlay("hwipe", NewHWipe()),
	)
	platform.Run(g)
}

type Overlay interface {
	engine.Overlay
	Update()
}

type OverlayDissolve struct {
	chunk [][]byte
	step  uint8
}

func NewDissolve() *OverlayDissolve {
	chunk := make([][]uint8, 8)
	for i := range chunk {
		chunk[i] = make([]uint8, 1)
	}
	return &OverlayDissolve{
		chunk: chunk,
		step:  0,
	}
}

func (o *OverlayDissolve) Update() {
	for i := range o.chunk {
		for j := range 8 {
			if draw.Dither(i, j, o.step) {
				o.chunk[i][0] |= 1 << j
			} else {
				o.chunk[i][0] &^= 1 << j
			}
		}
	}
	o.step += 1
}

func (o *OverlayDissolve) Get() [][]byte {
	return o.chunk
}
func (o *OverlayDissolve) Ink() bool {
	return true
}

type OverlayHWipe struct {
	chunk [][]byte
	step  uint8
	wipe  uint8
}

func NewHWipe() *OverlayHWipe {
	chunk := make([][]uint8, 1)
	chunk[0] = make([]uint8, 1)
	return &OverlayHWipe{
		chunk: chunk,
		step:  0,
	}
}

func (o *OverlayHWipe) Update() {
	o.chunk[0][0] = o.wipe

	o.step++
	if (o.step % 4) != 0 {
		return
	}

	if o.wipe == 0xff {
		o.wipe = 0x00
	} else {
		o.wipe = o.wipe<<1 + 1
	}
}

func (o *OverlayHWipe) Get() [][]byte {
	return o.chunk
}
func (o *OverlayHWipe) Ink() bool {
	return true
}
