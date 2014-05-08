package procker

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Process interface {
	Start() error
	Wait() error
	Kill() error
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

func NewProcessSet(processes ...Process) Process {
	return &processSet{processes: processes}
}

func (ps *processSet) Start() error {
	if ps.Started() {
		return errors.New("procker: already started")
	}

	var err error
	for _, process := range ps.processes {
		err = process.Start()
	}
	ps.started = true
	return err
}

func (ps *processSet) Wait() error {
	if !ps.Started() {
		return errors.New("procker: not started")
	}

	var err error
	for _, process := range ps.processes {
		err = process.Wait()
	}
	return err
}

func (ps *processSet) Kill() error {
	if !ps.Started() {
		return errors.New("procker: not started")
	}

	var err error
	for _, process := range ps.processes {
		err = process.Kill()
	}
	return err
}

func (ps *processSet) Started() bool {
	return ps.started
}
