package procker

import (
	"io"
	"os"
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

func (p *Process) Start(env []string, out, err io.Writer) error {
	args := strings.Fields(p.expandedCmd(env))
	p.cmd = exec.Command(args[0], args[1:]...)
	p.cmd.Env = env
	p.cmd.Stdout = out
	p.cmd.Stderr = err
	return p.cmd.Start()
}

func (p *Process) Wait() error {
	return p.cmd.Wait()
}

func (p *Process) expandedCmd(env []string) string {
	m := env2Map(env)
	return os.Expand(p.Command, func(name string) string {
		return m[name]
	})
}

func env2Map(env []string) map[string]string {
	m := make(map[string]string)
	for _, value := range env {
		pair := strings.SplitN(value, "=", 2)
		if len(pair) == 2 {
			m[pair[0]] = pair[1]
		}
	}
	return m
}
