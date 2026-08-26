//go:build !(js && wasm)

package desktop

import (
	"math"
	"os"

	"github.com/shirou/gopsutil/v4/process"
)

func NewProcess() (*process.Process, error) {
	return process.NewProcess(int32(os.Getpid()))
}

func (d *debugExtension) measureCpu() {
	if math.IsNaN(d.cpu) {
		d.cpu = 0.0
	}
	if c, err := d.proc.Percent(0); err == nil {
		d.cpu = c/128.0 + d.cpu*127.0/128.0
	}
}
