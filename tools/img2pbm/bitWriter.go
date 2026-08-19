package main

type bitWriter struct {
	index   int
	current byte
	pix     []byte
}

func NewBitWriter(w, h int) *bitWriter {
	s := ((w + 7) / 8) * h
	return &bitWriter{
		index:   0,
		current: 0,
		pix:     make([]byte, 0, s),
	}
}

func (w *bitWriter) Put(val bool) {
	if val {
		w.current |= (0x1 << (7 - w.index))
	}
	w.index++
	if w.index >= 8 {
		w.EndRow()
	}
}

func (w *bitWriter) EndRow() {
	if w.index != 0 {
		w.pix = append(w.pix, w.current)
		w.current = 0x00
		w.index = 0
	}
}
