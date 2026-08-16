package engine

import "github.com/hoani/3310_engine/engine/command"

type Debug interface {
	Console(format string, args ...any)
}

type GameInfo struct {
	Fps   int
	Debug bool
}

type Game interface {
	Setup(cmd *command.Command[Key], snd SoundPlayer, debug Debug)
	Info() *GameInfo
	Update() error
	Draw(canvas Canvas) error
}

type Key int

const (
	K1 Key = iota
	K2
	K3
	K4
	K5
	K6
	K7
	K8
	K9
	KStar
	K0
	KHash
)
