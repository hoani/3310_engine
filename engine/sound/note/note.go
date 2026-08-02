package note

type Index uint8

const (
	None Index = iota
	C1
	C4
	total
)

const Total = int(total)

var freq = [total]float32{
	2093.05, // None - frequency doesn't matter, we set the amp to 0
	523.25,  // C1
	2093.05, // C4
}

var amp = [total]uint8{
	0, // None
	0xFF,
	0xFF,
}

func Frequency(i Index) float32 {
	if i >= total {
		i = None
	}
	return freq[i]
}

func Amplitude(i Index) uint8 {
	if i >= total {
		i = None
	}
	return amp[i]
}
