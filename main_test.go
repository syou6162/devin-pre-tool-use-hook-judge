package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

type mockEngine struct {
	output schema.JudgeOutput
	err    error
}

func (m *mockEngine) Judge(_ context.Context, _ schema.JudgeInput, _ config.Config) (schema.JudgeOutput, error) {
	return m.output, m.err
}

func TestRunBlocksWhenConfigMissing(t *testing.T) {
	t.Parallel()

	input := `{"hook_event_name":"PreToolUse","tool_name":"exec","tool_input":{"command":"git push"}}`
	var stdout bytes.Buffer
	exitCode := run(nil, strings.NewReader(input), &stdout, &mockEngine{})

	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Fatalf("stdout = %q, want block decision", stdout.String())
	}
}

func TestRunApprovesWithMockEngine(t *testing.T) {
	t.Parallel()

	input := `{"hook_event_name":"PreToolUse","tool_name":"exec","tool_input":{"command":"git status"}}`
	var stdout bytes.Buffer
	engine := &mockEngine{
		output: schema.JudgeOutput{
			PermissionDecision:       schema.PermissionAllow,
			PermissionDecisionReason: "safe command",
		},
	}

	exitCode := run([]string{"--builtin", "validate_git_push"}, strings.NewReader(input), &stdout, engine)
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), `"decision":"approve"`) {
		t.Fatalf("stdout = %q, want approve decision", stdout.String())
	}
}

func TestRunBlocksWhenEngineFails(t *testing.T) {
	t.Parallel()

	input := `{"hook_event_name":"PreToolUse","tool_name":"exec","tool_input":{"command":"git push origin main"}}`
	var stdout bytes.Buffer
	engine := &mockEngine{err: context.DeadlineExceeded}

	exitCode := run([]string{"--builtin", "validate_git_push"}, strings.NewReader(input), &stdout, engine)
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
	if !strings.Contains(stdout.String(), `"decision":"block"`) {
		t.Fatalf("stdout = %q, want block decision", stdout.String())
	}
}
