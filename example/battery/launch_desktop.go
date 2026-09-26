//go:build !tinygo

package main

import (
	"fmt"

	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/platform/desktop"
)

type hal struct {
}

func (h *hal) VoltageMv() uint16 {
	return 3700
}

func (h *hal) Backlight(r, g, b uint8, enable bool) {
}

func Launch(g *Game) {
	fmt.Println("Launching desktop")

	g.Hal = &hal{}

	platform.Launch(g, desktop.NewRunner(desktop.NewDefaultPreLaunch(840, 640, "Battery"), desktop.NewDebugExtension(g)))
}
