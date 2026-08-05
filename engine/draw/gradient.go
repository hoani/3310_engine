package draw

type linearGradient struct {
	p0     Point
	p1     Point
	shade0 uint8
	shade1 uint8
}

func NewLinearGradient(p0, p1 Point, shade0, shade1 uint8) Gradient {
	return &linearGradient{
		p0:     p0,
		p1:     p1,
		shade0: shade0,
		shade1: shade1,
	}
}

func (g *linearGradient) Calculate(x, y int) uint8 {
	dx := g.p1.X - g.p0.X
	dy := g.p1.Y - g.p0.Y
	d2 := dx*dx + dy*dy
	if d2 == 0 {
		return g.shade0 // No better guess
	}

	// Project (x, y) onto the p0→p1 axis, normalized to [0, 1].
	t := ((x-g.p0.X)*dx + (y-g.p0.Y)*dy)
	if t < 0 {

		t = 0
	} else if t > d2 {
		t = d2
	}
	s0 := int(g.shade0)
	s1 := int(g.shade1)
	return uint8(s0 + t*(s1-s0)/d2)
}

type radialGradient struct {
	center  Point
	radius0 int
	radius1 int
	shade0  uint8
	shade1  uint8
}

func NewRadialGradient(center Point, radius0, radius1 int, shade0, shade1 uint8) Gradient {
	return &radialGradient{
		center:  center,
		radius0: radius0,
		radius1: radius1,
		shade0:  shade0,
		shade1:  shade1,
	}
}

func (g *radialGradient) Calculate(x, y int) uint8 {
	dx := (g.center.X - x) * 4
	dy := (g.center.Y - y) * 5
	d2 := (dx*dx + dy*dy) / 16
	if d2 <= g.radius0*g.radius0 {
		return g.shade0
	}
	if d2 >= g.radius1*g.radius1 {
		return g.shade1
	}
	// Binary search for the gradient - more efficient than a square root
	lo, hi := g.radius0, g.radius1
	for hi-lo > 1 {
		mid := (lo + hi) / 2
		if mid*mid <= d2 {
			lo = mid
		} else {
			hi = mid
		}
	}
	d := lo

	span := g.radius1 - g.radius0
	s0 := int(g.shade0)
	s1 := int(g.shade1)
	return uint8(s0 + (d-g.radius0)*(s1-s0)/span)
}
