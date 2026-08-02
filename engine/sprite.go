package engine

const ShadeWhite uint8 = 0xF9
const ShadeBlack uint8 = 0x00
const ShadeTransparent uint8 = 0xFF

type Sprite interface {
	At(i, j, index int) (shade uint8)
	Width() int
	Height() int
	Count() int
}
