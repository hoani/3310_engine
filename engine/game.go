package engine

import "image/color"

type Canvas interface {
	Clear()
	Width() int
	Height() int
	Set(x, y int, val bool)
	Get(x, y int) bool
	// Font Displayer Methods
	Size() (x, y int16)
	SetPixel(x, y int16, c color.RGBA)
	Display() error
}

type Game interface {
	Fps() int
	Update() error
	Draw(canvas Canvas) error
}
