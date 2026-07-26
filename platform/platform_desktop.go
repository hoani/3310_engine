//go:build !tinygo

package platform

import (
	"github.com/hoani/3310_engine/platform/desktop"
)

func Run() {
	desktop.Run()
}
