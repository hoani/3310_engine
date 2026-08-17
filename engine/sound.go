package engine

import (
	"time"

	"github.com/hoani/3310_engine/engine/sound/note"
)

type Note struct {
	Index     note.Index
	Amplitude uint8
	Duration  time.Duration
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
