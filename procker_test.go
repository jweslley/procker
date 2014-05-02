package procker

import (
	"bytes"
	"reflect"
	"testing"
)

func assert(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %+v, actual: %+v", expected, actual)
	}
}

func TestProcessStart(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("simple", "echo -n procker")
	err := p.Start("", env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Wait()
	if err != nil {
		t.Fatal("process failed")
	}

	assert(t, "procker", stdOut.String())
	assert(t, "", stdErr.String())
}

func TestProcessStartUsingEnv(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string = []string{"PROCKER_MSG=hello", "PROCKER_MSG2=world"}

	p := NewProcess("simple", "echo -n $PROCKER_MSG $PROCKER_MSG2")
	err := p.Start("", env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Wait()
	if err != nil {
		t.Fatal("process failed")
	}

	assert(t, "hello world", stdOut.String())
	assert(t, "", stdErr.String())
}

func TestProcessStartUsingWithCustomDir(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("cat", "cat README.md")
	err := p.Start("./test", env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Wait()
	if err != nil {
		t.Fatal("process failed")
	}

	assert(t, "test file!\n", stdOut.String())
	assert(t, "", stdErr.String())
}

func TestProcessCantBeStartTwice(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("cat", "cat README.md")
	err := p.Start("./test", env, stdOut, stdErr)

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Start("./test", env, stdOut, stdErr)
	if err == nil {
		t.Fatal("already started")
	}
}

func TestProcessWaitOnlyStarted(t *testing.T) {
	p := NewProcess("cat", "cat README.md")
	err := p.Wait()

	if err == nil {
		t.Fatal("not started")
	}
}
