package desktop

import (
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hoani/3310_engine/engine"
)

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

type Platform struct {
	game      engine.Game
	shader    Shaders
	colors    GameColors
	canvas    *canvas
	shadowing *ebiten.Image
	scale     float64
	ratio     float64
	resized   bool
}

func (p *Platform) Update() error {
	if err := p.game.Update(); err != nil {
		return err
	}
	return p.game.Draw(p.canvas) // This gets done here because we don't want to miss frames.
}

func (p *Platform) Draw(screen *ebiten.Image) {

	c := ebiten.NewImageFromImage(p.canvas.image)

	prev := ebiten.NewImageFromImage(p.shadowing)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = c
	opts.Images[1] = prev
	opts.Uniforms = map[string]any{
		// falltime's 125ms, one fall time is techncially 5*tau so 3.0/0.125 1/s * 1.0/60 s/frames
		"Tau": 4.0 / (64.0 * 0.125), // rounded frames up to 64.0 for better math
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
	screen.DrawImage(c, &ebiten.DrawImageOptions{})

	shadowOpt := &ebiten.DrawImageOptions{}
	shadowOpt.GeoM.Translate(0.0, float64(p.canvas.Height()))
	screen.DrawImage(p.shadowing, shadowOpt)

}

func (p *Platform) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var ratio = p.ratio * float64(outsideWidth) / float64(outsideHeight)
	if ratio <= 84.0/48.0 {
		p.scale = math.Floor(float64(outsideWidth) / 84.0)
	} else {
		p.scale = math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
	}
	p.scale = math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
	p.resized = true

	fmt.Printf("set scale %f\n", p.scale)

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

	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("Hoani's World")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	p := &Platform{game: game, colors: NewGameColors(), canvas: NewCanvas(), shadowing: ebiten.NewImage(84, 48), scale: 10.0, ratio: 1.25}

	var err error
	p.shader.shadowing, err = ebiten.NewShader(Shadowing_kage)
	handleError(err)

	p.shader.screen, err = ebiten.NewShader(Screen_kage)
	handleError(err)

	handleError(ebiten.RunGame(p))
}
