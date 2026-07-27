package desktop

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

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
