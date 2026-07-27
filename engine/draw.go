package engine

import (
	"image/color"

	"tinygo.org/x/tinyfont"
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
	HLine(x0, x1, y int, on bool)
	FillTriangle(p0, p1, p2 Point, on bool)
	Sprite(x, y int, spr Sprite, index int, invert bool)
	Text(x, y int, font tinyfont.Fonter, str string, on bool)
}

type draw struct {
	c  Canvas
	fc *FontCanvas
}

func NewDraw(c Canvas) Draw {
	return &draw{
		c:  c,
		fc: &FontCanvas{c: c, col: PixelOff},
	}
}

func AbsInt(val int) int {
	if val >= 0 {
		return val
	}
	return -val
}

func (d *draw) Sprite(x, y int, spr Sprite, index int, invert bool) {
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

func (d *draw) Text(x, y int, font tinyfont.Fonter, str string, on bool) {
	col := PixelOff
	if on {
		col = PixelOn
	}
	// Note: we use a font canvas to set the color because the generated Fonter ignores our color.
	d.fc.SetColor(col)
	tinyfont.WriteLine(d.fc, font, int16(x), int16(y), str, col)
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

func (d *draw) HLine(x0, x1, y int, on bool) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	for x := x0; x <= x1; x++ {
		d.c.Set(x, y, on)
	}
}

func (d *draw) FillTriangle(p0, p1, p2 Point, on bool) {
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

		d.HLine(x0, x1, y, on)
	}
}
