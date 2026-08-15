package draw

import (
	"github.com/hoani/3310_engine/engine"
)

type RectangleBuilder struct {
	p0 Point
	p1 Point
}

type LineBuilder struct {
	p0 Point
	p1 Point
}

type CircleBuilder struct {
	center   Point
	diameter int16
}

type OvalBuilder struct {
	center Point
	width  int16
	height int16
}

type TriangleBuilder struct {
	p0, p1, p2 Point
}

type TStripBuilder struct {
	points []Point
}

type TFanBuilder struct {
	center Point
	points []Point
}

type drawPixel func(x, y int)

type Gradient interface {
	Calculate(x, y int) uint8
}

type ShapeDrawer interface {
	Fill(drawPixel drawPixel)
	Outline(drawPixel drawPixel)
}

type ShapeBuilder struct {
	c         engine.Canvas
	rectangle RectangleBuilder
	circle    CircleBuilder
	oval      OvalBuilder
	triangle  TriangleBuilder
	line      LineBuilder
	active    ShapeDrawer
	shade     uint8
	gradient  Gradient
	opts      *Opts
}

func NewShapeBuilder(c engine.Canvas) *ShapeBuilder {
	return &ShapeBuilder{
		c:         c,
		rectangle: RectangleBuilder{},
		circle:    CircleBuilder{},
		oval:      OvalBuilder{},
		triangle:  TriangleBuilder{},
		line:      LineBuilder{},
		active:    nil,
		shade:     0x00,
		gradient:  nil,
	}
}

func (b *ShapeBuilder) DrawShade(shade uint8, opts *Opts) {
	if b.active == nil {
		return
	}
	b.opts = opts
	b.shade = shade
	switch shade {
	case 0xff:
		b.active.Fill(b.drawPaper)
	case 0x00:
		b.active.Fill(b.drawInk)
	default:
		b.active.Fill(b.drawDither)
	}

	if opts.Outline.Apply {
		b.active.Outline(b.drawPixelInk(opts.Outline.Ink))
	}
	b.active = nil
}

func (b *ShapeBuilder) Draw(on bool, opts *Opts) {
	if b.active == nil {
		return
	}
	b.opts = opts

	if opts.Outline.Only {
		b.active.Outline(b.drawPixelInk(on))
	} else {
		b.active.Fill(b.drawPixelInk(on))
	}

	if opts.Outline.Apply {
		b.active.Outline(b.drawPixelInk(opts.Outline.Ink))
	}
	b.active = nil
}

func (b *ShapeBuilder) drawPixelInk(on bool) drawPixel {
	if !on {
		return b.drawPaper
	}
	return b.drawInk
}

