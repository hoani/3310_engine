package engine

const SpriteWhite uint8 = 0xF9
const SpriteBlack uint8 = 0x00
const SpriteTransparent uint8 = 0xFF

type Sprite interface {
	At(i, j, index int) (shade uint8)
	Width() int
	Height() int
	Count() int
}
