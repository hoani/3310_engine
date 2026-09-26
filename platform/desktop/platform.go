//go:build !tinygo

package desktop

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/engine/command"
	"github.com/hoani/3310_engine/platform/desktop/soundplayer"
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
	game        engine.Game
	shader      Shaders
	colors      GameColors
	canvas      *canvas
	shadowing   *ebiten.Image
	scale       float64
	ratio       float64
	lastDraw    time.Time
	snd         *soundplayer.Player
	keypad      *command.CommandImpl[engine.Key]
	sw, sh      int
	drawPort    *image.Rectangle
	displayOpts *DisplayOpts
	display     Display
	config      engine.Config
	debug       Debug
}

type Debug struct {
	p *Platform
}

type Display struct {
	p       *Platform
	enabled bool
}

func (p *Platform) Config(c engine.Config) {
	p.config = c
	ebiten.SetTPS(p.config.Fps)
}

func (p *Platform) Cmd() command.Command[engine.Key] {
	return p.keypad
}

func (p *Platform) Snd() engine.SoundPlayer {
	return p.snd
}

func (p *Platform) Debug() engine.Debug {
	return &p.debug
}

func (d *Debug) Enabled() bool {
	return d.p.config.Debug
}

func (d *Debug) Console(format string, args ...any) {
	if len(args) == 0 {
		fmt.Printf(format)
	} else {
		fmt.Printf(format, args)
	}

	if !strings.HasSuffix(format, "\n") {
		fmt.Printf("\n")
	}
}

func (p *Platform) Display() engine.Display {
	return &p.display
}

func (d *Display) Enable(e bool) {
	d.enabled = e
}

func (p *Platform) Update() error {
	p.keypad.Update()
	if err := p.game.Update(); err != nil {
		return err
	}

	return p.game.Draw(p.canvas) // This gets done here because we don't want to miss frames.
}

func (p *Platform) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xce, 0xf9, 0xe0, 0xff})

	if p.display.enabled == false {
		return // All done
	}

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

	screen.DrawRectShader(p.canvas.Width(), p.canvas.Height(), p.shader.screen, opts)
}

func (p *Platform) computeDrawPort(screenWidth, screenHeight int) {
	x0 := math.Round(float64(screenWidth)-p.scale*float64(p.canvas.Width())) / 2.0
	y0 := math.Round(float64(screenHeight)-p.scale*p.ratio*float64(p.canvas.Height())) / 2.0
	w, h := p.scale*float64(p.canvas.Width()), p.scale*p.ratio*float64(p.canvas.Height())
	rect := image.Rect(int(x0), int(y0), int(x0+w), int(y0+h))
	p.drawPort = &rect
}

func (p *Platform) DrawPort() *image.Rectangle {
	return p.drawPort
}

func (p *Platform) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if p.sw == outsideWidth && p.sh == outsideHeight {
		return p.sw, p.sh // Nothing to change.
	}

	var ratio = p.ratio * float64(outsideWidth) / float64(outsideHeight)
	next := 0.0
	if ratio <= 84.0/48.0 {
		next = math.Floor(float64(outsideWidth) / 84.0)
	} else {
		next = math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
	}
	if next > 6 {
		next -= next * p.displayOpts.Margin
	}

	if next != p.scale {

		if p.config.Debug {
			fmt.Printf("%d set scale %f\n", time.Now().Second(), p.scale)
		}
		p.scale = next // math.Floor(float64(outsideHeight) / (p.ratio * 48.0))
		p.computeDrawPort(outsideWidth, outsideHeight)
	}

	p.sw = outsideWidth
	p.sh = outsideHeight

	return outsideWidth, outsideHeight
}

func HandleError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func Launch(game engine.Game, runner Runner) {

	cmd := runner.Keypad()
	snd, err := soundplayer.New()
	HandleError(err)

	p := &Platform{
		game:      game,
		colors:    NewGameColors(),
		canvas:    NewCanvas(),
		shadowing: ebiten.NewImage(84, 48),
		scale:     10.0, ratio: 1.25,
		lastDraw:    time.Now(),
		snd:         snd,
		keypad:      cmd,
		displayOpts: runner.DisplayOpts(),
		display:     Display{enabled: true},
	}

	p.debug = Debug{p: p}

	game.Setup(p)

	p.shader.shadowing, err = ebiten.NewShader(Shadowing_kage)
	HandleError(err)

	p.shader.screen, err = ebiten.NewShader(Screen_kage)
	HandleError(err)

	runner.Setup(p)
	HandleError(runner.Run())
}

func Run(game engine.Game) {
	r := NewDefaultRunner(game)

	Launch(game, r)
}
