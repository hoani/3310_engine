package note

func Frequency(i Index) float32 {
	if i >= total {
		i = None
	}
	return freq[i]
}

func IndexFromMidi(i uint8) Index {
	if i < uint8(midiOffset) {
		return None
	}
	val := i - uint8(midiOffset)
	if Index(val) > total {
		return None
	}
	return Index(val)
}
