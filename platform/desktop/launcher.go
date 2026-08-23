//go:build !tinygo

package desktop

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Launcher interface {
	Setup(p *Platform)
	Run() error
}

type PreLaunch interface {
	ebiten.Game
	Setup()
	Done() bool
}

type launcher struct {
	prelaunch PreLaunch
	platform  *Platform
}

func NewLauncher(prelaunch PreLaunch) Launcher {
	return &launcher{prelaunch: prelaunch}
}

func (l *launcher) Setup(p *Platform) {
	l.platform = p
}

func (l *launcher) Run() error {
	l.prelaunch.Setup()
	return ebiten.RunGame(l)
}

func (l *launcher) Update() error {
	if l.prelaunch.Done() {
		return l.platform.Update()
	}
	return l.prelaunch.Update()
}

func (l *launcher) Draw(screen *ebiten.Image) {
	if l.prelaunch.Done() {
		l.platform.Draw(screen)
		return
	}
	l.prelaunch.Draw(screen)
}

func (l *launcher) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if l.prelaunch.Done() {
		return l.platform.Layout(outsideWidth, outsideHeight)
	}
	l.platform.Layout(outsideWidth, outsideHeight) // At least keep it up to date.
	return l.prelaunch.Layout(outsideWidth, outsideHeight)
}

type defaultLauncher struct{ p *Platform }

func NewDefaultLauncher() *defaultLauncher {
	return &defaultLauncher{}
}

func (l *defaultLauncher) Setup(p *Platform) {
	l.p = p
}

func (l *defaultLauncher) Run() error {
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("Default Launcher")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(l.p.game.Info().Fps)

	ebiten.SetRunnableOnUnfocused(true)

	return ebiten.RunGame(l.p)
}
