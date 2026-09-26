package desktop

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform/desktop/overlay"
	"github.com/shirou/gopsutil/v4/process"
)

var Debug_ttf []byte = overlay.Font_ttf

type debugExtension struct {
	g    engine.Game
	p    engine.Platform
	fnt  *text.GoTextFaceSource
	proc *process.Process
	cpu  float64
}

func NewDebugExtension(g engine.Game) Extension {
	return &debugExtension{
		g: g,
	}
}

func (e *debugExtension) Setup(p engine.Platform) error {
	proc, err := NewProcess()
	if err != nil {
		return err
	}
	e.p = p
	e.proc = proc
	e.fnt, err = text.NewGoTextFaceSource(bytes.NewReader(Debug_ttf))
	return err
}
func (e *debugExtension) Update() error {
	if !e.p.Debug().Enabled() {
		return nil
	}
	e.measureCpu()
	return nil
}

func (e *debugExtension) Draw(screen *ebiten.Image, port *image.Rectangle) {
	if !e.p.Debug().Enabled() {
		return
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(4.0), float64(4.0))
	op.ColorScale.ScaleWithColor(color.Black)
	op.PrimaryAlign = text.AlignStart
	op.SecondaryAlign = text.AlignStart
	op.LineSpacing = 12 * 0.8

	statStr := fmt.Sprintf(
		"FPS: %.1f\n\nTPS: %.1f",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
	)
	if !math.IsNaN(e.cpu) {
		statStr += fmt.Sprintf("\n\nCPU: %.1f%%", e.cpu)
	}

	text.Draw(screen, statStr, &text.GoTextFace{Source: e.fnt, Size: 12}, op)
}
