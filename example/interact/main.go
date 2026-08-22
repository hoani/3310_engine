package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/engine/sound"
	"github.com/hoani/3310_engine/engine/sound/note"
	"github.com/hoani/3310_engine/engine/sprite"
	"github.com/hoani/3310_engine/example/fonts/cink"
	"github.com/hoani/3310_engine/example/sprite/sprites/pgm"
)

type Item struct {
	update func() error
	draw   func(canvas engine.Canvas) error
	name   string
}

type Game struct {
	count  int
	draw   draw.Draw
	snd    engine.SoundPlayer
	debug  engine.Debug
	keypad *command.Command[engine.Key]
	info   *engine.GameInfo
	index  int
	items  []Item
}

func (g *Game) Setup(keypad *command.Command[engine.Key], snd engine.SoundPlayer, debug engine.Debug) {
	g.keypad = keypad
	g.debug = debug
	g.snd = snd
}

func (g *Game) Info() *engine.GameInfo {
	return g.info
}

func (g *Game) Update() error {
	g.count++

	g.keypad.Update()

	if g.keypad.Pressed(engine.KB) {
		g.info.Illuminated = !g.info.Illuminated
	}
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

	return g.items[g.index].update()
}

func (g *Game) Draw(canvas engine.Canvas) error {
	if g.draw == nil {
		g.draw = draw.New(canvas)
	}
	canvas.Clear(false)
	g.draw.Text(42, 1, g.items[g.index].name).HAlign(draw.FaCenter).Draw(true, nil)

	return g.items[g.index].draw(canvas)
}

func (g *Game) itemKeypad(name string) Item {
	snds := make([]engine.Sound, 0, 16)
	for i := range 16 {
		snds = append(snds, sound.Sound(10, sound.Note(note.C3+note.Index(i), 0xFF, 8)))
	}

	text := ""
	return Item{
		name: name,
		update: func() error {
			for i := range 16 {
				if g.keypad.Pressed(engine.Key(i)) {
					g.snd.Play(snds[i])
					text = engine.KeyName(engine.Key(i))
				}
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			g.draw.Text(42, 24, text).Font(&cink.Frogotype).HAlign(draw.FaCenter).VAlign(draw.FaMiddle).Draw(true, nil)
			return nil
		},
	}
}

func (g *Game) itemSphere() Item {

	sphere, err := sprite.FromP5(pgm.Gradsphere)
	if err != nil {
		panic(err)
	}

	snds := make([]engine.Sound, 0, 16)
	for i := range 16 {
		snds = append(snds, sound.Sound(10, sound.Note(note.C3+note.Index(i), 0xFF, 8)))
	}

	pos := 0
	idx := 0
	dir := 1

	return Item{
		name: "sphere",
		update: func() error {
			pos += (idx + 1) * dir

			if g.keypad.Pressed(engine.Key(engine.K2)) {
				idx = (idx + 1) % len(snds)
				g.snd.Play(snds[idx])

			}

			if pos > 84-sphere.W/2 && dir > 0 {
				dir = -dir
				idx = (idx + 1) % len(snds)
				g.snd.Play(snds[idx])
			}
			if pos < -sphere.W/2 && dir < 0 {
				dir = -dir
				idx = (idx + 1) % len(snds)
				g.snd.Play(snds[idx])
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			g.draw.Sprite(pos, 0, sphere, 0, draw.NewSpriteOpts())
			return nil
		},
	}
}

func (g *Game) itemDeadPixel() Item {

	xpos := 84 / 2
	ypos := 48 / 2

	xspd := 1
	yspd := 1

	idx := 0

	snds := make([]engine.Sound, 0, note.Total-1)
	for i := range cap(snds) {
		snds = append(snds, sound.Sound(10, sound.Note(note.Index(1+i), 0xFF, 8)))
	}

	ink := true

	return Item{
		name: "",
		update: func() error {
			xpos += xspd
			ypos += yspd

			if g.keypad.Pressed(engine.Key(engine.K2)) {
				ink = !ink
			}

			if (xpos > 84 && xspd > 0) || (xpos < 0 && xspd < 0) {
				xspd = -xspd
				idx = (idx + 1) % len(snds)
				g.snd.Play(snds[idx])
			}
			if (ypos > 48 && yspd > 0) || (ypos < 0 && yspd < 0) {
				yspd = -yspd
				idx = (idx + 1) % len(snds)
				g.snd.Play(snds[idx])
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			canvas.Clear(ink)
			g.draw.Text(xpos, ypos, "pix").Font(&cink.Frogotype).HAlign(draw.FaCenter).VAlign(draw.FaMiddle).Draw(!ink, nil)
			return nil
		},
	}
}

func (g *Game) itemBoxer() Item {

	boxer, err := sprite.StripFromP5(pgm.Boxing32, 32)
	if err != nil {
		panic(err)
	}

	snds := make([]engine.Sound, 0, note.Total-1)
	for i := range cap(snds) {
		snds = append(snds, sound.Sound(10, sound.Note(note.Index(1+i), 0xFF, 8)))
	}

	img := 0
	cooldown := 0

	opts := draw.NewSpriteOpts().WithAlign(draw.SaCenter)

	return Item{
		name: "",
		update: func() error {
			if cooldown > 0 {
				cooldown--
				if cooldown == 0 {
					img = 0
				}
			}
			if g.keypad.Pressed(engine.Key(engine.K1)) {
				img = 1
				cooldown = 60
			}
			if g.keypad.Pressed(engine.Key(engine.K2)) {
				img = 2
				cooldown = 60
			}
			if g.keypad.Pressed(engine.Key(engine.K3)) {
				img = 3
				cooldown = 60
			}

			return nil
		},
		draw: func(canvas engine.Canvas) error {
			g.draw.Sprite(84/2, 48/2, boxer, img, opts)
			return nil
		},
	}
}

func main() {
	g := &Game{info: &engine.GameInfo{Debug: true, Fps: 60}}
	g.items = append(
		g.items,
		g.itemKeypad("keypad"),
		g.itemSphere(),
		g.itemDeadPixel(),
		g.itemBoxer(),
	)

	Launch(g)
}
