package note

// Range is C0..B8 (16 Hz - 7902 Hz), nine full octaves.

type Index uint8

// Names are scientific pitch notation; "s" means sharp (Cs4 == C#4).
const (
	None Index = iota
	C0
	Cs0
	D0
	Ds0
	E0
	F0
	Fs0
	G0
	Gs0
	A0
	As0
	B0
	C1
	Cs1
	D1
	Ds1
	E1
	F1
	Fs1
	G1
	Gs1
	A1
	As1
	B1
	C2
	Cs2
	D2
	Ds2
	E2
	F2
	Fs2
	G2
	Gs2
	A2
	As2
	B2
	C3
	Cs3
	D3
	Ds3
	E3
	F3
	Fs3
	G3
	Gs3
	A3
	As3
	B3
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
	C7
	Cs7
	D7
	Ds7
	E7
	F7
	Fs7
	G7
	Gs7
	A7
	As7
	B7
	C8
	Cs8
	D8
	Ds8
	E8
	F8
	Fs8
	G8
	Gs8
	A8
	As8
	B8
	total
	Custom
)

const Total = int(total)

const midiOffset Index = 11

// Note frequencies for the equal-tempered scale, A4 = 440 Hz (ISO 16). freq[n] = 440 * 2^((m-69)/12), where m is the MIDI note number.
var freq = [total]float32{
	1046.50, // None - frequency doesn't matter, we set the amp to 0
	16.35,   // C0   MIDI 12
	17.32,   // C#0  MIDI 13
	18.35,   // D0   MIDI 14
	19.45,   // D#0  MIDI 15
	20.60,   // E0   MIDI 16
	21.83,   // F0   MIDI 17
	23.12,   // F#0  MIDI 18
	24.50,   // G0   MIDI 19
	25.96,   // G#0  MIDI 20
	27.50,   // A0   MIDI 21
	29.14,   // A#0  MIDI 22
	30.87,   // B0   MIDI 23
	32.70,   // C1   MIDI 24
	34.65,   // C#1  MIDI 25
	36.71,   // D1   MIDI 26
	38.89,   // D#1  MIDI 27
	41.20,   // E1   MIDI 28
	43.65,   // F1   MIDI 29
	46.25,   // F#1  MIDI 30
	49.00,   // G1   MIDI 31
	51.91,   // G#1  MIDI 32
	55.00,   // A1   MIDI 33
	58.27,   // A#1  MIDI 34
	61.74,   // B1   MIDI 35
	65.41,   // C2   MIDI 36
	69.30,   // C#2  MIDI 37
	73.42,   // D2   MIDI 38
	77.78,   // D#2  MIDI 39
	82.41,   // E2   MIDI 40
	87.31,   // F2   MIDI 41
	92.50,   // F#2  MIDI 42
	98.00,   // G2   MIDI 43
	103.83,  // G#2  MIDI 44
	110.00,  // A2   MIDI 45
	116.54,  // A#2  MIDI 46
	123.47,  // B2   MIDI 47
	130.81,  // C3   MIDI 48
	138.59,  // C#3  MIDI 49
	146.83,  // D3   MIDI 50
	155.56,  // D#3  MIDI 51
	164.81,  // E3   MIDI 52
	174.61,  // F3   MIDI 53
	185.00,  // F#3  MIDI 54
	196.00,  // G3   MIDI 55
	207.65,  // G#3  MIDI 56
	220.00,  // A3   MIDI 57
	233.08,  // A#3  MIDI 58
	246.94,  // B3   MIDI 59
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
	2093.00, // C7   MIDI 96
	2217.46, // C#7  MIDI 97
	2349.32, // D7   MIDI 98
	2489.02, // D#7  MIDI 99
	2637.02, // E7   MIDI 100
	2793.83, // F7   MIDI 101
	2959.96, // F#7  MIDI 102
	3135.96, // G7   MIDI 103
	3322.44, // G#7  MIDI 104
	3520.00, // A7   MIDI 105
	3729.31, // A#7  MIDI 106
	3951.07, // B7   MIDI 107
	4186.01, // C8   MIDI 108
	4434.92, // C#8  MIDI 109
	4698.64, // D8   MIDI 110
	4978.03, // D#8  MIDI 111
	5274.04, // E8   MIDI 112
	5587.65, // F8   MIDI 113
	5919.91, // F#8  MIDI 114
	6271.93, // G8   MIDI 115
	6644.88, // G#8  MIDI 116
	7040.00, // A8   MIDI 117
	7458.62, // A#8  MIDI 118
	7902.13, // B8   MIDI 119
}

