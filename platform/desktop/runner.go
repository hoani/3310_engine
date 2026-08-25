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
	Draw(screen *ebiten.Image, port image.Rectangle)
}

type runner struct {
	prelaunch PreLaunch
	platform  *Platform
}

func NewRunner(prelaunch PreLaunch) Runner {
	return &runner{prelaunch: prelaunch}
}

func (r *runner) Setup(p *Platform) {
	r.platform = p
}

func (r *runner) Run() error {
	r.prelaunch.Setup()
	return ebiten.RunGame(r)
}

func (r *runner) Update() error {
	if r.prelaunch.Done() {
		return r.platform.Update()
	}
	return r.prelaunch.Update()
}

func (r *runner) Draw(screen *ebiten.Image) {
	if r.prelaunch.Done() {
		r.platform.Draw(screen)
		return
	}
	r.prelaunch.Draw(screen)
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
