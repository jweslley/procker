// +build !windows

package procker

import (
	"syscall"
	"time"
)

func (p *SysProcess) stop(timeout time.Duration) error {
	p.Signal(syscall.SIGTERM)

	select {
	case err := <-p.errc:
		return err
	case <-time.After(timeout):
		p.Signal(syscall.SIGKILL)
		return <-p.errc
	}
}
