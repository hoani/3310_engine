package sprite

type sprite struct {
	W     int
	H     int
	count int
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
