package procker

import (
	"io"
	"os/exec"
	"strings"
)

type Process struct {
	Name    string
	Command string
	cmd     *exec.Cmd
}

func NewProcess(name, command string) *Process {
	return &Process{Name: name, Command: command}
}

func (p *Process) Start(out, err io.Writer) error {
	args := strings.Fields(p.Command)
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Stdout = out
	p.cmd.Stderr = err
	return p.cmd.Start()
}

func (p *Process) Wait() error {
	return p.cmd.Wait()
}
