//go:build tinygo && pca10059

package board

import (
	"machine"
)

func init() {
	Set(&Definition{
		Pcd: Pcd{
			Spi:    machine.SPI0,
			SckPin: machine.P1_10,
			SdoPin: machine.P1_13,
			SdiPin: machine.P1_11, // Unused and only available via test point on dongle
			DcPin:  machine.P1_15,
			RstPin: machine.P0_29,
			ScePin: machine.P0_02,
			LedPin: machine.P0_31,
		},
		Buzzer: Buzzer{
			Pwm: &PwmAdaptor{*machine.PWM0},
			Pin: machine.P0_10,
		},
		Keypad: Keypad{
			Col: [4]machine.Pin{machine.P0_13, machine.P0_15, machine.P0_17, machine.P0_20},
			Row: [4]machine.Pin{machine.P0_22, machine.P0_24, machine.P1_00, machine.P0_09},
		},
	})
}

type PwmAdaptor struct {
	machine.PWM
}

func (p *PwmAdaptor) Enable(enable bool) {
}
