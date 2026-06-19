package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

var validPermissionModes = map[string]struct{}{
	"default":            {},
	"plan":               {},
	"acceptEdits":        {},
	"auto":               {},
	"dontAsk":            {},
	"bypassPermissions":  {},
}

var validPermissionDecisions = map[string]struct{}{
	PermissionAllow: {},
	PermissionDeny:  {},
	PermissionAsk:   {},
	PermissionDefer: {},
}

// ToJudgeInput converts DevinInput to JudgeInput, filling missing fields with defaults.
func ToJudgeInput(input *DevinInput) *JudgeInput {
	if input == nil {
		return nil
	}

	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = DefaultSessionID
	}

	transcriptPath := input.TranscriptPath
	if transcriptPath == "" {
		transcriptPath = DefaultTranscriptPath
	}

	cwd := input.CWD
	if cwd == "" {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		} else {
			cwd = DefaultCWD
		}
	}

	permissionMode := input.PermissionMode
	if permissionMode == "" {
		permissionMode = DefaultPermissionMode
	}

	toolInput := input.ToolInput
	if toolInput == nil {
		toolInput = map[string]any{}
	}

	return &JudgeInput{
		SessionID:      sessionID,
		TranscriptPath: transcriptPath,
		CWD:            cwd,
		PermissionMode: permissionMode,
		HookEventName:  input.HookEventName,
		ToolName:       input.ToolName,
		ToolInput:      toolInput,
	}
}

// ValidateDevinInput validates the raw Devin CLI hook input.
func ValidateDevinInput(input *DevinInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	}
	if input.HookEventName != HookEventName {
		return fmt.Errorf("hook_event_name must be %q", HookEventName)
	}
	if input.ToolName == "" {
		return fmt.Errorf("tool_name is required")
	}
	if input.ToolInput == nil {
		return fmt.Errorf("tool_input is required")
	}
	if input.PermissionMode != "" {
		if _, ok := validPermissionModes[input.PermissionMode]; !ok {
			return fmt.Errorf("invalid permission_mode: %q", input.PermissionMode)
		}
	}
	return nil
}

// ValidateJudgeInput validates the normalized judge input.
func ValidateJudgeInput(input *JudgeInput) error {
	if input == nil {
		return fmt.Errorf("input is nil")
	}
	if input.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if input.TranscriptPath == "" && input.SessionID == DefaultSessionID {
		// transcript_path may be empty when session_id is defaulted
	}
	if input.CWD == "" {
		return fmt.Errorf("cwd is required")
	}
	if input.HookEventName != HookEventName {
		return fmt.Errorf("hook_event_name must be %q", HookEventName)
	}
	if input.ToolName == "" {
		return fmt.Errorf("tool_name is required")
	}
	if input.ToolInput == nil {
		return fmt.Errorf("tool_input is required")
	}
	if _, ok := validPermissionModes[input.PermissionMode]; !ok {
		return fmt.Errorf("invalid permission_mode: %q", input.PermissionMode)
	}
	return nil
}

// ValidateHookOutput validates the PreToolUse hook output from the judgment engine.
func ValidateHookOutput(output *HookOutput) error {
	if output == nil {
		return fmt.Errorf("output is nil")
	}

	specific := output.HookSpecificOutput
	if specific.HookEventName != HookEventName {
		return fmt.Errorf("hookEventName must be %q", HookEventName)
	}
	if specific.PermissionDecision == "" {
		return fmt.Errorf("permissionDecision is required")
	}
	if _, ok := validPermissionDecisions[specific.PermissionDecision]; !ok {
		return fmt.Errorf("invalid permissionDecision: %q", specific.PermissionDecision)
	}
	if specific.PermissionDecisionReason == "" {
		return fmt.Errorf("permissionDecisionReason is required")
	}
	return nil
}

// ParseHookOutput unmarshals and validates hook output JSON.
func ParseHookOutput(data []byte) (*HookOutput, error) {
	var output HookOutput
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if err := ValidateHookOutput(&output); err != nil {
		return nil, err
	}
	return &output, nil
}

// BlockOutput creates a deny/block response in Devin CLI format.
func BlockOutput(reason string) *DevinOutput {
	return &DevinOutput{
		Decision: "block",
		Reason:   reason,
	}
}

// BlockHookOutput creates a deny response in PreToolUse hook format.
func BlockHookOutput(reason string) *HookOutput {
	return &HookOutput{
		HookSpecificOutput: HookSpecificOutput{
			HookEventName:            HookEventName,
			PermissionDecision:       PermissionDeny,
			PermissionDecisionReason: reason,
		},
	}
}

// ToDevinOutput converts a hook output to Devin CLI format.
func ToDevinOutput(output *HookOutput) *DevinOutput {
	if output == nil {
		return BlockOutput("判定結果がありません")
	}

	decision := "block"
	switch output.HookSpecificOutput.PermissionDecision {
	case PermissionAllow:
		decision = "approve"
	case PermissionDeny:
		decision = "block"
	case PermissionAsk:
		decision = "block"
	case PermissionDefer:
		decision = "approve"
	}

	return &DevinOutput{
		Decision: decision,
		Reason:   output.HookSpecificOutput.PermissionDecisionReason,
	}
}
