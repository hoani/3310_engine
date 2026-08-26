//go:build js && wasm

package desktop

import (
	"math"

	"github.com/shirou/gopsutil/v4/process"
)

func NewProcess() (*process.Process, error) {
	return &process.Process{}, nil
}

func (d *debugExtension) measureCpu() {
	d.cpu = math.NaN()
}
