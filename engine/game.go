package engine

import "github.com/hoani/3310_engine/engine/command"

type Debug interface {
	Console(format string, args ...any)
}

type Platform interface {
	Cmd() command.Command[Key]
	Snd() SoundPlayer
	Debug() Debug
}

type GameInfo struct {
	Fps   int
	Debug bool
}

type Game interface {
	Setup(Platform)
	Info() *GameInfo
	Update() error
	Draw(canvas Canvas) error
}
