package sound

import (
	"time"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/sound/note"
)

type BasicNote struct {
	index note.Index
	amp   uint8
	dur   uint8
}

func Note(index note.Index, amp uint8, dur uint8) BasicNote {
	return BasicNote{index: index, amp: amp, dur: dur}
}

func None(dur uint8) BasicNote {
	return BasicNote{index: note.None, amp: 0, dur: dur}
}

type Basic struct {
	msPerStep   uint16
	notes       []BasicNote
	index       int
	frameOffset uint16
}

func Sound(msPerStep uint16, notes ...BasicNote) engine.Sound {
	return &Basic{
		msPerStep:   msPerStep,
		notes:       notes,
		index:       0,
		frameOffset: 0,
	}
}

func (b *Basic) Reset() {
	b.index = 0
	b.frameOffset = 0
}

func (b *Basic) Next() engine.Note {

	if b.Done() {
		return engine.Note{Index: note.None}
	}

	n := b.notes[b.index]
	b.index++

	dur := time.Millisecond * time.Duration(n.dur) * time.Duration(b.msPerStep)
	return engine.Note{Index: n.index, Amplitude: n.amp, Duration: dur}
}

func (b *Basic) Done() bool {
	return b.index >= len(b.notes)
}
