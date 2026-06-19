package schema

import (
	"encoding/json"
	"testing"
)

func TestValidateDevinInput(t *testing.T) {
	t.Parallel()

	valid := DevinInput{
		HookEventName: HookEventPreToolUse,
		ToolName:      "exec",
		ToolInput:     json.RawMessage(`{"command":"git status"}`),
	}
	if err := ValidateDevinInput(valid); err != nil {
		t.Fatalf("expected valid input, got %v", err)
	}
}

func TestValidateDevinInputRejectsMissingToolName(t *testing.T) {
	t.Parallel()

	input := DevinInput{
		HookEventName: HookEventPreToolUse,
		ToolInput:     json.RawMessage(`{"command":"git status"}`),
	}
	if err := ValidateDevinInput(input); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestToDevinOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		judge    JudgeOutput
		decision string
		exitCode int
	}{
		{
			name: "allow maps to approve",
			judge: JudgeOutput{
				PermissionDecision:       PermissionAllow,
				PermissionDecisionReason: "safe command",
			},
			decision: DecisionApprove,
			exitCode: 0,
		},
		{
			name: "deny maps to block",
			judge: JudgeOutput{
				PermissionDecision:       PermissionDeny,
				PermissionDecisionReason: "unsafe command",
			},
			decision: DecisionBlock,
			exitCode: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			output := ToDevinOutput(tt.judge)
			if output.Decision != tt.decision {
				t.Fatalf("decision = %q, want %q", output.Decision, tt.decision)
			}
			if got := ExitCodeForOutput(output); got != tt.exitCode {
				t.Fatalf("exit code = %d, want %d", got, tt.exitCode)
			}
		})
	}
}
