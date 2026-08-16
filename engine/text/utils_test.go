package text

import (
	"testing"

	"github.com/hoani/3310_engine/font"
	"github.com/stretchr/testify/require"
)

func TestFit(t *testing.T) {
	result := Fit(&font.Tiny, "A quick brown fox jumps over", 48)
	require.Equal(t, "A quick brown\nfox jumps\nover", result)
}
