package engine

import "github.com/hoani/3310_engine/engine/command"

type Debug interface {
	Console(format string, args ...any)
}

type GameInfo struct {
	Fps         int
	Debug       bool
	Illuminated bool
}

type Game interface {
	Setup(cmd command.Command[Key], snd SoundPlayer, debug Debug)
	Info() *GameInfo
	Update() error
	Draw(canvas Canvas) error
}
