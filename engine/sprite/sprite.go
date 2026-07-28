package sprite

import (
	"errors"
	"strconv"
	"strings"
)

type sprite struct {
	Content string
	W       int
	H       int
	count   int
}

func (s *sprite) At(i, j, index int) (shade uint8) {
	index = index % s.count
	offset := s.count*s.W*j + i + (s.W * index)
	return uint8(s.Content[offset])
}

func (s *sprite) Width() int {
	return s.W
}

func (s *sprite) Height() int {
	return s.H
}

func (s *sprite) Count() int {
	return s.count
}

func FromP5(raw string) (*sprite, error) {
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

	return &sprite{
		Content: parts[3],
		W:       w,
		H:       h,
		count:   1,
	}, nil
}

func StripFromP5(raw string, w int) (*sprite, error) {
	spr, err := FromP5(raw)
	if err != nil {
		return nil, err
	}
	spr.count = spr.W / w
	spr.W = w
	return spr, nil
}
