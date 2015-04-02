package procker

import (
	"syscall"
	"time"
)

func (p *SysProcess) stop(timeout time.Duration) error {
	return p.Signal(syscall.SIGKILL)
}
