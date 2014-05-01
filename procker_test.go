package procker

import (
	"bytes"
	"testing"
)

func TestProcessStart(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("simple", "echo -n procker")
	err := p.Start(env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Wait()
	if err != nil {
		t.Fatal("process failed")
	}

	if stdOut.String() != "procker" {
		t.Error("bad output")
	}

	if stdErr.String() != "" {
		t.Error("bad output")
	}
}

func TestProcessStartUsingEnv(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string = []string{"PROCKER_MSG=hello", "PROCKER_MSG2=world"}

	p := NewProcess("simple", "echo -n $PROCKER_MSG $PROCKER_MSG2")
	err := p.Start(env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Wait()
	if err != nil {
		t.Fatal("process failed")
	}

	if stdOut.String() != "hello world" {
		t.Error("bad output")
	}

	if stdErr.String() != "" {
		t.Error("bad output")
	}
}
