//go:build tinygo && pico

package board

import "machine"

type PicoV0 struct {
	Base
	Backlight machine.Pin
}

// Based on https://github.com/hoani/3310_pico/releases/tag/v0
func V0() *PicoV0 {
	def := &PicoV0{
		Base:      *DefaultDefinition(),
		Backlight: machine.GPIO16,
	}
	def.Base.OnInitialize = func() {
		def.Backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})
		def.Backlight.High()
	}
	return def
}
