//go:build !tinygo

package desktop

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
)

func KeyPressedFunc(keys ...ebiten.Key) func() bool {
	return func() bool {
		result := false
		for _, key := range keys {
			result = result || (inpututil.KeyPressDuration(key) != 0)
		}
		return result
	}
}

func NewKeypad() *command.CommandImpl[engine.Key] {

	cmd := command.New[engine.Key]()
	cmd.Register(engine.K1, KeyPressedFunc(ebiten.Key1, ebiten.KeyNumpad7))
	cmd.Register(engine.K2, KeyPressedFunc(ebiten.Key2, ebiten.KeyNumpad8))
	cmd.Register(engine.K3, KeyPressedFunc(ebiten.Key3, ebiten.KeyNumpad9))
	cmd.Register(engine.K4, KeyPressedFunc(ebiten.KeyQ, ebiten.KeyNumpad4))
	cmd.Register(engine.K5, KeyPressedFunc(ebiten.KeyW, ebiten.KeyNumpad5))
	cmd.Register(engine.K6, KeyPressedFunc(ebiten.KeyE, ebiten.KeyNumpad6))
	cmd.Register(engine.K7, KeyPressedFunc(ebiten.KeyA, ebiten.KeyNumpad1))
	cmd.Register(engine.K8, KeyPressedFunc(ebiten.KeyS, ebiten.KeyNumpad2))
	cmd.Register(engine.K9, KeyPressedFunc(ebiten.KeyD, ebiten.KeyNumpad3))
	cmd.Register(engine.KStar, KeyPressedFunc(ebiten.KeyZ, ebiten.KeyNumpad0))
	cmd.Register(engine.K0, KeyPressedFunc(ebiten.KeyX, ebiten.KeyNumpadDecimal))
	cmd.Register(engine.KHash, KeyPressedFunc(ebiten.KeyC, ebiten.KeyNumpadAdd))
	cmd.Register(engine.KA, KeyPressedFunc(ebiten.KeyU))
	cmd.Register(engine.KB, KeyPressedFunc(ebiten.KeyI))
	cmd.Register(engine.KC, KeyPressedFunc(ebiten.KeyO))
	cmd.Register(engine.KD, KeyPressedFunc(ebiten.KeyP))

	return cmd
}
