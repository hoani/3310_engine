//go:build !tinygo && !(js && wasm)

package main

import (
	"fmt"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/platform/desktop"
)

func Launch(g engine.Game) {
	fmt.Println("Launching desktop")
	platform.Launch(g, desktop.NewDefaultLauncher())
}
