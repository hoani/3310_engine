//go:build tinygo && pico

package board

import "machine"

func DefaultDefinition() *Base {
	return &Base{
		pcd: Pcd{
			Spi:    machine.SPI0,
			SckPin: machine.GPIO18,
			SdoPin: machine.GPIO19,
			SdiPin: machine.GPIO16,
			DcPin:  machine.GPIO20,
			RstPin: machine.GPIO21,
			ScePin: machine.GPIO17,
		},
		// Note, we have PWM0 -> GPIO 0|1 or 16|17, PWM1 -> GPIO 2|3 or 18|19 and so on... up to PWM7
		buzzer: Buzzer{
			Pwm: machine.PWM7,
			Pin: machine.GPIO15,
		},
		keypad: Keypad{
			Col: [4]machine.Pin{machine.GP3, machine.GP4, machine.GP5, machine.GP2},
			Row: [4]machine.Pin{machine.GP6, machine.GP7, machine.GP8, machine.GP9},
		},
	}
}

func init() {
	Set(V1())
}
