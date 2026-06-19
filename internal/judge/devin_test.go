package judge

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func TestDevinEngineJudgeWithMock(t *testing.T) {
	t.Parallel()

	engine := &DevinEngine{
		Binary: "devin",
		Runner: func(ctx context.Context, name string, args ...string) *exec.Cmd {
			script := `#!/bin/sh
cat <<'EOF'
{"permissionDecision":"allow","permissionDecisionReason":"safe command"}
EOF
`
			cmd := exec.CommandContext(ctx, "sh", "-c", script)
			return cmd
		},
	}

	input := schema.JudgeInput{
		SessionID:      schema.DefaultSessionID,
		CWD:            ".",
		PermissionMode: schema.DefaultPermissionMode,
		HookEventName:  schema.HookEventPreToolUse,
		ToolName:       "exec",
		ToolInput:      []byte(`{"command":"git status"}`),
	}
	cfg := config.Config{Prompt: "allow safe commands", Model: "sonnet", Timeout: config.DefaultTimeout}

	output, err := engine.Judge(context.Background(), input, cfg)
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if output.PermissionDecision != schema.PermissionAllow {
		t.Fatalf("permissionDecision = %q, want %q", output.PermissionDecision, schema.PermissionAllow)
	}
}

func TestDevinEngineRetriesOnInvalidJSON(t *testing.T) {
	t.Parallel()

	attempts := 0
	engine := &DevinEngine{
		Binary: "devin",
		Runner: func(ctx context.Context, name string, args ...string) *exec.Cmd {
			attempts++
			var script string
			if attempts < 2 {
				script = `echo 'not json'`
			} else {
				script = `echo '{"permissionDecision":"deny","permissionDecisionReason":"blocked"}'`
			}
			return exec.CommandContext(ctx, "sh", "-c", script)
		},
	}

	input := schema.JudgeInput{
		SessionID:      schema.DefaultSessionID,
		CWD:            ".",
		PermissionMode: schema.DefaultPermissionMode,
		HookEventName:  schema.HookEventPreToolUse,
		ToolName:       "exec",
		ToolInput:      []byte(`{"command":"rm -rf /"}`),
	}
	cfg := config.Config{Prompt: "deny destructive commands", Model: "sonnet", Timeout: config.DefaultTimeout}

	output, err := engine.Judge(context.Background(), input, cfg)
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if output.PermissionDecision != schema.PermissionDeny {
		t.Fatalf("permissionDecision = %q, want %q", output.PermissionDecision, schema.PermissionDeny)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestBuildPromptIncludesRulesAndInput(t *testing.T) {
	t.Parallel()

	input := schema.JudgeInput{
		SessionID:      schema.DefaultSessionID,
		CWD:            ".",
		PermissionMode: schema.DefaultPermissionMode,
		HookEventName:  schema.HookEventPreToolUse,
		ToolName:       "exec",
		ToolInput:      []byte(`{"command":"git push origin main"}`),
	}

	prompt, err := buildPrompt(input, "deny main branch pushes")
	if err != nil {
		t.Fatalf("buildPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "deny main branch pushes") {
		t.Fatal("prompt should include validation rules")
	}
	if !strings.Contains(prompt, "git push origin main") {
		t.Fatal("prompt should include tool input")
	}
}

func TestWritePromptFile(t *testing.T) {
	t.Parallel()

	path, err := writePromptFile("test prompt")
	if err != nil {
		t.Fatalf("writePromptFile() error = %v", err)
	}
	defer os.Remove(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "test prompt" {
		t.Fatalf("prompt file = %q, want %q", string(data), "test prompt")
	}
	if filepath.Base(path) == "" {
		t.Fatal("expected non-empty prompt file path")
	}
}

func TestDevinEngineUsesPromptFileFlag(t *testing.T) {
	t.Parallel()

	var capturedArgs []string
	engine := &DevinEngine{
		Binary: "devin",
		Runner: func(ctx context.Context, name string, args ...string) *exec.Cmd {
			capturedArgs = append([]string{}, args...)
			cmd := exec.CommandContext(ctx, "sh", "-c", `echo '{"permissionDecision":"allow","permissionDecisionReason":"ok"}'`)
			cmd.Stdout = io.Discard
			return cmd
		},
	}

	input := schema.JudgeInput{
		SessionID:      schema.DefaultSessionID,
		CWD:            ".",
		PermissionMode: schema.DefaultPermissionMode,
		HookEventName:  schema.HookEventPreToolUse,
		ToolName:       "exec",
		ToolInput:      []byte(`{"command":"echo hi"}`),
	}
	cfg := config.Config{Prompt: "allow echo", Model: "haiku", Timeout: config.DefaultTimeout}

	if _, err := engine.Judge(context.Background(), input, cfg); err != nil {
		t.Fatalf("Judge() error = %v", err)
	}

	foundPromptFile := false
	for i, arg := range capturedArgs {
		if arg == "--prompt-file" && i+1 < len(capturedArgs) && capturedArgs[i+1] != "" {
			foundPromptFile = true
			break
		}
	}
	if !foundPromptFile {
		t.Fatalf("expected --prompt-file in args, got %v", capturedArgs)
	}
}
