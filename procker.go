package procker

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Process manages process's lifecycle.
type Process interface {
	// Start starts the process but does not wait for it to complete.
	Start() error

	// Wait waits for the command to exit.
	// It must have been started by Start.
	Wait() error

	// Kill causes the Process to exit immediately.
	// It must have been started by Start.
	Kill() error

	// Started reports whether the process was started.
	Started() bool
}

type sysProcess struct {
	name    string
	command string
	dir     string
	env     []string
	stdout  io.Writer
	stderr  io.Writer
	cmd     *exec.Cmd
}

// NewProcess creates a new process with the specified arguments.
func NewProcess(
	name, command, dir string,
	env []string,
	stdout, stderr io.Writer) Process {

	return &sysProcess{
		name:    name,
		command: command,
		dir:     dir,
		env:     env,
		stdout:  stdout,
		stderr:  stderr}
}

func (p *sysProcess) Start() error {
	if p.Started() {
		return errors.New("procker: already started")
	}

	args := strings.Fields(p.expandedCmd(p.env))
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Dir = p.dir
	p.cmd.Env = p.env
	p.cmd.Stdout = p.stdout
	p.cmd.Stderr = p.stderr
	return p.cmd.Start()
}

func (p *sysProcess) Wait() error {
	if !p.Started() {
		return errors.New("procker: not started")
	}
	return p.cmd.Wait()
}

func (p *sysProcess) Kill() error {
	if !p.Started() {
		return errors.New("procker: not started")
	}
	return p.cmd.Process.Kill()
}

func (p *sysProcess) Started() bool {
	return p.cmd != nil
}

func (p *sysProcess) expandedCmd(env []string) string {
	m := env2Map(env)
	return os.Expand(p.command, func(name string) string {
		return m[name]
	})
}

func (p *sysProcess) String() string {
	return p.name
}

type processSet struct {
	started   bool
	processes []Process
}

// NewProcessSet creates a process which controls other processes.
func NewProcessSet(processes ...Process) Process {
	return &processSet{processes: processes}
}

func (ps *processSet) Start() error {
	if ps.Started() {
		return errors.New("procker: already started")
	}

	ps.started = true
	return ps.each(func(p Process) error {
		return p.Start()
	})
}

func (ps *processSet) Wait() error {
	if !ps.Started() {
		return errors.New("procker: not started")
	}

	return ps.each(func(p Process) error {
		return p.Wait()
	})
}

func (ps *processSet) Kill() error {
	if !ps.Started() {
		return errors.New("procker: not started")
	}

	return ps.each(func(p Process) error {
		return p.Kill()
	})
}

func (ps *processSet) Started() bool {
	return ps.started
}

func (ps *processSet) each(f func(p Process) error) error {
	var err error
	for _, process := range ps.processes {
		err = f(process)
	}
	return err
}
