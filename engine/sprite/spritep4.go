package sprite

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hoani/3310_engine/engine"
)

type spriteP4 struct {
	sprite
	Content  string
	Mask     string
	RowWidth int
}

func (s *spriteP4) At(i, j, index int) (shade uint8) {

	if i < 0 || i >= s.W || j < 0 || j >= s.H {
		return engine.SpriteTransparent
	}

	rowBit := i + (s.W * index)
	shift := 7 - (rowBit % 8)

	index = index % s.count
	offset := s.RowWidth*j + rowBit/8

	if s.Mask != "" {
		maskBit := (uint8(s.Mask[offset]) >> shift) & 0x1
		if maskBit == 0 {
			return engine.SpriteTransparent
		}
	}
	inkBit := (uint8(s.Content[offset]) >> shift) & 0x1
	if inkBit != 0 {
		return engine.SpriteBlack
	}

	return engine.SpriteWhite
}

func parseP4(raw string) (int, int, string, error) {
	raw = strings.ReplaceAll(raw, "\r", "") // Deals with windows nonsense
	parts := strings.Split(raw, "\n")
	if len(parts) != 3 {
		return 0, 0, "", fmt.Errorf("invalid file format, got %d parts", len(parts))
	}
	if parts[0] != "P4" {
		return 0, 0, "", errors.New("unknown file format")
	}
	sizes := strings.Split(parts[1], " ")
	if len(sizes) != 2 {
		return 0, 0, "", errors.New("P4 sizes are invalid")
	}
	w, err := strconv.Atoi(sizes[0])
	if err != nil {
		return 0, 0, "", err
	}
	h, err := strconv.Atoi(sizes[1])
	if err != nil {
		return 0, 0, "", err
	}

	return w, h, parts[2], nil
}

func FromP4(raw string, mask string) (*spriteP4, error) {
	w, h, content, err := parseP4(raw)
	if err != nil {
		return nil, err
	}

	var mcontent string
	if mask != "" {
		var mw, mh int
		mw, mh, mcontent, err = parseP4(mask)
		if err != nil {
			return nil, err
		}
		if mw != w || mh != h {
			return nil, errors.New("mismatched raw and mask sizes")
		}
		mask = mcontent
	}

	return &spriteP4{
		sprite: sprite{
			W:     w,
			H:     h,
			count: 1,
		},
		Content:  content,
		Mask:     mcontent,
		RowWidth: (w + 7) / 8,
	}, nil
}

func StripFromP4(raw, mask string, w int) (*spriteP4, error) {
	spr, err := FromP4(raw, mask)
	if err != nil {
		return nil, err
	}
	spr.count = spr.W / w
	spr.W = w
	return spr, nil
}
