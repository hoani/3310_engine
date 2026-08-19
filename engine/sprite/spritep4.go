package sprite

import (
	"errors"

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

	index = index % s.count

	rowBit := i + (s.W * index)
	shift := 7 - (rowBit % 8)

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
	magic, offset, err := nextLine(raw, 0)
	if err != nil {
		return 0, 0, "", errors.New("missing magic header")
	}

	if magic != "P4" {
		return 0, 0, "", errors.New("file format is not P4")
	}

	sizeLine, offset, err := nextLine(raw, offset)
	if err != nil {
		return 0, 0, "", errors.New("missing size line")
	}

	w, h, err := parseSizes(sizeLine)
	if err != nil {
		return 0, 0, "", err
	}

	return w, h, raw[offset:], nil // body verbatim
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
