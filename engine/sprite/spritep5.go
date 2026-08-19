package sprite

import (
	"errors"
	"strconv"
	"strings"

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
	raw = strings.ReplaceAll(raw, "\r", "") // Deals with windows nonsense
	parts := strings.Split(raw, "\n")
	if len(parts) != 4 {
		return nil, errors.New("invalid file format")
	}
	if parts[0] != "P5" {
		return nil, errors.New("unknown file format")
	}
	sizes := strings.Split(parts[1], " ")
	if len(sizes) != 2 {
		return nil, errors.New("P5 sizes are invalid")
	}
	w, err := strconv.Atoi(sizes[0])
	if err != nil {
		return nil, err
	}
	h, err := strconv.Atoi(sizes[1])
	if err != nil {
		return nil, err
	}

	return &spriteP5{
		sprite: sprite{
			W:     w,
			H:     h,
			count: 1,
		},
		Content: parts[3],
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
