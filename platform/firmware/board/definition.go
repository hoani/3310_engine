//go:build tinygo

package board

import (
	"errors"
	"machine"
)

type Pcd struct {
	Spi    *machine.SPI
	SckPin machine.Pin
	SdoPin machine.Pin
	SdiPin machine.Pin
	DcPin  machine.Pin
	RstPin machine.Pin
	ScePin machine.Pin
}

type Buzzer struct {
	Pwm PwmGroup
	Pin machine.Pin
}

type Keypad struct {
	Col [3]machine.Pin
	Row [4]machine.Pin
}

type Definition struct {
	Pcd    Pcd
	Led    machine.Pin
	Buzzer Buzzer
	Keypad Keypad
}

type PwmGroup interface {
	Configure(config machine.PWMConfig) error
	Enable(enable bool)
	SetPeriod(period uint64) error
	Set(channel uint8, value uint32)
	Top() uint32
	Channel(pin machine.Pin) (channel uint8, err error)
}

var Board *Definition

func Set(input *Definition) {
	Board = input
}

func Get() (*Definition, error) {
	if Board == nil {
		return nil, errors.New("Board not supported")
	}
	return Board, nil
}
