package procker

import (
	"syscall"
	"time"
)

func (p *sysProcess) stop(timeout time.Duration) error {
	return p.signal(syscall.SIGKILL)
}
