package main_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
	mainpkg "github.com/syou6162/devin-pre-tool-use-hook-judge"
)

func TestRun_ValidInputApproves(t *testing.T) {
	input := `{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`

	var stdout bytes.Buffer
	code := mainpkg.Run(strings.NewReader(input), &stdout)
	if code != schema.ExitCodeApprove {
		t.Fatalf("exit code = %d, want %d", code, schema.ExitCodeApprove)
	}

	if !strings.Contains(stdout.String(), `"decision":"approve"`) {
		t.Fatalf("stdout = %q, want approve decision", stdout.String())
	}
}

func TestRun_InvalidInputBlocks(t *testing.T) {
	var stdout bytes.Buffer
	code := mainpkg.Run(strings.NewReader(`{invalid`), &stdout)
	if code != schema.ExitCodeBlock {
		t.Fatalf("exit code = %d, want %d", code, schema.ExitCodeBlock)
	}

	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Fatalf("stdout = %q, want block decision", stdout.String())
	}
}

func TestRun_MissingRequiredFieldBlocks(t *testing.T) {
	input := `{
		"hook_event_name": "PreToolUse",
		"tool_input": {"command": "ls"}
	}`

	var stdout bytes.Buffer
	code := mainpkg.Run(strings.NewReader(input), &stdout)
	if code != schema.ExitCodeBlock {
		t.Fatalf("exit code = %d, want %d", code, schema.ExitCodeBlock)
	}

	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Fatalf("stdout = %q, want block decision", stdout.String())
	}
}
