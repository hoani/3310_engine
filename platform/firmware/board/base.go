//go:build tinygo

package board

type Base struct {
	pcd          Pcd
	buzzer       Buzzer
	keypad       Keypad
	OnInitialize func()
}

func (b *Base) Initialize() {
	if b.OnInitialize != nil {
		b.OnInitialize()
	}
}

func (b *Base) Pcd() *Pcd {
	return &b.pcd
}

func (b *Base) Buzzer() *Buzzer {
	return &b.buzzer
}

func (b *Base) Keypad() *Keypad {
	return &b.keypad
}
