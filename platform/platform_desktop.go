//go:build !tinygo

package platform

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform/desktop"
)

func Run(game engine.Game) {
	desktop.Launch(game, desktop.NewDefaultLauncher())
}

func Launch(game engine.Game, launcher desktop.Launcher) {
	desktop.Launch(game, launcher)
}
