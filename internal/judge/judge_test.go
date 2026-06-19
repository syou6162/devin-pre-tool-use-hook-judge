package judge

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func TestExtractJSON(t *testing.T) {
	t.Parallel()

	validOutput := `{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "allow",
    "permissionDecisionReason": "safe"
  }
}`

	tests := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "plain json",
			raw:  validOutput,
		},
		{
			name: "json code fence",
			raw: "```json\n" + validOutput + "\n```",
		},
		{
			name: "generic code fence",
			raw: "```\n" + validOutput + "\n```",
		},
		{
			name: "leading and trailing text",
			raw:  "Here is the result:\n" + validOutput + "\nDone.",
		},
		{
			name:    "empty response",
			raw:     "   ",
			wantErr: true,
		},
		{
			name:    "no json object",
			raw:     "not json at all",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractJSON(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("extractJSON() error = %v", err)
			}

			if _, err := schema.ParseHookOutput([]byte(got)); err != nil {
				t.Fatalf("parsed output invalid: %v", err)
			}
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	t.Parallel()

	input := &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "ls"},
	}

	prompt, err := buildPrompt(input, "deny dangerous commands")
	if err != nil {
		t.Fatalf("buildPrompt() error = %v", err)
	}

	for _, want := range []string{
		"<custom_validation_rules>",
		"deny dangerous commands",
		"# Current Tool Usage",
		`"tool_name": "exec"`,
		"<output_json_schema>",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

type mockDevinRunner struct {
	responses []string
	errs      []error
	calls     int
}

func (m *mockDevinRunner) Run(ctx context.Context, promptFile string, model string) (string, error) {
	idx := m.calls
	m.calls++

	if idx < len(m.errs) && m.errs[idx] != nil {
		return "", m.errs[idx]
	}
	if idx < len(m.responses) {
		return m.responses[idx], nil
	}
	return "", errors.New("no more mock responses")
}

func validHookOutputJSON() string {
	return `{
  "hookSpecificOutput": {
    "hookEventName": "PreToolUse",
    "permissionDecision": "allow",
    "permissionDecisionReason": "safe command"
  }
}`
}

func TestDevinEngineJudgeSuccess(t *testing.T) {
	t.Parallel()

	engine := &DevinEngine{
		Runner: &mockDevinRunner{
			responses: []string{validHookOutputJSON()},
		},
	}

	input := &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "ls"},
	}

	output, err := engine.Judge(context.Background(), input, JudgeOptions{
		CustomPrompt: "allow safe commands",
	})
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if output.HookSpecificOutput.PermissionDecision != schema.PermissionAllow {
		t.Fatalf("got decision %q, want allow", output.HookSpecificOutput.PermissionDecision)
	}
}

func TestDevinEngineJudgeRetriesOnInvalidJSON(t *testing.T) {
	t.Parallel()

	runner := &mockDevinRunner{
		responses: []string{
			"not json",
			"still not json",
			validHookOutputJSON(),
		},
	}
	engine := &DevinEngine{Runner: runner}

	input := &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "ls"},
	}

	output, err := engine.Judge(context.Background(), input, JudgeOptions{
		CustomPrompt: "allow safe commands",
	})
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if output.HookSpecificOutput.PermissionDecision != schema.PermissionAllow {
		t.Fatalf("got decision %q, want allow", output.HookSpecificOutput.PermissionDecision)
	}
	if runner.calls != 3 {
		t.Fatalf("calls = %d, want 3", runner.calls)
	}
}

func TestDevinEngineJudgeFailsAfterMaxRetries(t *testing.T) {
	t.Parallel()

	engine := &DevinEngine{
		Runner: &mockDevinRunner{
			responses: []string{"bad", "also bad", "still bad"},
		},
	}

	input := &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "ls"},
	}

	_, err := engine.Judge(context.Background(), input, JudgeOptions{
		CustomPrompt: "allow safe commands",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var judgeErr *JudgeError
	if !errors.As(err, &judgeErr) {
		t.Fatalf("expected JudgeError, got %T", err)
	}
	if judgeErr.Attempts != maxRetries {
		t.Fatalf("attempts = %d, want %d", judgeErr.Attempts, maxRetries)
	}
}

func TestBlockOnError(t *testing.T) {
	t.Parallel()

	output := BlockOnError(&JudgeError{
		Attempts: maxRetries,
		Cause:    errors.New("parse failed"),
	})
	if output.HookSpecificOutput.PermissionDecision != schema.PermissionDeny {
		t.Fatalf("got decision %q, want deny", output.HookSpecificOutput.PermissionDecision)
	}
}

type mockJudgeEngine struct {
	output *schema.HookOutput
	err    error
	called bool
}

func (m *mockJudgeEngine) Judge(ctx context.Context, input *schema.JudgeInput, opts JudgeOptions) (*schema.HookOutput, error) {
	m.called = true
	return m.output, m.err
}

func TestMockJudgeEngine(t *testing.T) {
	t.Parallel()

	mock := &mockJudgeEngine{
		output: &schema.HookOutput{
			HookSpecificOutput: schema.HookSpecificOutput{
				HookEventName:            schema.HookEventName,
				PermissionDecision:       schema.PermissionDeny,
				PermissionDecisionReason: "blocked by mock",
			},
		},
	}

	var engine JudgeEngine = mock
	output, err := engine.Judge(context.Background(), &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "rm -rf /"},
	}, JudgeOptions{CustomPrompt: "test"})
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if !mock.called {
		t.Fatal("mock engine was not called")
	}
	if output.HookSpecificOutput.PermissionDecision != schema.PermissionDeny {
		t.Fatalf("got decision %q, want deny", output.HookSpecificOutput.PermissionDecision)
	}
}

func TestDevinEngineJudgeRetriesOnCodeFence(t *testing.T) {
	t.Parallel()

	engine := &DevinEngine{
		Runner: &mockDevinRunner{
			responses: []string{"```json\n" + validHookOutputJSON() + "\n```"},
		},
	}

	input := &schema.JudgeInput{
		SessionID:      "session-1",
		TranscriptPath: "/tmp/session.jsonl",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  schema.HookEventName,
		ToolName:       "exec",
		ToolInput:      map[string]any{"command": "ls"},
	}

	output, err := engine.Judge(context.Background(), input, JudgeOptions{
		CustomPrompt: "allow safe commands",
	})
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if output.HookSpecificOutput.PermissionDecision != schema.PermissionAllow {
		t.Fatalf("got decision %q, want allow", output.HookSpecificOutput.PermissionDecision)
	}
}
