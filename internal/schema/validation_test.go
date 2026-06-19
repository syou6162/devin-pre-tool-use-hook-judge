package schema

import (
	"testing"
)

func TestValidateDevinInput(t *testing.T) {
	t.Parallel()

	valid := &DevinInput{
		HookEventName: HookEventName,
		ToolName:      "exec",
		ToolInput:     map[string]any{"command": "ls"},
	}
	if err := ValidateDevinInput(valid); err != nil {
		t.Fatalf("ValidateDevinInput() error = %v", err)
	}
}

func TestToJudgeInputFillsDefaults(t *testing.T) {
	t.Parallel()

	input := &DevinInput{
		HookEventName: HookEventName,
		ToolName:      "exec",
		ToolInput:     map[string]any{"command": "ls"},
	}

	judgeInput := ToJudgeInput(input)
	if judgeInput.SessionID != DefaultSessionID {
		t.Fatalf("session_id = %q, want %q", judgeInput.SessionID, DefaultSessionID)
	}
	if judgeInput.PermissionMode != DefaultPermissionMode {
		t.Fatalf("permission_mode = %q, want %q", judgeInput.PermissionMode, DefaultPermissionMode)
	}
	if judgeInput.CWD == "" {
		t.Fatal("cwd should be filled")
	}
}

func TestValidateHookOutput(t *testing.T) {
	t.Parallel()

	output := &HookOutput{
		HookSpecificOutput: HookSpecificOutput{
			HookEventName:            HookEventName,
			PermissionDecision:       PermissionAllow,
			PermissionDecisionReason: "safe",
		},
	}
	if err := ValidateHookOutput(output); err != nil {
		t.Fatalf("ValidateHookOutput() error = %v", err)
	}
}

func TestToDevinOutput(t *testing.T) {
	t.Parallel()

	output := ToDevinOutput(&HookOutput{
		HookSpecificOutput: HookSpecificOutput{
			HookEventName:            HookEventName,
			PermissionDecision:       PermissionAllow,
			PermissionDecisionReason: "safe",
		},
	})
	if output.Decision != "approve" {
		t.Fatalf("decision = %q, want approve", output.Decision)
	}
}
