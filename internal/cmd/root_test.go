package cmd

import (
	"bytes"
	"testing"
)

func TestRunVersion(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{stdout: &stdout, stderr: &bytes.Buffer{}}

	if got := root.Run([]string{"-version"}); got != ExitSuccess {
		t.Fatalf("Run(-version) = %d, want %d", got, ExitSuccess)
	}

	want := "devin-pre-tool-use-hook-judge 0.1.0\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunHelloWorld(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{stdout: &stdout, stderr: &bytes.Buffer{}}

	if got := root.Run(nil); got != ExitSuccess {
		t.Fatalf("Run(nil) = %d, want %d", got, ExitSuccess)
	}

	want := "hello, world\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}
