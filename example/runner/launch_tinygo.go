//go:build tinygo

package main

import (
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform"
)

func Launch(g engine.Game) {
	platform.Run(g)
}
