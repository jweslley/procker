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

	p := NewProcess("simple", "echo -n procker", "", env, stdOut, stdErr)
	err := p.Start()

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
	env := []string{"PROCKER_MSG=hello", "PROCKER_MSG2=world"}

	p := NewProcess("simple",
		"echo -n $PROCKER_MSG $PROCKER_MSG2",
		"", env, stdOut, stdErr)
	err := p.Start()

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

	p := NewProcess("cat", "cat README.md", "./test", env, stdOut, stdErr)
	err := p.Start()

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

	p := NewProcess("cat", "cat README.md", "./test", env, stdOut, stdErr)
	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Start()
	if err == nil {
		t.Fatal("already started")
	}
}

func TestProcessWaitOnlyStarted(t *testing.T) {
	p := NewProcess("cat", "cat README.md", "", nil, nil, nil)
	err := p.Wait()

	if err == nil {
		t.Fatal("not started")
	}
}

func TestProcessKill(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("lazyecho",
		"sh test/lazyecho.sh 5 procker",
		"", env, stdOut, stdErr)

	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	go func() {
		erw := p.Wait()
		if erw == nil {
			t.Fatalf("not killed")
		}
	}()

	err = p.Kill()
	if err != nil {
		t.Fatal("not killed")
	}

	assert(t, "", stdOut.String())
	assert(t, "", stdErr.String())
}
