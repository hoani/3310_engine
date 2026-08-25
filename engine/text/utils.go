package text

import (
	"strings"

	"github.com/hoani/3310_engine/engine"
	"tinygo.org/x/tinyfont"
)

func LineWidth(f engine.Font, s string) int {
	parts := strings.Split(s, "\n")
	w := 0
	for _, part := range parts {
		pw, _ := tinyfont.LineWidth(f, part)
		if int(pw) > w {
			w = int(pw)
		}
	}
	return w
}

func Fit(f engine.Font, s string, w int) string {
	parts := strings.Split(s, " ")
	var result strings.Builder
	line := ""
	for _, p := range parts {
		if line == "" {
			line = p
			continue
		}
		next := line + " " + p
		if LineWidth(f, next) > w {
			result.WriteString(line)
			result.WriteString("\n")
			line = p
			continue
		}
		line = next
	}
	if line != "" {
		result.WriteString(line)
	}
	return result.String()
}
