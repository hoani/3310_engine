package draw

import (
	"image/color"

	"github.com/hoani/3310_engine/engine"
	"github.com/hoani/3310_engine/font"
)

var PixelOff = color.RGBA{0, 0, 0, 255}
var PixelOn = color.RGBA{255, 255, 255, 255}

type Point struct {
	X int
	Y int
}

func P(x, y int) Point {
	return Point{X: x, Y: y}
}

type Draw interface {
	Triangle(p0, p1, p2 Point) *TriangleBuilder
	Sprite(x, y int, spr engine.Sprite, index int, invert bool)
	Text(x, y int, str string) *TextBuilder
}

type draw struct {
	c               engine.Canvas
	textBuilder     *TextBuilder
	triangleBuilder *TriangleBuilder
}

func New(c engine.Canvas) Draw {
	return &draw{
		c:               c,
		textBuilder:     NewTextBuilder(&FontCanvas{c: c, ink: false}, &font.Tiny),
		triangleBuilder: NewTriangleBuilder(c),
	}
}

func AbsInt(val int) int {
	if val >= 0 {
		return val
	}
	return -val
}

func (d *draw) Sprite(x, y int, spr engine.Sprite, index int, invert bool) {
	for i := 0; i < spr.Width(); i++ {
		for j := 0; j < spr.Height(); j++ {
			show, on := spr.At(i, j, index)
			if show {
				if invert {
					on = !on
				}
				d.c.Set(x+i, y+j, on)
			}
		}
	}
}

func (d *draw) Text(x, y int, str string) *TextBuilder {
	return d.textBuilder.New(x, y, str)
}

// func (d *Draw) Line(x0, y0, x1, y1 int, w int, c bool) {
// 	dx := x1 - x0
// 	dy := y1 - y0
// 	absDx := AbsInt(dx)
// 	absDy := AbsInt(dy)
// 	if absDx > absDy {
// 		for i := 0; i < absDx; i++ {
// 			x := x0 + dx*i/absDx
// 			y := y0 + dy*i/absDx
// 			d.c.Set(x, y, c)
// 		}
// 	}
// }

func (d *draw) Triangle(p0, p1, p2 Point) *TriangleBuilder {
	return d.triangleBuilder.New(p0, p1, p2)
}
