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

func (p *Platform) measureCpu() {
	if math.IsNaN(p.debug.cpu) {
		p.debug.cpu = 0.0
	}
	if c, err := p.debug.proc.Percent(0); err == nil {
		p.debug.cpu = c/128.0 + p.debug.cpu*127.0/128.0
	}
}
