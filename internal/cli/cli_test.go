package cli

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestRunVersion(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	err = Run([]string{"-version"})
	if closeErr := w.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := "devin-pre-tool-use-hook-judge 0.1.0\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRunHelloWorld(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	err = Run(nil)
	if closeErr := w.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}

	got := buf.String()
	want := "hello, world\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
