package engine

import "github.com/hoani/3310_engine/engine/command"

type Debug interface {
	Console(format string, args ...any)
	Enabled() bool
}

type Display interface {
	Enable(bool)
}

type Runtime interface {
	SetFps(int)
}

type Platform interface {
	Config(Config)
	Cmd() command.Command[Key]
	Snd() SoundPlayer
	Debug() Debug
	// Display() Display
	// Runtime() Runtime
}

type Config struct {
	Fps   int
	Debug bool
}

type Game interface {
	Setup(Platform)
	Update() error
	Draw(canvas Canvas) error
}
