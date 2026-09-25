//go:build tinygo

package firmware

type Extension interface {
	Setup() error
	Update() error
}