func Name(i Index) string {
	switch i {
	case None:
		return "None"
	case C0:
		return "C0"
	case Cs0:
		return "Cs0"
	case D0:
		return "D0"
	case Ds0:
		return "Ds0"
	case E0:
		return "E0"
	case F0:
		return "F0"
	case Fs0:
		return "Fs0"
	case G0:
		return "G0"
	case Gs0:
		return "Gs0"
	case A0:
		return "A0"
	case As0:
		return "As0"
	case B0:
		return "B0"
	case C1:
		return "C1"
	case Cs1:
		return "Cs1"
	case D1:
		return "D1"
	case Ds1:
		return "Ds1"
	case E1:
		return "E1"
	case F1:
		return "F1"
	case Fs1:
		return "Fs1"
	case G1:
		return "G1"
	case Gs1:
		return "Gs1"
	case A1:
		return "A1"
	case As1:
		return "As1"
	case B1:
		return "B1"
	case C2:
		return "C2"
	case Cs2:
		return "Cs2"
	case D2:
		return "D2"
	case Ds2:
		return "Ds2"
	case E2:
		return "E2"
	case F2:
		return "F2"
	case Fs2:
		return "Fs2"
	case G2:
		return "G2"
	case Gs2:
		return "Gs2"
	case A2:
		return "A2"
	case As2:
		return "As2"
	case B2:
		return "B2"
	case C3:
		return "C3"
	case Cs3:
		return "Cs3"
	case D3:
		return "D3"
	case Ds3:
		return "Ds3"
	case E3:
		return "E3"
	case F3:
		return "F3"
	case Fs3:
		return "Fs3"
	case G3:
		return "G3"
	case Gs3:
		return "Gs3"
	case A3:
		return "A3"
	case As3:
		return "As3"
	case B3:
		return "B3"
	case C4:
		return "C4"
	case Cs4:
		return "Cs4"
	case D4:
		return "D4"
	case Ds4:
		return "Ds4"
	case E4:
		return "E4"
	case F4:
		return "F4"
	case Fs4:
		return "Fs4"
	case G4:
		return "G4"
	case Gs4:
		return "Gs4"
	case A4:
		return "A4"
	case As4:
		return "As4"
	case B4:
		return "B4"
	case C5:
		return "C5"
	case Cs5:
		return "Cs5"
	case D5:
		return "D5"
	case Ds5:
		return "Ds5"
	case E5:
		return "E5"
	case F5:
		return "F5"
	case Fs5:
		return "Fs5"
	case G5:
		return "G5"
	case Gs5:
		return "Gs5"
	case A5:
		return "A5"
	case As5:
		return "As5"
	case B5:
		return "B5"
	case C6:
		return "C6"
	case Cs6:
		return "Cs6"
	case D6:
		return "D6"
	case Ds6:
		return "Ds6"
	case E6:
		return "E6"
	case F6:
		return "F6"
	case Fs6:
		return "Fs6"
	case G6:
		return "G6"
	case Gs6:
		return "Gs6"
	case A6:
		return "A6"
	case As6:
		return "As6"
	case B6:
		return "B6"
	case C7:
		return "C7"
	case Cs7:
		return "Cs7"
	case D7:
		return "D7"
	case Ds7:
		return "Ds7"
	case E7:
		return "E7"
	case F7:
		return "F7"
	case Fs7:
		return "Fs7"
	case G7:
		return "G7"
	case Gs7:
		return "Gs7"
	case A7:
		return "A7"
	case As7:
		return "As7"
	case B7:
		return "B7"
	case C8:
		return "C8"
	case Cs8:
		return "Cs8"
	case D8:
		return "D8"
	case Ds8:
		return "Ds8"
	case E8:
		return "E8"
	case F8:
		return "F8"
	case Fs8:
		return "Fs8"
	case G8:
		return "G8"
	case Gs8:
		return "Gs8"
	case A8:
		return "A8"
	case As8:
		return "As8"
	case B8:
		return "B8"
	default:
		return ""
	}
}
