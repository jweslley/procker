// +build !windows

package procker

import (
	"syscall"
	"time"
)

func (p *sysProcess) stop(timeout time.Duration) error {
	p.signal(syscall.SIGTERM)

	select {
	case err := <-p.errc:
		return err
	case <-time.After(timeout):
		return p.signal(syscall.SIGKILL)
	}

	return nil
}
