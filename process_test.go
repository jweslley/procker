package procker

import (
	"bytes"
	"reflect"
	"testing"
	"time"
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

	p := NewProcess("echo -n procker", "", env, stdOut, stdErr)
	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.wait()
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

	p := NewProcess("echo -n $PROCKER_MSG $PROCKER_MSG2",
		"", env, stdOut, stdErr)
	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.wait()
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

	p := NewProcess("cat README.md", "./test", env, stdOut, stdErr)
	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.wait()
	if err != nil {
		t.Fatal("process failed")
	}

	assert(t, "test file!\n", stdOut.String())
	assert(t, "", stdErr.String())
}

func TestProcessCantBeStartedTwiceWhileRunning(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("sleep 1000", "", env, stdOut, stdErr)
	err := p.Start()

	if err != nil {
		t.Fatal("process failed")
	}

	err = p.Start()
	if err == nil {
		t.Fatal("already started")
	}

	p.Stop(0)
}

func TestProcessCanRunMultipleTimes(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("echo -n procker", "", env, stdOut, stdErr)

	i := 0
	for i < 5 {
		if p.Running() {
			t.Fatal("process running")
		}

		err := p.Start()
		if err != nil {
			t.Fatal("process failed")
		}

		err = p.wait()
		if err != nil {
			t.Fatal("process failed")
		}

		assert(t, "procker", stdOut.String())
		assert(t, "", stdErr.String())

		stdOut.Reset()
		stdErr.Reset()
		i += 1
	}
}

func TestProcessWaitOnlyStarted(t *testing.T) {
	p := NewProcess("cat README.md", "", nil, nil, nil)
	err := p.wait()

	if err == nil {
		t.Fatal("not started")
	}
}

func TestProcessStop(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("sh test/lazyecho.sh 5 procker", "", env, stdOut, stdErr)

	err := p.Start()
	if err != nil {
		t.Fatal("process failed")
	}

	go func() {
		erw := p.wait()
		if erw == nil {
			t.Fatalf("not stopped")
		}
	}()

	assert(t, true, p.Running())
	p.Stop(1 * time.Second)
	assert(t, false, p.Running())

	assert(t, "", stdOut.String())
	assert(t, "", stdErr.String())
}

func TestProcessForceStopIfTimeoutExpires(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	var env []string

	p := NewProcess("sh test/trapecho.sh 10 procker", "", env, stdOut, stdErr)

	err := p.Start()
	if err != nil {
		t.Fatal("process failed")
	}

	started := make(chan bool)
	finished := make(chan bool)
	go func() {
		started <- true
		erw := p.wait()
		if erw == nil {
			t.Fatalf("not stopped")
		}
		finished <- true
	}()

	// wait goroutine to start
	<-started

	assert(t, true, p.Running())
	p.Stop(5 * time.Second)
	assert(t, false, p.Running())

	<-finished

	assert(t, "", stdOut.String())
	assert(t, "", stdErr.String())
}
