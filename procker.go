package procker

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Process struct {
	Name    string
	Command string
	Dir     string
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
	cmd     *exec.Cmd
}

func NewProcess(name, command string) *Process {
	return &Process{Name: name, Command: command}
}

func (p *Process) Start() error {
	if p.Started() {
		return errors.New("procker: already started")
	}

	args := strings.Fields(p.expandedCmd(p.Env))
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = p.Env
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	return p.cmd.Start()
}

func (p *Process) Wait() error {
	if !p.Started() {
		return errors.New("procker: not started")
	}
	return p.cmd.Wait()
}

func (p *Process) Kill() error {
	if !p.Started() {
		return errors.New("procker: not started")
	}
	return p.cmd.Process.Kill()
}

func (p *Process) Pid() int {
	if !p.Started() {
		return 0
	}
	return p.cmd.Process.Pid
}

func (p *Process) Started() bool {
	return p.cmd != nil
}

func (p *Process) expandedCmd(env []string) string {
	m := env2Map(env)
	return os.Expand(p.Command, func(name string) string {
		return m[name]
	})
}
