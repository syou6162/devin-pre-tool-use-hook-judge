package schema_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func TestParseDevinInput_MinimalValid(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	input, err := schema.ParseDevinInput(raw)
	if err != nil {
		t.Fatalf("ParseDevinInput() error = %v", err)
	}

	if input.HookEventName != schema.HookEventPreToolUse {
		t.Fatalf("HookEventName = %q, want %q", input.HookEventName, schema.HookEventPreToolUse)
	}
	if input.ToolName != "exec" {
		t.Fatalf("ToolName = %q, want exec", input.ToolName)
	}
}

func TestParseDevinInput_FullClaudeCodeCompatible(t *testing.T) {
	raw := []byte(`{
		"session_id": "abc123",
		"transcript_path": "/path/to/session.jsonl",
		"cwd": "/workspace",
		"permission_mode": "default",
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	input, err := schema.ParseDevinInput(raw)
	if err != nil {
		t.Fatalf("ParseDevinInput() error = %v", err)
	}

	judgeInput, err := schema.ToJudgeInput(input)
	if err != nil {
		t.Fatalf("ToJudgeInput() error = %v", err)
	}

	if judgeInput.SessionID != "abc123" {
		t.Fatalf("SessionID = %q, want abc123", judgeInput.SessionID)
	}
	if judgeInput.CWD != "/workspace" {
		t.Fatalf("CWD = %q, want /workspace", judgeInput.CWD)
	}
	if judgeInput.PermissionMode != "default" {
		t.Fatalf("PermissionMode = %q, want default", judgeInput.PermissionMode)
	}
}

func TestParseDevinInput_EmptyInput(t *testing.T) {
	_, err := schema.ParseDevinInput(nil)
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseDevinInput_InvalidJSON(t *testing.T) {
	_, err := schema.ParseDevinInput([]byte(`{invalid`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseDevinInput_MissingHookEventName(t *testing.T) {
	raw := []byte(`{
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	_, err := schema.ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), "hook_event_name") {
		t.Fatalf("expected hook_event_name error, got %v", err)
	}
}

func TestParseDevinInput_InvalidHookEventName(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PostToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	_, err := schema.ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), schema.HookEventPreToolUse) {
		t.Fatalf("expected PreToolUse validation error, got %v", err)
	}
}

func TestParseDevinInput_MissingToolName(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_input": {"command": "ls"}
	}`)

	_, err := schema.ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), "tool_name") {
		t.Fatalf("expected tool_name error, got %v", err)
	}
}

func TestParseDevinInput_MissingToolInput(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec"
	}`)

	_, err := schema.ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), "tool_input") {
		t.Fatalf("expected tool_input error, got %v", err)
	}
}

func TestParseDevinInput_ToolInputMustBeObject(t *testing.T) {
	cases := []string{
		`"not an object"`,
		`[]`,
		`null`,
	}

	for _, toolInput := range cases {
		raw := []byte(`{
			"hook_event_name": "PreToolUse",
			"tool_name": "exec",
			"tool_input": ` + toolInput + `
		}`)

		_, err := schema.ParseDevinInput(raw)
		if err == nil {
			t.Fatalf("expected error for tool_input %s", toolInput)
		}
	}
}

func TestParseDevinInput_InvalidPermissionMode(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"},
		"permission_mode": "unknownMode"
	}`)

	_, err := schema.ParseDevinInput(raw)
	if err == nil {
		t.Fatal("expected error for invalid permission_mode")
	}
}

func TestToJudgeInput_DefaultValues(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	input, err := schema.ParseDevinInput(raw)
	if err != nil {
		t.Fatalf("ParseDevinInput() error = %v", err)
	}

	judgeInput, err := schema.ToJudgeInput(input)
	if err != nil {
		t.Fatalf("ToJudgeInput() error = %v", err)
	}

	if judgeInput.SessionID != "" {
		t.Fatalf("SessionID = %q, want empty string", judgeInput.SessionID)
	}
	if judgeInput.TranscriptPath != "" {
		t.Fatalf("TranscriptPath = %q, want empty string", judgeInput.TranscriptPath)
	}
	if judgeInput.PermissionMode != schema.DefaultPermissionMode {
		t.Fatalf("PermissionMode = %q, want %q", judgeInput.PermissionMode, schema.DefaultPermissionMode)
	}

	expectedCWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	if judgeInput.CWD != expectedCWD {
		t.Fatalf("CWD = %q, want %q", judgeInput.CWD, expectedCWD)
	}
}

func TestToJudgeInput_UsesProvidedCWD(t *testing.T) {
	input := schema.DevinInput{
		HookEventName: schema.HookEventPreToolUse,
		ToolName:      "exec",
		ToolInput:     json.RawMessage(`{"command":"ls"}`),
		CWD:           filepath.Clean("/tmp/project"),
	}

	judgeInput, err := schema.ToJudgeInput(input)
	if err != nil {
		t.Fatalf("ToJudgeInput() error = %v", err)
	}

	if judgeInput.CWD != "/tmp/project" {
		t.Fatalf("CWD = %q, want /tmp/project", judgeInput.CWD)
	}
}

func TestBlockOutput(t *testing.T) {
	output := schema.BlockOutput("blocked for safety")

	if output.Decision != schema.DecisionBlock {
		t.Fatalf("Decision = %q, want %q", output.Decision, schema.DecisionBlock)
	}
	if output.Reason != "blocked for safety" {
		t.Fatalf("Reason = %q, want blocked for safety", output.Reason)
	}
}

func TestApproveOutput(t *testing.T) {
	output := schema.ApproveOutput("looks safe")

	if output.Decision != schema.DecisionApprove {
		t.Fatalf("Decision = %q, want %q", output.Decision, schema.DecisionApprove)
	}
}

func TestExitCodeForDecision(t *testing.T) {
	tests := []struct {
		decision string
		want     int
	}{
		{schema.DecisionApprove, schema.ExitCodeApprove},
		{schema.DecisionBlock, schema.ExitCodeBlock},
		{schema.DecisionDeny, schema.ExitCodeBlock},
		{"unknown", schema.ExitCodeBlock},
	}

	for _, tt := range tests {
		if got := schema.ExitCodeForDecision(tt.decision); got != tt.want {
			t.Fatalf("ExitCodeForDecision(%q) = %d, want %d", tt.decision, got, tt.want)
		}
	}
}

func TestWriteOutput(t *testing.T) {
	var buf strings.Builder
	output := schema.BlockOutput("test reason")

	if err := schema.WriteOutput(&buf, output); err != nil {
		t.Fatalf("WriteOutput() error = %v", err)
	}

	var decoded schema.DevinOutput
	if err := json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &decoded); err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}

	if decoded.Decision != schema.DecisionBlock {
		t.Fatalf("Decision = %q, want %q", decoded.Decision, schema.DecisionBlock)
	}
	if decoded.Reason != "test reason" {
		t.Fatalf("Reason = %q, want test reason", decoded.Reason)
	}
}

func TestValidPermissionModes(t *testing.T) {
	modes := []string{
		"default",
		"plan",
		"acceptEdits",
		"auto",
		"dontAsk",
		"bypassPermissions",
	}

	for _, mode := range modes {
		raw := []byte(`{
			"hook_event_name": "PreToolUse",
			"tool_name": "exec",
			"tool_input": {"command": "ls"},
			"permission_mode": "` + mode + `"
		}`)

		if _, err := schema.ParseDevinInput(raw); err != nil {
			t.Fatalf("permission_mode %q should be valid, got error: %v", mode, err)
		}
	}
}
