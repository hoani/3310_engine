package draw

import "github.com/hoani/3310_engine/engine"

// Thresholds for 8x8 bayer dithering
var bayer8 = [8][8]uint8{
	{0, 32, 8, 40, 2, 34, 10, 42},
	{48, 16, 56, 24, 50, 18, 58, 26},
	{12, 44, 4, 36, 14, 46, 6, 38},
	{60, 28, 52, 20, 62, 30, 54, 22},
	{3, 35, 11, 43, 1, 33, 9, 41},
	{51, 19, 59, 27, 49, 17, 57, 25},
	{15, 47, 7, 39, 13, 45, 5, 37},
	{63, 31, 55, 23, 61, 29, 53, 21},
}

func Bayer64(x, y int, amount uint8) bool {
	x = x & 0x07
	y = y & 0x07
	return bayer8[y][x] <= amount
}

func DitherSprite(x, y int, sample uint8) (draw bool, value bool) {
	if sample == engine.SpriteTransparent {
		return false, false
	}
	if sample >= engine.SpriteWhite {
		return true, false
	}
	level := ((sample >> 3) & 0x1f)
	inkCoverage := 64 - 2*level
	return true, Bayer64(x, y, inkCoverage)
}

func Dither(x, y int, shade uint8) (active bool) {
	if shade >= 0xfc {
		return false
	}
	if shade == 0x00 {
		return true
	}
	level := (shade >> 2)
	inkCoverage := 63 - level
	return Bayer64(x, y, inkCoverage)
}

func DitherOffset(x, y int, shade uint8) (active bool) {
	return Dither((x + y + 1), (x - y), shade)
}
