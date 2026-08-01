package draw_test

import (
	"testing"

	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/platform/desktop"
	"github.com/stretchr/testify/require"
)

func TestDrawTriangle(t *testing.T) {
	c := desktop.NewCanvas()
	d := draw.New(c)

	d.Triangle(draw.P(0, 0), draw.P(3, 3), draw.P(0, 3)).Draw(true)

	require.Equal(t, c.Get(0, 0), true)
	require.Equal(t, c.Get(3, 3), true)
	require.Equal(t, c.Get(0, 3), true)

	require.Equal(t, c.Get(1, 0), false)
	require.Equal(t, c.Get(3, 2), false)
	require.Equal(t, c.Get(0, 4), false)
}
