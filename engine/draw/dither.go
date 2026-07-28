package draw

const ShadeWhite uint8 = 0xF9
const ShadeBlack uint8 = 0x00
const ShadeTransparent uint8 = 0xFF

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

func Dither(x, y int, shade uint8) (draw bool, value bool) {
	if shade == ShadeTransparent {
		return false, false
	}
	if shade >= ShadeWhite {
		return true, false
	}
	level := (shade >> 3) & 0x1f
	inkCoverage := 64 - 2*int(level)
	return true, int(bayer8[y&7][x&7]) <= inkCoverage
}
