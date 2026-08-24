package draw

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/font"
)

var PixelOff = color.RGBA{0, 0, 0, 255}
var PixelOn = color.RGBA{255, 255, 255, 255}

type Rotation uint8

const (
	Rot0 Rotation = iota
	Rot90
	Rot180
	Rot270
)

type Point struct {
	X int
	Y int
}

func P(x, y int) Point {
	return Point{X: x, Y: y}
}

type Draw interface {
	Triangle(p0, p1, p2 Point) *ShapeBuilder
	Shape(s ShapeDrawer) *ShapeBuilder
	Circle(center Point, diameter uint16) *ShapeBuilder
	Oval(center Point, width uint16, height uint16) *ShapeBuilder
	Rectangle(x0, y0, x1, y1 int) *ShapeBuilder
	Line(x0, y0, x1, y1 int) *ShapeBuilder
	Sprite(x, y int, spr engine.Sprite, index int, opts *SpriteOpts)
	Text(x, y int, str string) *TextBuilder
}

type draw struct {
	c            engine.Canvas
	textBuilder  *TextBuilder
	shapeBuilder *ShapeBuilder
}

func New(c engine.Canvas) Draw {
	return &draw{
		c:            c,
		textBuilder:  NewTextBuilder(&FontCanvas{c: c, ink: false}, &font.Tiny),
		shapeBuilder: NewShapeBuilder(c),
	}
}

func AbsInt(val int) int {
	if val >= 0 {
		return val
	}
	return -val
}

func (d *draw) Text(x, y int, str string) *TextBuilder {
	return d.textBuilder.New(x, y, str)
}

func (d *draw) Rectangle(x0, y0, x1, y1 int) *ShapeBuilder {
	return d.shapeBuilder.Rectangle(P(x0, y0), P(x1, y1))
}

func (d *draw) Shape(s ShapeDrawer) *ShapeBuilder {
	return d.shapeBuilder.Shape(s)
}

func (d *draw) Line(x0, y0, x1, y1 int) *ShapeBuilder {
	return d.shapeBuilder.Line(P(x0, y0), P(x1, y1))
}

func (d *draw) Triangle(p0, p1, p2 Point) *ShapeBuilder {
	return d.shapeBuilder.Triangle(p0, p1, p2)
}

func (d *draw) Circle(center Point, diameter uint16) *ShapeBuilder {
	return d.shapeBuilder.Circle(center, diameter)
}

func (d *draw) Oval(center Point, width uint16, height uint16) *ShapeBuilder {
	return d.shapeBuilder.Oval(center, width, height)
}
