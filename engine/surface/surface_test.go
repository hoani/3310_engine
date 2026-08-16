package surface

import (
	"testing"

	"github.com/hoani/3310_engine/engine"
	"github.com/stretchr/testify/assert"
)

func TestSurfaceReset(t *testing.T) {
	s := New(8, 8)
	s.Clear(true)

	for _, b := range s.Ink() {
		assert.Equal(t, uint8(0xFF), b)
	}

	s.Reset()

	for _, b := range s.Ink() {
		assert.Equal(t, uint8(0x00), b)
	}

	s.Clear(false)

	for _, b := range s.Paper() {
		assert.Equal(t, uint8(0xFF), b)
	}

	s.Reset()

	for _, b := range s.Paper() {
		assert.Equal(t, uint8(0x00), b)
	}
}

func TestSurface(t *testing.T) {
	testCases := []struct {
		name       string
		ink        bool
		canvasFunc func(x, y int) engine.Canvas
	}{
		{name: "ink", ink: true, canvasFunc: engine.NewCanvas},
		{name: "paper", ink: false, canvasFunc: engine.NewCanvas},
		{name: "surface x surface", ink: true, canvasFunc: func(x, y int) engine.Canvas { return New(x, y) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			t.Run("clear", func(t *testing.T) {
				s := New(8, 8)
				s.Clear(tc.ink)

				for i := range s.Width() {
					for j := range s.Height() {
						assert.Equal(t, s.Get(i, j), tc.ink)
					}
				}
			})

			t.Run("basic", func(t *testing.T) {
				c := tc.canvasFunc(8, 8)
				c.Clear(!tc.ink)

				s := New(8, 8)
				s.Set(0, 0, tc.ink)

				c.DrawSurface(0, 0, s)

				assert.Equal(t, c.Get(0, 0), tc.ink)
				assert.Equal(t, c.Get(1, 0), !tc.ink)
			})

			t.Run("x offset", func(t *testing.T) {
				c := tc.canvasFunc(8, 8)
				c.Clear(!tc.ink)

				s := New(8, 8)
				s.Set(0, 0, tc.ink)

				c.DrawSurface(2, 0, s)

				assert.Equal(t, c.Get(2, 0), tc.ink)
				assert.Equal(t, c.Get(0, 0), !tc.ink)

				// negative value
				s.Set(7, 0, tc.ink)
				c.DrawSurface(-2, 0, s)

				assert.Equal(t, c.Get(5, 0), tc.ink)
			})

			t.Run("y offset", func(t *testing.T) {
				c := tc.canvasFunc(16, 16)
				c.Clear(!tc.ink)

				s := New(16, 16)
				s.Set(0, 0, tc.ink)

				c.DrawSurface(0, 2, s)

				assert.Equal(t, c.Get(0, 2), tc.ink)
				assert.Equal(t, c.Get(0, 0), !tc.ink)

				s.Set(0, 7, tc.ink)
				c.DrawSurface(0, 2, s) // Have it wrap over to the next block

				assert.Equal(t, c.Get(0, 9), tc.ink)
				assert.Equal(t, c.Get(0, 7), !tc.ink)

				// negative value
				c.DrawSurface(0, -2, s)

				assert.Equal(t, c.Get(0, 5), tc.ink)
			})

			t.Run("different sizes", func(t *testing.T) {
				c := tc.canvasFunc(32, 32)
				c.Clear(!tc.ink)

				s := New(16, 16)
				s.Set(8, 8, tc.ink)

				c.DrawSurface(4, 4, s)

				assert.Equal(t, c.Get(12, 12), tc.ink)
				assert.Equal(t, c.Get(0, 0), !tc.ink)

				c.DrawSurface(18, 18, s)

				assert.Equal(t, c.Get(26, 26), tc.ink)
				assert.Equal(t, c.Get(0, 0), !tc.ink)

				// Drawing surface off canvas
				s.Set(15, 0, tc.ink)

				c.DrawSurface(32, 0, s)
				assert.Equal(t, c.Get(15, 8), !tc.ink)

				c.DrawSurface(-16, 8, s)
				assert.Equal(t, c.Get(31, 0), !tc.ink)

				c.DrawSurface(0, -8, s)
				assert.Equal(t, c.Get(15, 0), !tc.ink)
			})
		})
	}
}
