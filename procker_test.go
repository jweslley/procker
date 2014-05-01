package procker

import (
	"bytes"
	"testing"
)

func TestProcessStart(t *testing.T) {
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}

	p := NewProcess("simple", "echo -n procker")
	err := p.Start(stdOut, stdErr)

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
