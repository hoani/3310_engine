package note

// Range is C4..B6 (261 Hz - 1976 Hz), three full octaves, sits below the ~2 kHz point where a buzzer starts losing output.
// C4-C5 is thin but legible: square-wave drive puts strong harmonics in the buzzer's efficient band even when the fundamental barely radiates.

type Index uint8

// Names are scientific pitch notation; "s" means sharp (Cs4 == C#4).
const (
	None Index = iota
	C4
	Cs4
	D4
	Ds4
	E4
	F4
	Fs4
	G4
	Gs4
	A4
	As4
	B4
	C5
	Cs5
	D5
	Ds5
	E5
	F5
	Fs5
	G5
	Gs5
	A5
	As5
	B5
	C6
	Cs6
	D6
	Ds6
	E6
	F6
	Fs6
	G6
	Gs6
	A6
	As6
	B6
	total
)

const Total = int(total)

const midiOffset Index = 59

// Note frequencies for the equal-tempered scale, A4 = 440 Hz (ISO 16). freq[n] = 440 * 2^((m-69)/12), where m is the MIDI note number.
var freq = [total]float32{
	1046.50, // None - frequency doesn't matter, we set the amp to 0
	261.63,  // C4   MIDI 60
	277.18,  // C#4  MIDI 61
	293.66,  // D4   MIDI 62
	311.13,  // D#4  MIDI 63
	329.63,  // E4   MIDI 64
	349.23,  // F4   MIDI 65
	369.99,  // F#4  MIDI 66
	392.00,  // G4   MIDI 67
	415.30,  // G#4  MIDI 68
	440.00,  // A4   MIDI 69
	466.16,  // A#4  MIDI 70
	493.88,  // B4   MIDI 71
	523.25,  // C5   MIDI 72
	554.37,  // C#5  MIDI 73
	587.33,  // D5   MIDI 74
	622.25,  // D#5  MIDI 75
	659.26,  // E5   MIDI 76
	698.46,  // F5   MIDI 77
	739.99,  // F#5  MIDI 78
	783.99,  // G5   MIDI 79
	830.61,  // G#5  MIDI 80
	880.00,  // A5   MIDI 81
	932.33,  // A#5  MIDI 82
	987.77,  // B5   MIDI 83
	1046.50, // C6   MIDI 84
	1108.73, // C#6  MIDI 85
	1174.66, // D6   MIDI 86
	1244.51, // D#6  MIDI 87
	1318.51, // E6   MIDI 88
	1396.91, // F6   MIDI 89
	1479.98, // F#6  MIDI 90
	1567.98, // G6   MIDI 91
	1661.22, // G#6  MIDI 92
	1760.00, // A6   MIDI 93
	1864.66, // A#6  MIDI 94
	1975.53, // B6   MIDI 95
}
