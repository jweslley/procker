package procker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	signal(os.Signal) error

	// Wait waits for the command to exit.
	wait() error
}

type sysProcess struct {
	command string
	dir     string
	env     []string
	stdout  io.Writer
	stderr  io.Writer
	cmd     *exec.Cmd
	errc    chan error
}

// NewProcess creates a new process with the specified arguments.
func NewProcess(
	command, dir string,
	env []string,
	stdout, stderr io.Writer) Process {

	return &sysProcess{
		command: command,
		dir:     dir,
		env:     env,
		stdout:  stdout,
		stderr:  stderr,
		errc:    make(chan error),
	}
}

func (p *sysProcess) Start() error {
	if p.Running() {
		return errors.New("procker: already started")
	}

	args := strings.Fields(p.expandedCmd(p.env))
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Dir = p.dir
	p.cmd.Env = p.env
	p.cmd.Stdout = p.stdout
	p.cmd.Stderr = p.stderr
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

func (p *sysProcess) Stop(timeout time.Duration) error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return p.stop(timeout)
}

func (p *sysProcess) signal(sig os.Signal) error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return p.cmd.Process.Signal(sig)
}

func (p *sysProcess) wait() error {
	if !p.Running() {
		return errors.New("procker: not started")
	}

	return <-p.errc
}

func (p *sysProcess) Running() bool {
	return p.cmd != nil
}

func (p *sysProcess) expandedCmd(env []string) string {
	m := env2Map(env)
	return os.Expand(p.command, func(name string) string {
		return m[name]
	})
}

func (p *sysProcess) String() string {
	return p.command
}

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

func (pg *processGroup) signal(sig os.Signal) error {
	if !pg.Running() {
		return errors.New("procker: not started")
	}

	return pg.each(func(p Process) error {
		return p.signal(sig)
	})
}

func (pg *processGroup) wait() error {
	if !pg.Running() {
		return errors.New("procker: not started")
	}

	return pg.each(func(p Process) error {
		return p.wait()
	})
}

func (pg *processGroup) Running() bool {
	return pg.running
}

func (pg *processGroup) each(func(p Process) error) error {
	var err error
	//for _, process := range pg.processes {
	//err = f(process)
	//}
	return err
}
