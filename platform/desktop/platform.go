//go:build !tinygo

package desktop

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/platform/desktop/soundplayer"
	"github.com/shirou/gopsutil/v4/process"
)

//go:embed Diary_of_an_8-bit_mage.otf
var Debug_ttf []byte

//go:embed shadowing.kage
var Shadowing_kage []byte

//go:embed shader.kage
var Screen_kage []byte

// Pass in a color like 0xffaabb
func NormalizeColor(rgb uint32) []float32 {
	return []float32{
		float32(0xff&(rgb>>16)) / 255.0,
		float32(0xff&(rgb>>8)) / 255.0,
		float32(0xff&(rgb>>0)) / 255.0,
		1.0,
	}
}

type GameColors struct {
	On   []float32
	Off  []float32
	Back []float32
}

func NewGameColors() GameColors {
	return GameColors{
		On:   NormalizeColor(0x43523d),
		Off:  NormalizeColor(0xc7f0d8),
		Back: NormalizeColor(0xc9f2da),
	}
}

type Shaders struct {
	shadowing *ebiten.Shader
	screen    *ebiten.Shader
}

type Debug struct {
	fnt  *text.GoTextFaceSource
	proc *process.Process
	cpu  float64
}

type Platform struct {
	game      engine.Game
	shader    Shaders
	colors    GameColors
	canvas    *canvas
	shadowing *ebiten.Image
	scale     float64
	ratio     float64
	resized   bool
	debug     *Debug
	lastDraw  time.Time
	snd       *soundplayer.Player
}

func (p *Platform) Console(format string, args ...any) {
	fmt.Printf(format, args...)

	if !strings.HasSuffix(format, "\n") {
		fmt.Printf("\n")
	}
}

func (p *Platform) Update() error {
	if err := p.game.Update(); err != nil {
		return err
	}
	if p.debug != nil {
		p.measureCpu()
	}

	return p.game.Draw(p.canvas) // This gets done here because we don't want to miss frames.
}

func (p *Platform) Draw(screen *ebiten.Image) {

	// Very pendantic, but ensures consistent shadowing with varying screen FPS.
	dt := time.Since(p.lastDraw).Seconds()
	p.lastDraw = time.Now()
	fps := float64(ebiten.ActualFPS())
	dt = math.Min(2.0/fps, math.Max(dt, 0.5/fps)) // clamp to ride through stalls.

	c := ebiten.NewImageFromImage(p.canvas.Image())

	prev := ebiten.NewImageFromImage(p.shadowing)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = c
	opts.Images[1] = prev
	opts.Uniforms = map[string]any{
		// Rise time tau is around 60ms, 1/16 is close enough, so 1.0/0.0625 * 1/s * 1.0/60 s/frames
		"TauRise": 1.0 - math.Exp(-dt/0.060), // rounded frames up to 64.0 for better math
		// falltime's tau is 125ms, one fall time is 1.0/0.125 * 1/s * 1.0/60 s/frames
		"TauFall": 1.0 - math.Exp(-dt/0.125), // rounded frames up to 64.0 for better math
	}

	p.shadowing.DrawRectShader(p.canvas.Width(), p.canvas.Height(), p.shader.shadowing, opts)

	xOffset := math.Round(float64(screen.Bounds().Size().X)-p.scale*float64(p.canvas.Width())) / 2.0
	yOffset := math.Round(float64(screen.Bounds().Size().Y)-p.scale*p.ratio*float64(p.canvas.Height())) / 2.0

	opts = &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"PixelOn":   p.colors.On,
		"PixelOff":  p.colors.Off,
		"PixelBack": p.colors.Back,
		"XScale":    p.scale,
		"YScale":    p.ratio * p.scale,
		"XOffset":   xOffset,
		"YOffset":   yOffset,
	}
	opts.Images[0] = p.shadowing
	opts.GeoM.Scale(p.scale, p.scale*p.ratio)

	opts.GeoM.Translate(xOffset, yOffset)

	if p.resized {
		screen.Fill(color.RGBA{0xce, 0xf9, 0xe0, 0xff})
	}

	screen.DrawRectShader(p.canvas.Width(), p.canvas.Height(), p.shader.screen, opts)

	if p.debug != nil {
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
		if !math.IsNaN(p.debug.cpu) {
			statStr += fmt.Sprintf(
				"\n\nCPU: %.1f%%",
				p.debug.cpu,
			)
		}

		text.Draw(screen, statStr, &text.GoTextFace{
			Source: p.debug.fnt,
			Size:   12,
		}, op)

		miniX := float64(screen.Bounds().Dx()) - float64(p.canvas.Width())

		miniOpt := &ebiten.DrawImageOptions{}
		miniOpt.GeoM.Translate(miniX, 0.0)
		screen.DrawImage(c, miniOpt)

		miniOpt.GeoM.Translate(0.0, float64(p.canvas.Height()))
		screen.DrawImage(p.shadowing, miniOpt)
	}
}

func (p *Platform) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var ratio = p.ratio * float64(outsideWidth) / float64(outsideHeight)
	next := 0.0
	if ratio <= 84.0/48.0 {
		next = math.Floor(float64(outsideWidth) / 84.0)
	} else {
		next = math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
	}
	if next != p.scale {
		fmt.Printf("set scale %f\n", p.scale)
		p.scale = math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
		p.resized = true
	}

	return outsideWidth, outsideHeight
}

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func Run(game engine.Game) {

	cmd := NewKeypad()

	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("Hoani's World")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	snd, err := soundplayer.New()
	handleError(err)

	p := &Platform{game: game, colors: NewGameColors(), canvas: NewCanvas(), shadowing: ebiten.NewImage(84, 48), scale: 10.0, ratio: 1.25, lastDraw: time.Now(), snd: snd}

	game.Setup(cmd, snd, p)

	if game.Info().Debug {
		proc, err := NewProcess()
		handleError(err)
		fnt, err := text.NewGoTextFaceSource(bytes.NewReader(Debug_ttf))
		handleError(err)
		p.debug = &Debug{
			proc: proc,
			fnt:  fnt,
		}
	}

	p.shader.shadowing, err = ebiten.NewShader(Shadowing_kage)
	handleError(err)

	p.shader.screen, err = ebiten.NewShader(Screen_kage)
	handleError(err)

	ebiten.SetRunnableOnUnfocused(true)

	handleError(ebiten.RunGame(p))

}
