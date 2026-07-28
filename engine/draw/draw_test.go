package draw_test

import (
	"testing"

	"github.com/hoani/3310_engine/engine/draw"
	"github.com/hoani/3310_engine/platform/desktop"
	"github.com/stretchr/testify/assert"
)

func TestDrawHline(t *testing.T) {
	c := desktop.NewCanvas()
	d := draw.New(c)

	d.HLine(3, 5, 0, true)

	assert.Equal(t, c.Get(2, 0), false)
	assert.Equal(t, c.Get(3, 0), true)
	assert.Equal(t, c.Get(5, 0), true)
	assert.Equal(t, c.Get(6, 0), false)
	assert.Equal(t, c.Get(5, 1), false)
}

func TestDrawTriangle(t *testing.T) {
	c := desktop.NewCanvas()
	d := draw.New(c)

	d.FillTriangle(draw.P(0, 0), draw.P(3, 3), draw.P(0, 3), true)

	assert.Equal(t, c.Get(0, 0), true)
	assert.Equal(t, c.Get(3, 3), true)
	assert.Equal(t, c.Get(0, 3), true)

	assert.Equal(t, c.Get(1, 0), false)
	assert.Equal(t, c.Get(3, 2), false)
	assert.Equal(t, c.Get(0, 4), false)
}
