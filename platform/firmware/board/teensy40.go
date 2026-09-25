//go:build tinygo && teensy40

package board

import (
	"machine"
)

func init() {
	Set(&Base{
		Pcd: Pcd{
			Spi:    machine.SPI0,
			SckPin: machine.D13,
			SdoPin: machine.D11,
			SdiPin: machine.D12, // Unused but maps to teensy 4.0 pinout
			DcPin:  machine.D14,
			RstPin: machine.D15,
			ScePin: machine.D10,
		},
		Buzzer: Buzzer{
			Pwm: &PwmAdaptor{}, // Not actualy supported... oh dear.
			Pin: machine.D23,
		},
		Keypad: Keypad{
			Col: [4]machine.Pin{machine.D2, machine.D3, machine.D4, machine.D5},
			Row: [4]machine.Pin{machine.D6, machine.D7, machine.D8, machine.D9},
		},
	})
}

type PwmAdaptor struct {
}

func (p *PwmAdaptor) Configure(config machine.PWMConfig) error           { return nil }
func (p *PwmAdaptor) Enable(enable bool)                                 {}
func (p *PwmAdaptor) SetPeriod(period uint64) error                      { return nil }
func (p *PwmAdaptor) Set(channel uint8, value uint32)                    {}
func (p *PwmAdaptor) Top() uint32                                        { return 0xFFFF }
func (p *PwmAdaptor) Channel(pin machine.Pin) (channel uint8, err error) { return 0, nil }
