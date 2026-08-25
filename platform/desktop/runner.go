//go:build !tinygo

package desktop

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Runner interface {
	Setup(p *Platform)
	Run() error
}

type PreLaunch interface {
	ebiten.Game
	Setup()
	Done() bool
}

type Extension interface {
	Update() error
	Draw(screen *ebiten.Image, port *image.Rectangle)
}

type runner struct {
	prelaunch  PreLaunch
	platform   *Platform
	extensions []Extension
}

func NewRunner(prelaunch PreLaunch, extensions ...Extension) *runner {
	return &runner{prelaunch: prelaunch, extensions: extensions}
}

func (r *runner) Setup(p *Platform) {
	r.platform = p
}

func (r *runner) Run() error {
	r.prelaunch.Setup()
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

type defaultRunner struct{ p *Platform }

func NewDefaultRunner() *defaultRunner {
	return &defaultRunner{}
}

func (r *defaultRunner) Setup(p *Platform) {
	r.p = p
}

func (r *defaultRunner) Run() error {
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("Default Launcher")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(r.p.game.Info().Fps)

	ebiten.SetRunnableOnUnfocused(true)

	return ebiten.RunGame(r.p)
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

func (l *defaultPreLaunch) Setup() {
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle(l.name)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
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
