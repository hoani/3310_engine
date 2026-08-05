package draw

import (
	"github.com/hoani/3310_engine/engine"
)

type RectangleBuilder struct {
	p0 Point
	p1 Point
}

type CircleBuilder struct {
	center   Point
	diameter int16
}

type drawPixel func(x, y int)

type Gradient interface {
	Calculate(x, y int) uint8
}

type ShapeDrawer interface {
	Draw(drawPixel drawPixel)
}

type ShapeBuilder struct {
	c         engine.Canvas
	rectangle RectangleBuilder
	circle    CircleBuilder
	triangle  TriangleBuilder
	active    ShapeDrawer
	shade     uint8
	gradient  Gradient
}

func NewShapeBuilder(c engine.Canvas) *ShapeBuilder {
	return &ShapeBuilder{
		c:         c,
		rectangle: RectangleBuilder{},
		circle:    CircleBuilder{},
		triangle:  TriangleBuilder{},
		active:    nil,
		shade:     0x00,
		gradient:  nil,
	}
}

func (sb *ShapeBuilder) Rectangle(p0, p1 Point) *ShapeBuilder {
	b := &sb.rectangle
	b.p0, b.p1 = p0, p1
	if b.p0.X > b.p1.X {
		b.p1.X, b.p0.X = b.p0.X, b.p1.X
	}
	if b.p0.Y > b.p1.Y {
		b.p1.Y, b.p0.Y = b.p0.Y, b.p1.Y
	}
	sb.active = b
	return sb
}

func (b *RectangleBuilder) Draw(drawPixel drawPixel) {
	for x := b.p0.X; x < b.p1.X; x++ {
		for y := b.p0.Y; y < b.p1.Y; y++ {
			drawPixel(x, y)
		}
	}
}

func (sb *ShapeBuilder) Circle(center Point, diameter int16) *ShapeBuilder {
	b := &sb.circle
	b.center = center
	b.diameter = diameter
	sb.active = b
	return sb
}

func (b *ShapeBuilder) DrawShade(shade uint8) {
	if b.active == nil {
		return
	}
	b.shade = shade
	if shade == 0xff {
		b.active.Draw(b.drawPaper)
	} else if shade == 0x00 {
		b.active.Draw(b.drawInk)
	} else {
		b.active.Draw(b.drawDither)
	}
	b.active = nil
}

func (b *ShapeBuilder) Draw(on bool) {
	if b.active == nil {
		return
	}
	if !on {
		b.active.Draw(b.drawPaper)
	} else {
		b.active.Draw(b.drawInk)
	}
	b.active = nil
}

func (b *ShapeBuilder) DrawGradient(gradient Gradient) {
	if b.active == nil {
		return
	}
	b.gradient = gradient
	b.active.Draw(b.drawGradient)
	b.active = nil
}

func (b *CircleBuilder) Draw(drawPixel drawPixel) {
	d2 := int(b.diameter) * int(b.diameter)
	offset := int((1 + b.diameter) % 2)
	radius := int((1 + b.diameter) / 2)

	for j := 0; j < radius; j++ {
		h2 := 4 * j * j
		for i := 0; i < radius; i++ {
			w2 := 4 * i * i
			if h2+w2 >= d2 {
				break
			}
			x0 := (b.center.X - i) - offset
			y0 := (b.center.Y - j) - offset
			x1 := b.center.X + i
			y1 := b.center.Y + j

			drawPixel(x0, y0)
			drawPixel(x1, y0)
			drawPixel(x0, y1)
			drawPixel(x1, y1)
		}
	}
}

type TriangleBuilder struct {
	p0, p1, p2 Point
}

func (sb *ShapeBuilder) Triangle(p0, p1, p2 Point) *ShapeBuilder {
	b := &sb.triangle
	b.p0 = p0
	b.p1 = p1
	b.p2 = p2
	sb.active = b
	return sb
}

func (b *TriangleBuilder) Draw(drawPixel drawPixel) {
	filledTriangle(drawPixel, b.p0, b.p1, b.p2)
}

func hLine(drawPixel drawPixel, x0, x1, y int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	for x := x0; x <= x1; x++ {
		drawPixel(x, y)
	}
}

func filledTriangle(drawPixel drawPixel, p0, p1, p2 Point) {
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

		hLine(drawPixel, x0, x1, y)
	}
}

func (s *ShapeBuilder) drawInk(x, y int) {
	s.c.Set(x, y, true)
}

func (s *ShapeBuilder) drawPaper(x, y int) {
	s.c.Set(x, y, false)
}

func (s *ShapeBuilder) drawDither(x, y int) {
	s.c.Set(x, y, Dither(x, y, s.shade))
}

func (s *ShapeBuilder) drawGradient(x, y int) {
	s.c.Set(x, y, Dither(x, y, s.gradient.Calculate(x, y)))
}
