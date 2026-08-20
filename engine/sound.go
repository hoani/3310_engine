package engine

import (
	"time"

	"github.com/hoani/3310_engine/engine/sound/note"
)

type Note interface {
	Frequency() float32
	Index() note.Index
	Amplitude() uint8
	Duration() time.Duration
}

type Sound interface {
	Reset()
	Next() Note
	Done() bool
}

// Provided by platform.
type SoundPlayer interface {
	Play(s Sound)
	Track(s Sound, loop bool)
	Stop()
}

type basicNote struct {
	index     note.Index
	amplitude uint8
	duration  time.Duration
}

func BasicNote(index note.Index, amplitude uint8, duration time.Duration) Note {
	return basicNote{index: index, amplitude: amplitude, duration: duration}
}

func NoneNote(duration time.Duration) Note {
	return basicNote{index: note.None, amplitude: 0, duration: duration}
}

func (n basicNote) Frequency() float32 {
	return note.Frequency(n.index)
}

func (n basicNote) Index() note.Index {
	return n.index
}

func (n basicNote) Amplitude() uint8 {
	return n.amplitude
}

func (n basicNote) Duration() time.Duration {
	return n.duration
}

type customNote struct {
	frequency float32
	amplitude uint8
	duration  time.Duration
}

func CustomNote(frequency float32, amplitude uint8, duration time.Duration) Note {
	return customNote{frequency: frequency, amplitude: amplitude, duration: duration}
}

func (n customNote) Index() note.Index {
	return note.Custom
}

func (n customNote) Frequency() float32 {
	return n.frequency
}

func (n customNote) Amplitude() uint8 {
	return n.amplitude
}

func (n customNote) Duration() time.Duration {
	return n.duration
}
