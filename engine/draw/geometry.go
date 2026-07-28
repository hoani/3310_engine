package draw

import "github.com/hoani/3310_engine/engine"

type TriangleBuilder struct {
	c          engine.Canvas
	p0, p1, p2 Point
	shade      uint8
}

func NewTriangleBuilder(c engine.Canvas) *TriangleBuilder {
	return &TriangleBuilder{
		c:     c,
		shade: engine.ShadeBlack,
	}
}

func (b *TriangleBuilder) New(p0, p1, p2 Point) *TriangleBuilder {
	b.p0 = p0
	b.p1 = p1
	b.p2 = p2
	return b
}

func (b *TriangleBuilder) Draw(on bool) {
	shade := engine.ShadeBlack
	if !on {
		shade = engine.ShadeWhite
	}
	b.DrawShade(shade)
}

func (b *TriangleBuilder) DrawShade(shade uint8) {
	filledTriangle(b.c, b.p0, b.p1, b.p2, shade)
}

func hLine(c engine.Canvas, x0, x1, y int, shade uint8) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	for x := x0; x <= x1; x++ {
		_, v := Dither(x, y, shade)
		c.Set(x, y, v)
	}
}

func filledTriangle(c engine.Canvas, p0, p1, p2 Point, shade uint8) {
	// sort points by y ascending: (x0,y0) top ... (x2,y2) bottom
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	}
	if p0.Y > p2.Y {
		p0, p2 = p2, p0
	}
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}

	height := p2.Y - p0.Y
	if height == 0 {
		return // zero height triangle... just ignore it.
	}

	firstHeight := p1.Y - p0.Y
	secondHeight := p2.Y - p1.Y

	for y := p0.Y; y <= p2.Y; y++ {
		secondHalf := y > p1.Y || p1.Y == p0.Y

		// x0 runs the long edge (v0->v2); x1 runs the current short edge
		x0 := p0.X + (p2.X-p0.X)*(y-p0.Y)/height
		var x1 int
		if secondHalf {
			x1 = p1.X
			if secondHeight > 0 {
				x1 += (p2.X - p1.X) * (y - p1.Y) / secondHeight
			}
		} else {
			x1 = p0.X
			if firstHeight > 0 {
				x1 += (p1.X - p0.X) * (y - p0.Y) / firstHeight
			}
		}

		hLine(c, x0, x1, y, shade)
	}
}
