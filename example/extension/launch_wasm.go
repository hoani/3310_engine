//go:build !tinygo && js && wasm

package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform"
	"github.com/hoani/3310_engine/platform/desktop"
)

type WasmPrelaunch struct {
	done bool
	name string
	fnt  *text.GoTextFaceSource
	w, h int
}

func NewWasmPrelaunch(name string) desktop.PreLaunch {
	fnt, err := text.NewGoTextFaceSource(bytes.NewReader(desktop.Debug_ttf))
	desktop.HandleError(err)

	return &WasmPrelaunch{
		name: name,
		done: false,
		fnt:  fnt,
	}
}

func (l *WasmPrelaunch) Setup() error {
	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle(l.name)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
	return nil
}

func (l *WasmPrelaunch) Done() bool {
	return l.done
}

func (l *WasmPrelaunch) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		l.done = true
	}
	return nil
}

func (l *WasmPrelaunch) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(l.w/2), float64(l.h/2))
	op.ColorScale.ScaleWithColor(color.White)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	op.LineSpacing = 12 * 0.8

	text.Draw(screen, "Left Click to Continue", &text.GoTextFace{
		Source: l.fnt,
		Size:   12,
	}, op)
}

func (l *WasmPrelaunch) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	l.w = outsideWidth
	l.h = outsideHeight
	return outsideWidth, outsideHeight
}

func Launch(g engine.Game) {
	platform.Launch(g, desktop.NewRunner(NewWasmPrelaunch("Interact")))
}
