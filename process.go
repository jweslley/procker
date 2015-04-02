package procker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Process manages process's lifecycle.
type Process interface {

	// Start starts the process but does not wait for it to complete.
	Start() error

	// Stop stops the process gracefully within a given timeout.
	// Stop will kill the process if timeout expires.
	Stop(timeout time.Duration) error

	// Running reports whether the process was running.
	Running() bool

	// Signal sends a signal to the Process.
	Signal(os.Signal) error

	// Wait waits for the command to exit.
	Wait() error
}

// SysProcess represents an external command.
// Please check exec.Cmd for more information about exported fields.
type SysProcess struct {
	Command     string
	Dir         string
	Env         []string
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	ExtraFiles  []*os.File
	SysProcAttr *syscall.SysProcAttr

	cmd  *exec.Cmd
	errc chan error
}

func (p *SysProcess) Start() error {
	if p.Running() {
		return errors.New("procker: already started")
	}

	args := strings.Fields(p.expandedCmd())
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = p.Env
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	p.cmd.ExtraFiles = p.ExtraFiles
	p.cmd.SysProcAttr = p.SysProcAttr

	if p.errc == nil {
		p.errc = make(chan error)
	}

	err := p.cmd.Start()
	if err != nil {
		return fmt.Errorf("procker: failed to start: %v", err)
	}

	go func() {
		p.errc <- p.cmd.Wait()
		p.cmd = nil
	}()
	return nil
}

func (p *SysProcess) Stop(timeout time.Duration) error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return p.stop(timeout)
}

func (p *SysProcess) Signal(sig os.Signal) error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return p.cmd.Process.Signal(sig)
}

func (p *SysProcess) Wait() error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return <-p.errc
}

func (p *SysProcess) Running() bool {
	return p.cmd != nil
}

func (p *SysProcess) expandedCmd() string {
	m := env2Map(p.Env)
	return os.Expand(p.Command, func(name string) string {
		return m[name]
	})
}

func (p *SysProcess) String() string {
	return p.Command
}

// FIXME processGroup is too simplistic, needs improvements
type processGroup struct {
	running   bool
	processes []Process
}

// NewProcessGroup creates a process which controls other processes.
func NewProcessGroup(processes ...Process) Process {
	return &processGroup{processes: processes}
}

func (pg *processGroup) Start() error {
	if pg.Running() {
		return errors.New("procker: already started")
	}

	pg.running = true
	return pg.each(func(p Process) error {
		return p.Start()
	})
}

func (pg *processGroup) Stop(timeout time.Duration) error {
	if !pg.Running() {
		return errors.New("procker: not started")
	}

	return pg.each(func(p Process) error {
		return p.Stop(timeout)
	})
}

func (pg *processGroup) Signal(sig os.Signal) error {
	if !pg.Running() {
		return errors.New("procker: not started")
	}

	return pg.each(func(p Process) error {
		return p.Signal(sig)
	})
}

func (pg *processGroup) Wait() error {
	if !pg.Running() {
		return errors.New("procker: not started")
	}

	return pg.each(func(p Process) error {
		return p.Wait()
	})
}

func (pg *processGroup) Running() bool {
	return pg.running
}

func (pg *processGroup) each(f func(p Process) error) error {
	var err error
	var wg sync.WaitGroup
	for _, process := range pg.processes {
		wg.Add(1)
		go func(p Process) {
			err = f(p)
			wg.Done()
		}(process)
	}
	wg.Wait()
	return err
}
