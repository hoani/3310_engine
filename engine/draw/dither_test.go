package draw_test

import (
	"fmt"
	"testing"

	"github.com/hoani/3310_engine/engine/draw"
	"github.com/stretchr/testify/require"
)

func TestDitherSprite(t *testing.T) {
	testCases := []struct {
		x, y  []int
		shade uint8
		draw  bool
		value bool
	}{
		{x: []int{0, 1, 2, 3}, y: []int{0, 0, 0, 0}, shade: 0xFF, draw: false, value: false},
		{x: []int{0, 1, 2, 3}, y: []int{0, 0, 0, 0}, shade: 0xF9, draw: true, value: false},
		{x: []int{0, 1, 2, 3}, y: []int{0, 0, 0, 0}, shade: 0xFA, draw: true, value: false},
		// 50%
		{x: []int{0, 1, 2, 3}, y: []int{1, 0, 1, 0}, shade: 0x80, draw: true, value: false},
		{x: []int{0, 1, 2, 3}, y: []int{0, 1, 0, 1}, shade: 0x80, draw: true, value: true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			for i, x := range tc.x {
				y := tc.y[i]
				d, v := draw.DitherSprite(x, y, tc.shade)
				require.Equal(t, d, tc.draw)
				require.Equal(t, v, tc.value)
			}
		})
	}
}
