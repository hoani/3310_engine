package engine

type Key int

const (
	K1 Key = iota
	K2
	K3
	KA
	K4
	K5
	K6
	KB
	K7
	K8
	K9
	KC
	KStar
	K0
	KHash
	KD
)

const (
	nameK1    = "1"
	nameK2    = "2"
	nameK3    = "3"
	nameK4    = "4"
	nameK5    = "5"
	nameK6    = "6"
	nameK7    = "7"
	nameK8    = "8"
	nameK9    = "9"
	nameKStar = "*"
	nameK0    = "0"
	nameKHash = "#"
	nameKA    = "A"
	nameKB    = "B"
	nameKC    = "C"
	nameKD    = "D"
)

func KeyName(k Key) string {
	switch k {
	case K1:
		return nameK1
	case K2:
		return nameK2
	case K3:
		return nameK3
	case K4:
		return nameK4
	case K5:
		return nameK5
	case K6:
		return nameK6
	case K7:
		return nameK7
	case K8:
		return nameK8
	case K9:
		return nameK9
	case KStar:
		return nameKStar
	case K0:
		return nameK0
	case KHash:
		return nameKHash
	case KA:
		return nameKA
	case KB:
		return nameKB
	case KC:
		return nameKC
	case KD:
		return nameKD
	}
	return ""
}
