package engine

import "github.com/hoani/3310_engine/engine/command"

type Canvas interface {
	Clear()
	Width() int
	Height() int
	Set(x, y int, val bool)
	Get(x, y int) bool
}

type Game interface {
	Setup(cmd *command.Command[Key])
	Fps() int
	Update() error
	Draw(canvas Canvas) error
}

type Sprite interface {
	At(i, j, index int) (visible bool, on bool)
	Width() int
	Height() int
	Count() int
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
