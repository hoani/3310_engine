//go:build tinygo

package board

import (
	"errors"
	"machine"
)

var Board Definition

func Set(input Definition) {
	Board = input
}

func Get() (Definition, error) {
	if Board == nil {
		return nil, errors.New("Board not supported")
	}
	return Board, nil
}

type Definition interface {
	Initialize()
	Pcd() *Pcd
	Buzzer() *Buzzer
	Keypad() *Keypad
}

type Pcd struct {
	Spi    *machine.SPI
	SckPin machine.Pin
	SdoPin machine.Pin
	SdiPin machine.Pin
	DcPin  machine.Pin
	RstPin machine.Pin
	ScePin machine.Pin
}

type Pwm struct {
	Pwm PwmGroup
	Pin machine.Pin
}

type Buzzer Pwm

type Keypad struct {
	Col [4]machine.Pin
	Row [4]machine.Pin
}

type PwmGroup interface {
	Configure(config machine.PWMConfig) error
	Enable(enable bool)
	SetPeriod(period uint64) error
	Set(channel uint8, value uint32)
	Top() uint32
	Channel(pin machine.Pin) (channel uint8, err error)
}
