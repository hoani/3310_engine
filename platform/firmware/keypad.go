package firmware

import (
	"machine"
	"time"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
)

type Keypad struct {
	col     [3]machine.Pin
	row     [4]machine.Pin
	state   [12]bool
	changed [12]bool
}

func NewKeypad(col [3]machine.Pin, row [4]machine.Pin) (*Keypad, *command.Command[engine.Key]) {
	for _, in := range row {
		in.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	for _, out := range col {
		out.Configure(machine.PinConfig{Mode: machine.PinOutput})
		out.Low()
	}

	kp := &Keypad{
		col:   col,
		row:   row,
		state: [12]bool{},
	}

	cmd := command.New[engine.Key]()
	cmd.Register(engine.K1, func() bool { return kp.Get(engine.K1) })
	cmd.Register(engine.K2, func() bool { return kp.Get(engine.K2) })
	cmd.Register(engine.K3, func() bool { return kp.Get(engine.K3) })
	cmd.Register(engine.K4, func() bool { return kp.Get(engine.K4) })
	cmd.Register(engine.K5, func() bool { return kp.Get(engine.K5) })
	cmd.Register(engine.K6, func() bool { return kp.Get(engine.K6) })
	cmd.Register(engine.K7, func() bool { return kp.Get(engine.K7) })
	cmd.Register(engine.K8, func() bool { return kp.Get(engine.K8) })
	cmd.Register(engine.K9, func() bool { return kp.Get(engine.K9) })
	cmd.Register(engine.KStar, func() bool { return kp.Get(engine.KStar) })
	cmd.Register(engine.K0, func() bool { return kp.Get(engine.K0) })
	cmd.Register(engine.KHash, func() bool { return kp.Get(engine.KHash) })

	return kp, cmd
}

func (k *Keypad) Update() {
	index := 0
	for i, out := range k.col {
		out.High()
		// Measured on scope - rise time is actually about 20ns. I don't think these delays really matter...
		time.Sleep(time.Nanosecond * 200)

		for j, in := range k.row {
			index = j*len(k.col) + i
			val := in.Get()
			// Debounce routine, this is probably enough given we only check at 60Hz.
			if val != k.state[index] && !k.changed[index] {
				k.changed[index] = true
				k.state[index] = val
			} else {
				k.changed[index] = false
			}
			index++
		}
		out.Low()
		time.Sleep(time.Nanosecond * 200)
	}
}

func (k *Keypad) Get(key engine.Key) bool {
	return k.state[int(key)]
}
