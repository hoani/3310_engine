//go:build !tinygo

package desktop

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
)

type Runner interface {
	Keypad() *command.CommandImpl[engine.Key]
	Setup(p *Platform)
	Run() error
}

type PreLaunch interface {
	ebiten.Game
	Setup() error
	Done() bool
}

type Extension interface {
	Setup() error
	Update() error
	Draw(screen *ebiten.Image, port *image.Rectangle)
}

type runner struct {
	prelaunch  PreLaunch
	platform   *Platform
	extensions []Extension
	keypad     *command.CommandImpl[engine.Key]
}

func NewRunner(prelaunch PreLaunch, extensions ...Extension) *runner {
	return &runner{prelaunch: prelaunch, extensions: extensions}
}

func (r *runner) WithKeypad(k *command.CommandImpl[engine.Key]) *runner {
	r.keypad = k
	return r
}

func (r *runner) Keypad() *command.CommandImpl[engine.Key] {
	if r.keypad == nil {
		r.keypad = NewKeypad()
	}
	return r.keypad
}

func (r *runner) Setup(p *Platform) {
	r.platform = p
}

func (r *runner) Run() error {
	if err := r.prelaunch.Setup(); err != nil {
		return err
	}

	for _, extension := range r.extensions {
		if err := extension.Setup(); err != nil {
			return err
		}
	}

	return ebiten.RunGame(r)
}

func (r *runner) Update() error {
	if !r.prelaunch.Done() {
		return r.prelaunch.Update()
	}

	if err := r.platform.Update(); err != nil {
		return err
	}

	for _, extension := range r.extensions {
		if err := extension.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (r *runner) Draw(screen *ebiten.Image) {
	if !r.prelaunch.Done() {
		r.prelaunch.Draw(screen)
		return
	}

	r.platform.Draw(screen)
	for _, extension := range r.extensions {
		extension.Draw(screen, r.platform.DrawPort())
	}
}

func (r *runner) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if r.prelaunch.Done() {
		return r.platform.Layout(outsideWidth, outsideHeight)
	}
	r.platform.Layout(outsideWidth, outsideHeight) // At least keep it up to date.
	return r.prelaunch.Layout(outsideWidth, outsideHeight)
}

func NewDefaultRunner(g engine.Game) *runner {
	return NewRunner(NewDefaultPreLaunch(840, 480, ""), NewDebugExtension(g))
}

type defaultPreLaunch struct {
	w, h int
	name string
}

func NewDefaultPreLaunch(w, h int, name string) PreLaunch {
	return &defaultPreLaunch{
		w: w, h: h, name: name,
	}
}

func (l *defaultPreLaunch) Setup() error {
	ebiten.SetWindowSize(l.w, l.h)
	ebiten.SetWindowTitle(l.name)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetFullscreen(true)
	return nil
}

func (l *defaultPreLaunch) Done() bool {
	return true
}

func (l *defaultPreLaunch) Update() error {
	return nil
}

func (l *defaultPreLaunch) Draw(screen *ebiten.Image) {
}

func (l *defaultPreLaunch) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
