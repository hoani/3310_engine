package board

import "machine"

type LedPwm struct {
	Pwm
	Channel uint8
}

type RgbBacklight struct {
	R LedPwm
	G LedPwm
	B LedPwm
}

func (l *RgbBacklight) Setup() error {
	if err := l.R.Setup(); err != nil {
		return err
	}
	if err := l.B.Setup(); err != nil {
		return err
	}
	if err := l.G.Setup(); err != nil {
		return err
	}
	return nil
}

func (l *RgbBacklight) Set(r, g, b uint8) {
	l.R.Set(r)
	l.G.Set(g)
	l.B.Set(b)
}

func (l *RgbBacklight) Enable(e bool) {
	l.R.Enable(e)
	l.G.Enable(e)
	l.B.Enable(e)
}

func (l *LedPwm) Setup() error {
	l.Pwm.Pwm.Configure(machine.PWMConfig{
		Period: 1e6,
	})
	ch, err := l.Pwm.Pwm.Channel(l.Pin)
	if err != nil {
		return err
	}
	l.Channel = ch
	l.Pwm.Pwm.Enable(false)
	return nil
}

func (l *LedPwm) Set(duty uint8) {
	top := l.Pwm.Pwm.Top()
	d := (uint32(duty) * top) / 255
	l.Pwm.Pwm.Set(l.Channel, d)
}

func (l *LedPwm) Enable(e bool) {
	l.Pwm.Pwm.Enable(e)
}
