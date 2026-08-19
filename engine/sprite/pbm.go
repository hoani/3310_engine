package sprite

import (
	"errors"
	"io"
	"strconv"
	"strings"
)

func nextLine(raw string, off int) (line string, next int, err error) {
	i := strings.IndexByte(raw[off:], '\n')
	if i < 0 {
		return "", 0, io.EOF
	}
	end := off + i
	line = raw[off:end]
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line, end + 1, nil
}

func parseSizes(line string) (int, int, error) {
	sp := strings.IndexByte(line, ' ')
	if sp < 0 {
		return 0, 0, errors.New("sizes are invalid")
	}
	w, err := strconv.Atoi(line[:sp])
	if err != nil {
		return 0, 0, err
	}
	h, err := strconv.Atoi(line[sp+1:])
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}
