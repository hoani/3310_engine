package engine

type Canvas interface {
	Clear()
	Width() int
	Height() int
	Set(x, y int, val bool)
	Get(x, y int) bool
}

type Game interface {
	Fps() int
	Update() error
	Draw(canvas Canvas) error
}
