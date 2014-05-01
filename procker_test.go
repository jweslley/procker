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
	err := p.Start(env, stdOut, stdErr)

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
	err := p.Start(env, stdOut, stdErr)

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
