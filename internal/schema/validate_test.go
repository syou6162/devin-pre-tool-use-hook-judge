package schema

import (
	"strings"
	"testing"
)

func TestParseDevinInput_ValidMinimal(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)

	input, err := ParseDevinInput(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.SessionID != DefaultSessionID {
		t.Fatalf("session_id = %q, want %q", input.SessionID, DefaultSessionID)
	}
	if input.CWD != DefaultCWD {
		t.Fatalf("cwd = %q, want %q", input.CWD, DefaultCWD)
	}
	if input.PermissionMode != DefaultPermissionMode {
		t.Fatalf("permission_mode = %q, want %q", input.PermissionMode, DefaultPermissionMode)
	}
}

func TestParseDevinInput_ValidFull(t *testing.T) {
	raw := []byte(`{
		"session_id": "abc123",
		"transcript_path": "/tmp/session.jsonl",
		"cwd": "/workspace",
		"permission_mode": "default",
		"hook_event_name": "PreToolUse",
		"tool_name": "edit",
		"tool_input": {"file_path": "/tmp/a.txt", "content": "hello"}
	}`)

	input, err := ParseDevinInput(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.SessionID != "abc123" {
		t.Fatalf("session_id = %q", input.SessionID)
	}
}

func TestParseDevinInput_MissingHookEventName(t *testing.T) {
	raw := []byte(`{"tool_name":"exec","tool_input":{"command":"ls"}}`)
	_, err := ParseDevinInput(raw)
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
	_, err := ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), "PreToolUse") {
		t.Fatalf("expected PreToolUse error, got %v", err)
	}
}

func TestParseDevinInput_MissingToolName(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_input": {"command": "ls"}
	}`)
	_, err := ParseDevinInput(raw)
	if err == nil || !strings.Contains(err.Error(), "tool_name") {
		t.Fatalf("expected tool_name error, got %v", err)
	}
}

func TestParseDevinInput_InvalidJSON(t *testing.T) {
	_, err := ParseDevinInput([]byte(`{invalid`))
	if err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestParseDevinInput_InvalidPermissionMode(t *testing.T) {
	raw := []byte(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"},
		"permission_mode": "unknownMode"
	}`)
	_, err := ParseDevinInput(raw)
	if err == nil {
		t.Fatal("expected permission_mode error")
	}
}

func TestValidateJudgeResult(t *testing.T) {
	tests := []struct {
		name    string
		result  JudgeResult
		wantErr bool
	}{
		{name: "approve", result: JudgeResult{Decision: DecisionApprove, Reason: "ok"}, wantErr: false},
		{name: "deny", result: JudgeResult{Decision: DecisionDeny, Reason: "no"}, wantErr: false},
		{name: "invalid", result: JudgeResult{Decision: "maybe", Reason: "x"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJudgeResult(&tt.result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJudgeResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToDevinOutput(t *testing.T) {
	approve := ToDevinOutput(&JudgeResult{Decision: DecisionApprove, Reason: "safe"})
	if approve.Decision != DecisionApprove || approve.Reason != "safe" {
		t.Fatalf("approve mapping failed: %+v", approve)
	}

	deny := ToDevinOutput(&JudgeResult{Decision: DecisionDeny, Reason: "unsafe"})
	if deny.Decision != DecisionDeny || deny.Reason != "unsafe" {
		t.Fatalf("deny mapping failed: %+v", deny)
	}
}

func TestExitCodeForOutput(t *testing.T) {
	if ExitCodeForOutput(DevinOutput{Decision: DecisionApprove}) != ExitApprove {
		t.Fatal("approve should exit 0")
	}
	if ExitCodeForOutput(DevinOutput{Decision: DecisionBlock}) != ExitBlock {
		t.Fatal("block should exit 2")
	}
}

func TestToJudgeInput(t *testing.T) {
	input := &DevinInput{
		SessionID:      "s1",
		TranscriptPath: "/tmp/t",
		CWD:            "/workspace",
		PermissionMode: "default",
		HookEventName:  HookEventNamePreToolUse,
		ToolName:       "exec",
		ToolInput:      []byte(`{"command":"ls"}`),
	}

	judgeInput := ToJudgeInput(input)
	if judgeInput.ToolName != "exec" {
		t.Fatalf("tool_name = %q", judgeInput.ToolName)
	}
}
