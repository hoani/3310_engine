//go:build tinygo

package platform

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform/firmware"
)

func Run(game engine.Game) {
	firmware.Run(game)
}
