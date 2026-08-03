//go:build tinygo && pico

package board

import "machine"

func init() {
	Set(&Definition{
		Pcd: Pcd{
			Spi:    machine.SPI0,
			SckPin: machine.GPIO18,
			SdoPin: machine.GPIO19,
			SdiPin: machine.GPIO16,
			DcPin:  machine.GPIO20,
			RstPin: machine.GPIO21,
			ScePin: machine.GPIO17,
		},
		Led: machine.LED,
		// Note, we have PWM0 -> GPIO 0|1 or 16|17, PWM1 -> GPIO 2|3 or 18|19 and so on... up to PWM7
		Buzzer: Buzzer{
			Pwm: machine.PWM3,
			Pin: machine.GPIO22,
		},
		Keypad: Keypad{
			Col: [3]machine.Pin{machine.GP3, machine.GP4, machine.GP5},
			Row: [4]machine.Pin{machine.GP6, machine.GP7, machine.GP8, machine.GP9},
		},
	})
}
