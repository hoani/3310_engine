package desktop

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hoani/3310_engine/engine/canvas"
)

//go:embed runner.png
var Runner_png []byte

//go:embed shadowing.kage
var Shadowing_kage []byte

//go:embed shader.kage
var Screen_kage []byte

type Sheet struct {
	image  *ebiten.Image
	ox, oy int
	W, H   int
	Count  int
}

func NewSheet(image *ebiten.Image, ox, oy, w, h, count int) *Sheet {
	return &Sheet{
		image: image,
		ox:    ox,
		oy:    oy,
		W:     w,
		H:     h,
		Count: count,
	}
}

func (s *Sheet) Frame(index int) image.Image {
	index = index % s.Count
	sx, sy := s.ox+index*s.W, s.oy
	return s.image.SubImage(image.Rect(sx, sy, sx+s.W, sy+s.H))
}

type Sprite struct {
	sheet *Sheet
	count int
}

func NewSprite(sh *Sheet) *Sprite {
	return &Sprite{
		sheet: sh,
		count: 0,
	}
}

func (s *Sprite) Update() error {
	s.count++
	return nil
}

func (s *Sprite) Frame() *ebiten.Image {
	i := (s.count / 60) % s.sheet.Count
	return s.sheet.Frame(i).(*ebiten.Image)
}

func (s *Sprite) Draw(x, y float64, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(s.sheet.W)/2, -float64(s.sheet.H)/2)
	op.GeoM.Translate(x, y)
	screen.DrawImage(s.Frame(), op)
}

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

type Game struct {
	sprites   []*Sprite
	count     int
	shader    Shaders
	colors    GameColors
	canvas    canvas.Canvas
	shadowing *ebiten.Image
	scale     float64
	ratio     float64
	resized   bool
}

func (g *Game) Update() error {
	g.count++
	for _, s := range g.sprites {
		s.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x0 := float64((g.count / 8) % g.canvas.Image().Bounds().Dx())
	y0 := float64(g.canvas.Image().Bounds().Dy())/2 - float64(32*(len(g.sprites)/2))

	g.canvas.Clear()
	c := ebiten.NewImageFromImage(g.canvas.Image())
	for i, s := range g.sprites {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(s.sheet.W)/2, -float64(s.sheet.H)/2)
		op.GeoM.Translate(x0, y0+float64(i)*32.0)
		c.DrawImage(s.Frame(), op)
	}

	prev := ebiten.NewImageFromImage(g.shadowing)
	opts := &ebiten.DrawRectShaderOptions{}
	opts.Images[0] = c
	opts.Images[1] = prev
	opts.Uniforms = map[string]any{
		// falltime's 125ms, one fall time is techncially 5*tau so 3.0/0.125 1/s * 1.0/60 s/frames
		"Tau": 4.0 / (64.0 * 0.125), // rounded frames up to 64.0 for better math
	}

	g.shadowing.DrawRectShader(g.canvas.Size().X, g.canvas.Size().Y, g.shader.shadowing, opts)

	xOffset := math.Round(float64(screen.Bounds().Size().X)-g.scale*float64(g.canvas.Size().X)) / 2.0
	yOffset := math.Round(float64(screen.Bounds().Size().Y)-g.scale*g.ratio*float64(g.canvas.Size().Y)) / 2.0

	opts = &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"PixelOn":  g.colors.On,
		"PixelOff": g.colors.Off,
		"XScale":   g.scale,
		"YScale":   g.ratio * g.scale,
		"XOffset":  xOffset,
		"YOffset":  yOffset,
	}
	opts.Images[0] = g.shadowing
	opts.GeoM.Scale(g.scale, g.scale*g.ratio)

	opts.GeoM.Translate(xOffset, yOffset)

	if g.resized {
		screen.Fill(color.RGBA{0xce, 0xf9, 0xe0, 0xff})
	}

	screen.DrawRectShader(g.canvas.Size().X, g.canvas.Size().Y, g.shader.screen, opts)
	screen.DrawImage(c, &ebiten.DrawImageOptions{})

	shadowOpt := &ebiten.DrawImageOptions{}
	shadowOpt.GeoM.Translate(0.0, float64(g.canvas.Size().Y))
	screen.DrawImage(g.shadowing, shadowOpt)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	var ratio = g.ratio * float64(outsideWidth) / float64(outsideHeight)
	if ratio <= 84.0/48.0 {
		g.scale = math.Floor(float64(outsideWidth) / 84.0)
	} else {
		g.scale = math.Floor(float64(outsideHeight) / (g.ratio * 48.0))
	}
	g.scale = math.Floor(float64(outsideHeight) / (g.ratio * 48.0))
	g.resized = true

	fmt.Printf("set scale %f\n", g.scale)

	return outsideWidth, outsideHeight
}

func handleError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func Run() {

	ebiten.SetWindowSize(840, 480)
	ebiten.SetWindowTitle("Hoani's World")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := &Game{colors: NewGameColors(), canvas: canvas.New().Build(), shadowing: ebiten.NewImage(84, 48), scale: 10.0, ratio: 1.25}
	runnerPng, _, err := image.Decode(bytes.NewReader(Runner_png))
	handleError(err)
	runnerImg := ebiten.NewImageFromImage(runnerPng)

	g.shader.shadowing, err = ebiten.NewShader(Shadowing_kage)
	handleError(err)

	g.shader.screen, err = ebiten.NewShader(Screen_kage)
	handleError(err)

	g.sprites = append(g.sprites,
		NewSprite(NewSheet(runnerImg, 0, 0, 32, 32, 5)),
		// NewSprite(NewSheet(runnerImg, 0, 32, 32, 32, 8)),
		// NewSprite(NewSheet(runnerImg, 0, 64, 32, 32, 4)),
	)

	handleError(ebiten.RunGame(g))
}
