package sound

import (
	"time"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/sound/note"
)

type NoteData interface {
	Index() note.Index
	Frequency() float32
	Amplitude() uint8
	Duration() uint8
}

type CustomData struct {
	frequency float32
	amp       uint8
	dur       uint8
}

func (d CustomData) Index() note.Index {
	return note.Custom
}
func (d CustomData) Frequency() float32 {
	return d.frequency
}
func (d CustomData) Amplitude() uint8 {
	return d.amp
}
func (d CustomData) Duration() uint8 {
	return d.dur
}

type BasicData struct {
	index note.Index
	amp   uint8
	dur   uint8
}

func (d BasicData) Index() note.Index {
	return d.index
}
func (d BasicData) Frequency() float32 {
	return note.Frequency(d.index)
}
func (d BasicData) Amplitude() uint8 {
	return d.amp
}
func (d BasicData) Duration() uint8 {
	return d.dur
}

func Custom(frequency float32, amp uint8, dur uint8) CustomData {
	return CustomData{frequency: frequency, amp: amp, dur: dur}
}

func Note(index note.Index, amp uint8, dur uint8) BasicData {
	return BasicData{index: index, amp: amp, dur: dur}
}

func None(dur uint8) BasicData {
	return BasicData{index: note.None, amp: 0, dur: dur}
}

type Basic struct {
	msPerStep   uint16
	notes       []NoteData
	index       int
	frameOffset uint16
}

func Sound(msPerStep uint16, notes ...NoteData) engine.Sound {
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
		return engine.NoneNote(0)
	}

	n := b.notes[b.index]
	b.index++

	dur := time.Millisecond * time.Duration(n.Duration()) * time.Duration(b.msPerStep)
	if n.Index() == note.Custom {
		return engine.CustomNote(n.Frequency(), n.Amplitude(), dur)
	}

	return engine.BasicNote(n.Index(), n.Amplitude(), dur)
}

func (b *Basic) Done() bool {
	return b.index >= len(b.notes)
}
