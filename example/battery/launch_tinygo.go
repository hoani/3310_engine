//go:build tinygo

package main

import (
	"fmt"

	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/platform/firmware/board"
)

type hal struct {
	Board *board.PicoV1
}

func (h *hal) VoltageMv() uint16 {
	fmt.Println("sampling voltage")
	return uint16(2 * (uint32(h.Board.Battery.Get()) * 3300) / 65536)
}

func (h *hal) Backlight(r, g, b uint8, enable bool) {
	h.Board.Backlight.Enable(enable)
	if enable {
		h.Board.Backlight.Set(r, g, b)
	}
}

func Launch(g *Game) {
	b := board.V1()
	board.Set(b)

	g.Hal = &hal{Board: b}

	platform.Run(g)
}
