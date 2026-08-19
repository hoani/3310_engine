package sprite

import (
	"errors"
	"strconv"

	"github.com/hoani/3310_engine/engine"
)

type spriteP5 struct {
	sprite
	Content string
}

func (s *spriteP5) At(i, j, index int) (shade uint8) {
	if i < 0 || i >= s.W || j < 0 || j >= s.H {
		return engine.SpriteTransparent
	}
	index = index % s.count
	offset := s.count*s.W*j + i + (s.W * index)
	return uint8(s.Content[offset])
}

func FromP5(raw string) (*spriteP5, error) {
	magic, offset, err := nextLine(raw, 0)
	if err != nil {
		return nil, errors.New("missing magic header")
	}
	if magic != "P5" {
		return nil, errors.New("file format is not P5")
	}

	sizeLine, offset, err := nextLine(raw, offset)
	if err != nil {
		return nil, errors.New("missing size header")
	}
	w, h, err := parseSizes(sizeLine)
	if err != nil {
		return nil, err
	}

	maxLine, offset, err := nextLine(raw, offset)
	if err != nil {
		return nil, errors.New("missing maxvalue header")
	}
	maxval, err := strconv.Atoi(maxLine)
	if err != nil {
		return nil, err
	}
	if maxval != 255 {
		return nil, errors.New("Only 8-bit P5 is supported")
	}

	return &spriteP5{
		sprite: sprite{
			W:     w,
			H:     h,
			count: 1,
		},
		Content: raw[offset:], // body verbatim
	}, nil
}

func StripFromP5(raw string, w int) (*spriteP5, error) {
	spr, err := FromP5(raw)
	if err != nil {
		return nil, err
	}
	spr.count = spr.W / w
	spr.W = w
	return spr, nil
}
