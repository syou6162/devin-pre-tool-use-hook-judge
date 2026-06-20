package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{stdin: strings.NewReader(""), stdout: &stdout, stderr: &bytes.Buffer{}}

	if got := root.Run([]string{"-version"}); got != ExitSuccess {
		t.Fatalf("Run(-version) = %d, want %d", got, ExitSuccess)
	}

	want := "devin-pre-tool-use-hook-judge 0.1.0\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunApproveValidInput(t *testing.T) {
	t.Parallel()

	input := `{
		"hook_event_name": "PreToolUse",
		"tool_name": "bash",
		"tool_input": {"command": "ls"}
	}`

	var stdout bytes.Buffer
	root := &Root{
		stdin:  strings.NewReader(input),
		stdout: &stdout,
		stderr: &bytes.Buffer{},
	}

	if got := root.Run(nil); got != ExitSuccess {
		t.Fatalf("Run(nil) = %d, want %d", got, ExitSuccess)
	}

	want := `{"decision":"approve"}` + "\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunBlockInvalidInput(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{
		stdin:  strings.NewReader(`{"tool_name":"bash"}`),
		stdout: &stdout,
		stderr: &bytes.Buffer{},
	}

	if got := root.Run(nil); got != ExitBlock {
		t.Fatalf("Run(nil) = %d, want %d", got, ExitBlock)
	}

	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Errorf("stdout = %q, want block decision", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"reason":`) {
		t.Errorf("stdout = %q, want reason field", stdout.String())
	}
}

func TestRunBlockEmptyInput(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{
		stdin:  strings.NewReader(""),
		stdout: &stdout,
		stderr: &bytes.Buffer{},
	}

	if got := root.Run(nil); got != ExitBlock {
		t.Fatalf("Run(nil) = %d, want %d", got, ExitBlock)
	}

	want := `{"decision":"block","reason":"input is empty"}` + "\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunBlockInvalidJSON(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	root := &Root{
		stdin:  strings.NewReader("{"),
		stdout: &stdout,
		stderr: &bytes.Buffer{},
	}

	if got := root.Run(nil); got != ExitBlock {
		t.Fatalf("Run(nil) = %d, want %d", got, ExitBlock)
	}

	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Errorf("stdout = %q, want block decision", stdout.String())
	}
}
