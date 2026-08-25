//go:build !tinygo

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hoani/3310_engine/platform/desktop"
)

type frameExtension struct{}

type countExtension struct {
	count int
	fnt   *text.GoTextFaceSource
}

func (e *frameExtension) Update() error {

	return nil
}

func (e *frameExtension) Draw(screen *ebiten.Image, port *image.Rectangle) {
	x0, y0, x1, y1 := float64(port.Min.X), float64(port.Min.Y), float64(port.Max.X), float64(port.Max.Y)
	ebitenutil.DrawLine(screen, x0, y0, x1, y0, color.RGBA{255, 0, 0, 0})
	ebitenutil.DrawLine(screen, x0, y0, x0, y1, color.RGBA{255, 0, 0, 0})
	ebitenutil.DrawLine(screen, x1, y0, x1, y1, color.RGBA{255, 0, 0, 0})
	ebitenutil.DrawLine(screen, x0, y1, x1, y1, color.RGBA{255, 0, 0, 0})
}

func (e *countExtension) Update() error {
	e.count++
	if e.fnt == nil {
		fnt, err := text.NewGoTextFaceSource(bytes.NewReader(desktop.Debug_ttf))
		if err != nil {
			return err
		}
		e.fnt = fnt
	}
	return nil
}

func (e *countExtension) Draw(screen *ebiten.Image, port *image.Rectangle) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(port.Max.X), float64(port.Max.Y))
	op.ColorScale.ScaleWithColor(color.Black)
	op.PrimaryAlign = text.AlignStart
	op.SecondaryAlign = text.AlignEnd
	op.LineSpacing = 12 * 0.8

	str := fmt.Sprintf("Count: %d", e.count)

	text.Draw(screen, str, &text.GoTextFace{
		Source: e.fnt,
		Size:   12,
	}, op)
}
