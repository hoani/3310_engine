package surface

import (
	"testing"

	"github.com/hoani/3310_engine/engine"
	"github.com/stretchr/testify/assert"
)

func TestSurface(t *testing.T) {
	testCases := []struct {
		name string
		ink  bool
	}{
		{name: "ink", ink: true},
		{name: "paper", ink: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("basic", func(t *testing.T) {
				c := engine.NewCanvas(8, 8)
				c.Clear(!tc.ink)

				s := New(8, 8)
				s.Set(0, 0, tc.ink)

				c.DrawSurface(0, 0, s)

				assert.Equal(t, c.Get(0, 0), tc.ink)
				assert.Equal(t, c.Get(1, 0), !tc.ink)
			})

			t.Run("x offset", func(t *testing.T) {
				c := engine.NewCanvas(8, 8)
				c.Clear(!tc.ink)

				s := New(8, 8)
				s.Set(0, 0, tc.ink)

				c.DrawSurface(2, 0, s)

				assert.Equal(t, c.Get(2, 0), tc.ink)
				assert.Equal(t, c.Get(0, 0), !tc.ink)
			})

			t.Run("y offset", func(t *testing.T) {
				c := engine.NewCanvas(16, 16)
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
			})

			t.Run("different sizes", func(t *testing.T) {
				c := engine.NewCanvas(32, 32)
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
			})
		})
	}
}