func (b *ShapeBuilder) DrawGradient(gradient Gradient, opts *Opts) {
	if b.active == nil {
		return
	}
	b.opts = opts
	b.gradient = gradient
	b.active.Fill(b.drawGradient)

	if opts.Outline.Apply {
		b.active.Outline(b.drawPixelInk(opts.Outline.Ink))
	}

	b.active = nil
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

func (sb *ShapeBuilder) Shape(s ShapeDrawer) *ShapeBuilder {
	sb.active = s
	return sb
}

func (b *RectangleBuilder) Fill(drawPixel drawPixel) {
	for x := b.p0.X; x < b.p1.X; x++ {
		for y := b.p0.Y; y < b.p1.Y; y++ {
			drawPixel(x, y)
		}
	}
}

func (b *RectangleBuilder) Outline(drawPixel drawPixel) {
	line(drawPixel, b.p0, Point{b.p1.X, b.p0.Y})
	line(drawPixel, Point{b.p1.X, b.p0.Y}, b.p1)
	line(drawPixel, b.p0, Point{b.p0.X, b.p1.Y})
	line(drawPixel, Point{b.p0.X, b.p1.Y}, b.p1)
}

func (sb *ShapeBuilder) Circle(center Point, diameter int16) *ShapeBuilder {
	b := &sb.circle
	b.center = center
	b.diameter = diameter
	sb.active = b
	return sb
}

func (b *CircleBuilder) draw(drawPixel drawPixel, fill bool) {
	d := int(b.diameter)
	d2 := d * d
	offset := (1 + d) % 2
	radius := (1 + d) / 2

	last := 0

	for j := radius - 1; j >= 0; j-- {
		v := 2*j + offset
		v2 := v * v
		lineDone := false
		for i := 0; i < radius; i++ {
			x0 := (b.center.X - i) - offset
			y0 := (b.center.Y - j) - offset
			x1 := b.center.X + i
			y1 := b.center.Y + j

			if i >= last {
				u := 2*(i+1) + offset
				u2 := u * u
				if u2+v2 > d2 {
					last = i
					lineDone = true
				}
			}

			if fill || j == radius-1 || lineDone || i > last {
				drawPixel(x0, y0)
				drawPixel(x1, y0)
				drawPixel(x0, y1)
				drawPixel(x1, y1)
			}

			if lineDone {
				break
			}
		}
	}
}

func (b *CircleBuilder) Fill(drawPixel drawPixel) {
	b.draw(drawPixel, true)
}

func (b *CircleBuilder) Outline(drawPixel drawPixel) {
	b.draw(drawPixel, false)
}

func (sb *ShapeBuilder) Oval(center Point, width, height int16) *ShapeBuilder {
	b := &sb.oval
	b.center = center
	b.width = width
	b.height = height
	sb.active = b
	return sb
}

func (b *OvalBuilder) draw(drawPixel drawPixel, fill bool) {
	// Uses equation 1 = (x-xc)^2/a^2 + (y-yc)^2/b^2
	// Where:
	//   a = horizontal radius
	//   b = vertical radius
	//   xc, yc = the center point
	//
	// To handle widths and heights better though we actually use:
	// 1 = (2*(x-xc))^2/w^2 + (2*(y-yc))^2/h^2
	// The use of widths and heights allow us to continue using integers

	w := int(b.width)
	h := int(b.height)

	w2 := w * w
	h2 := h * h

	jmax := (h + 1) / 2
	imax := (w + 1) / 2

	xoffset := int((1 + b.width) % 2)
	yoffset := int((1 + b.height) % 2)

	last := 0

	for j := jmax - 1; j >= 0; j-- {
		limit := w2*h2 - 4*j*j*w2
		lineDone := false
		for i := 0; i < imax; i++ {
			if i >= last {
				check := 4 * (i + 1) * (i + 1) * h2
				if check > limit {
					last = i
					lineDone = true
				}
			}

			x0 := (b.center.X - i) - xoffset
			y0 := (b.center.Y - j) - yoffset
			x1 := b.center.X + i
			y1 := b.center.Y + j

			if fill || j == jmax-1 || lineDone || i > last {
				drawPixel(x0, y0)
				drawPixel(x1, y0)
				drawPixel(x0, y1)
				drawPixel(x1, y1)
			}

			if lineDone {
				break // We are done on this line
			}
		}
	}
}

func (b *OvalBuilder) Fill(drawPixel drawPixel) {
	b.draw(drawPixel, true)
}

func (b *OvalBuilder) Outline(drawPixel drawPixel) {
	b.draw(drawPixel, false)
}

func (sb *ShapeBuilder) Triangle(p0, p1, p2 Point) *ShapeBuilder {
	b := &sb.triangle
	b.p0 = p0
	b.p1 = p1
	b.p2 = p2
	sb.active = b
	return sb
}

func (b *TriangleBuilder) Fill(drawPixel drawPixel) {
	filledTriangle(drawPixel, b.p0, b.p1, b.p2)
	b.Outline(drawPixel) // Smooth out some edges
}

func (b *TriangleBuilder) Outline(drawPixel drawPixel) {
	line(drawPixel, b.p0, b.p1)
	line(drawPixel, b.p0, b.p2)
	line(drawPixel, b.p1, b.p2)
}

func (sb *ShapeBuilder) Line(p0, p1 Point) *ShapeBuilder {
	b := &sb.line
	b.p0 = p0
	b.p1 = p1
	sb.active = b
	return sb
}

func (b *LineBuilder) Fill(drawPixel drawPixel) {
	line(drawPixel, b.p0, b.p1)
}

func (b *LineBuilder) Outline(drawPixel drawPixel) {
	line(drawPixel, b.p0, b.p1)
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

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func imin(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func line(drawPixel drawPixel, p0, p1 Point) {
	// sort points by y ascending: (x0,y0) top ... (x1,y1) bottom
	if p0.Y > p1.Y {
		p0, p1 = p1, p0
	}

	height := p1.Y - p0.Y

	if height == 0 {
		x0 := p0.X
		x1 := p1.X
		hLine(drawPixel, x0, x1, p0.Y)
		return
	}

	xmin := imin(p0.X, p1.X)
	xmax := imax(p0.X, p1.X)

	width := p1.X - p0.X
	sign := 1
	if width < 0 {
		sign = -1
	}

	gradient := 0xffff * (width + sign) / (height + 1)

	for y := p0.Y; y <= p1.Y; y++ {

		x0 := p0.X + (gradient*((y)-p0.Y))/0xFFFF
		x1 := p0.X + (gradient*((y+1)-p0.Y))/0xFFFF

		if y < p1.Y {
			if x1 > x0 {
				x1--
			} else if x1 < x0 {
				x1++
			}
		}

		if x0 < x1 {
			x0 = imax(x0, xmin)
			x1 = imin(x1, xmax)
		} else {
			x1 = imax(x1, xmin)
			x0 = imin(x0, xmax)
		}

		hLine(drawPixel, x0, x1, y)
	}
}

func (s *ShapeBuilder) drawInk(x, y int) {
	if !s.opts.Show(x, y) {
		return
	}
	s.c.Set(x, y, !s.opts.Invert)
}

func (s *ShapeBuilder) drawPaper(x, y int) {
	if !s.opts.Show(x, y) {
		return
	}
	s.c.Set(x, y, s.opts.Invert)
}

func (s *ShapeBuilder) drawDither(x, y int) {
	if !s.opts.Show(x, y) {
		return
	}
	on := Dither(x, y, s.shade)
	if s.opts.Invert {
		on = !on
	}
	s.c.Set(x, y, on)
}

func (s *ShapeBuilder) drawGradient(x, y int) {
	if !s.opts.Show(x, y) {
		return
	}
	on := Dither(x, y, s.gradient.Calculate(x, y))
	if s.opts.Invert {
		on = !on
	}
	s.c.Set(x, y, on)
}
