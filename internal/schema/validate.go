package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseDevinInput parses and validates stdin JSON, applying defaults for optional fields.
func ParseDevinInput(raw []byte) (*DevinInput, error) {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	var input DevinInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if err := validateDevinInput(&input); err != nil {
		return nil, err
	}

	applyDevinInputDefaults(&input)
	return &input, nil
}

func validateDevinInput(input *DevinInput) error {
	if input.HookEventName == "" {
		return fmt.Errorf("hook_event_name is required")
	}
	if input.HookEventName != HookEventNamePreToolUse {
		return fmt.Errorf("hook_event_name must be %q", HookEventNamePreToolUse)
	}
	if input.ToolName == "" {
		return fmt.Errorf("tool_name is required")
	}
	if len(input.ToolInput) == 0 {
		return fmt.Errorf("tool_input is required")
	}
	if !json.Valid(input.ToolInput) {
		return fmt.Errorf("tool_input must be valid JSON")
	}
	if input.PermissionMode != "" {
		if _, ok := validPermissionModes[input.PermissionMode]; !ok {
			return fmt.Errorf("invalid permission_mode: %q", input.PermissionMode)
		}
	}
	return nil
}

func applyDevinInputDefaults(input *DevinInput) {
	if input.SessionID == "" {
		input.SessionID = DefaultSessionID
	}
	if input.TranscriptPath == "" {
		input.TranscriptPath = DefaultTranscriptPath
	}
	if input.CWD == "" {
		input.CWD = DefaultCWD
	}
	if input.PermissionMode == "" {
		input.PermissionMode = DefaultPermissionMode
	}
}

// ToJudgeInput converts validated Devin input into the internal judge format.
func ToJudgeInput(input *DevinInput) *JudgeInput {
	return &JudgeInput{
		SessionID:      input.SessionID,
		TranscriptPath: input.TranscriptPath,
		CWD:            input.CWD,
		PermissionMode: input.PermissionMode,
		HookEventName:  input.HookEventName,
		ToolName:       input.ToolName,
		ToolInput:      input.ToolInput,
	}
}

// ValidateJudgeResult validates the judgment engine response.
func ValidateJudgeResult(result *JudgeResult) error {
	if result == nil {
		return fmt.Errorf("judge result is nil")
	}
	switch result.Decision {
	case DecisionApprove, DecisionDeny:
		return nil
	default:
		return fmt.Errorf("invalid decision: %q", result.Decision)
	}
}

// ToDevinOutput maps a validated judge result to Devin hook output.
func ToDevinOutput(result *JudgeResult) DevinOutput {
	switch result.Decision {
	case DecisionApprove:
		return DevinOutput{Decision: DecisionApprove, Reason: result.Reason}
	default:
		decision := DecisionBlock
		if result.Decision == DecisionDeny {
			decision = DecisionDeny
		}
		return DevinOutput{Decision: decision, Reason: result.Reason}
	}
}

// BlockOutput creates a safe-side block response.
func BlockOutput(reason string) DevinOutput {
	return DevinOutput{Decision: DecisionBlock, Reason: reason}
}

// ExitCodeForOutput returns the process exit code for a Devin hook output.
func ExitCodeForOutput(output DevinOutput) int {
	if output.Decision == DecisionApprove {
		return ExitApprove
	}
	return ExitBlock
}
