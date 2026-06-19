package schema

import (
	"encoding/json"
	"fmt"
)

func BlockOutput(reason string) DevinOutput {
	return DevinOutput{
		Decision: DecisionBlock,
		Reason:   reason,
	}
}

func ApproveOutput(reason string) DevinOutput {
	output := DevinOutput{Decision: DecisionApprove}
	if reason != "" {
		output.Reason = reason
	}
	return output
}

func ParseDevinInput(raw []byte) (DevinInput, error) {
	var input DevinInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return DevinInput{}, fmt.Errorf("invalid JSON input: %w", err)
	}
	return input, nil
}

func ValidateDevinInput(input DevinInput) error {
	if input.HookEventName != HookEventPreToolUse {
		return fmt.Errorf("unsupported hook event %q", input.HookEventName)
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
	return nil
}

func ToJudgeInput(input DevinInput) JudgeInput {
	toolInput := input.ToolInput
	if len(toolInput) == 0 {
		toolInput = json.RawMessage("{}")
	}

	return JudgeInput{
		SessionID:      DefaultSessionID,
		TranscriptPath: "",
		CWD:            DefaultCWD,
		PermissionMode: DefaultPermissionMode,
		HookEventName:  input.HookEventName,
		ToolName:       input.ToolName,
		ToolInput:      toolInput,
	}
}

func ToDevinOutput(judge JudgeOutput) DevinOutput {
	switch judge.PermissionDecision {
	case PermissionAllow:
		return ApproveOutput(judge.PermissionDecisionReason)
	case PermissionDeny, PermissionAsk:
		return BlockOutput(judge.PermissionDecisionReason)
	default:
		return BlockOutput(DefaultPermissionReason)
	}
}

func ValidateJudgeOutput(output JudgeOutput) error {
	switch output.PermissionDecision {
	case PermissionAllow, PermissionDeny, PermissionAsk:
		if output.PermissionDecisionReason == "" {
			return fmt.Errorf("permissionDecisionReason is required")
		}
		return nil
	default:
		return fmt.Errorf("invalid permissionDecision %q", output.PermissionDecision)
	}
}

func ExitCodeForOutput(output DevinOutput) int {
	switch output.Decision {
	case DecisionApprove:
		return 0
	case DecisionBlock, DecisionDeny:
		return 2
	default:
		return 2
	}
}
