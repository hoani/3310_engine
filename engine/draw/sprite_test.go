package draw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpriteOptsIdentity(t *testing.T) {

	opt := SpriteOpts{}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			xpos, ypos := opt.Transform(x, y, 8, 8)

			assert.Equal(t, x, xpos)
			assert.Equal(t, y, ypos)
		}
	}

}

func TestSpriteOptsFlip(t *testing.T) {
	testCases := []struct {
		name string
		in   int
		out  int
		w    int
	}{
		{name: "basic", in: 0, out: 7, w: 8},
		{name: "odd", in: 0, out: 8, w: 9},
		{name: "even center", in: 4, out: 3, w: 8},
		{name: "odd center", in: 4, out: 4, w: 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			t.Run("HFlip", func(t *testing.T) {
				opt := SpriteOpts{
					HFlip: true,
				}
				for y := range 8 {
					xpos, ypos := opt.Transform(tc.in, y, tc.w, 8)

					assert.Equal(t, tc.out, xpos)
					assert.Equal(t, y, ypos)
				}
			})

			t.Run("VFlip", func(t *testing.T) {
				opt := SpriteOpts{
					VFlip: true,
				}
				for x := range 8 {
					xpos, ypos := opt.Transform(x, tc.in, 8, tc.w)

					assert.Equal(t, x, xpos)
					assert.Equal(t, tc.out, ypos)
				}
			})

			t.Run("Both Flip", func(t *testing.T) {
				opt := SpriteOpts{
					HFlip: true,
					VFlip: true,
				}
				xpos, ypos := opt.Transform(tc.in, tc.in, tc.w, tc.w)

				assert.Equal(t, tc.out, xpos)
				assert.Equal(t, tc.out, ypos)
			})
		})
	}
}

func TestSpriteOptsRotate(t *testing.T) {
	testCases := []struct {
		name   string
		w      int
		xi, yi int
		xo, yo int
		rot    Rotation
	}{
		{name: "zero", w: 8, xi: 0, yi: 0, xo: 0, yo: 0, rot: Rot0},
		{name: "90", w: 8, xi: 0, yi: 0, xo: 0, yo: 7, rot: Rot90},
		{name: "180", w: 8, xi: 0, yi: 0, xo: 7, yo: 7, rot: Rot180},
		{name: "270", w: 8, xi: 0, yi: 0, xo: 7, yo: 0, rot: Rot270},
		{name: "c90", w: 8, xi: 3, yi: 3, xo: 3, yo: 4, rot: Rot90},
		{name: "c180", w: 8, xi: 3, yi: 3, xo: 4, yo: 4, rot: Rot180},
		{name: "c270", w: 8, xi: 3, yi: 3, xo: 4, yo: 3, rot: Rot270},
		{name: "odd c90", w: 7, xi: 3, yi: 3, xo: 3, yo: 3, rot: Rot90},
		{name: "odd c180", w: 7, xi: 3, yi: 3, xo: 3, yo: 3, rot: Rot180},
		{name: "odd c270", w: 7, xi: 3, yi: 3, xo: 3, yo: 3, rot: Rot270},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			opt := SpriteOpts{
				Rotation: tc.rot,
			}
			xpos, ypos := opt.Transform(tc.xi, tc.yi, tc.w, tc.w)

			assert.Equal(t, tc.xo, xpos)
			assert.Equal(t, tc.yo, ypos)
		})
	}
}
