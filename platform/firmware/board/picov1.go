//go:build tinygo && pico

package board

import "machine"

type PicoV1 struct {
	Base
	Backlight RgbBacklight
	Battery   machine.ADC
}

// Based on https://github.com/hoani/3310_pico/releases/tag/v1
func V1() *PicoV1 {
	def := &PicoV1{
		Base: *DefaultDefinition(),

		Backlight: RgbBacklight{
			// Note, we have PWM0 -> GPIO 0|1 or 16|17, PWM1 -> GPIO 2|3 or 18|19 and so on... up to PWM7
			R: LedPwm{Pwm: Pwm{Pwm: machine.PWM6, Pin: machine.GP13}},
			G: LedPwm{Pwm: Pwm{Pwm: machine.PWM6, Pin: machine.GP12}},
			B: LedPwm{Pwm: Pwm{Pwm: machine.PWM5, Pin: machine.GP11}},
		},
		Battery: machine.ADC{Pin: machine.ADC0},
	}
	def.Base.OnInitialize = func() {
		machine.GP16.Configure(machine.PinConfig{Mode: machine.PinOutput})
		machine.GP16.Set(true)

		machine.InitADC()
		def.Battery.Pin.Configure(machine.PinConfig{Mode: machine.PinAnalog})
		def.Battery.Configure(machine.ADCConfig{})

		def.Backlight.Setup()
		def.Backlight.Enable(true)
		def.Backlight.Set(0x88, 0x44, 0x88)
	}
	return def
}
