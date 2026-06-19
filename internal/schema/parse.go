package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseDevinInput parses and validates raw JSON into a DevinInput.
func ParseDevinInput(raw []byte) (DevinInput, error) {
	if len(raw) == 0 {
		return DevinInput{}, fmt.Errorf("入力が空です")
	}

	var input DevinInput
	if err := json.Unmarshal(raw, &input); err != nil {
		return DevinInput{}, fmt.Errorf("JSONパースエラー: %w", err)
	}

	if err := validateDevinInput(input); err != nil {
		return DevinInput{}, err
	}

	return input, nil
}

func validateDevinInput(input DevinInput) error {
	if input.HookEventName == "" {
		return fmt.Errorf("必須フィールド hook_event_name が欠けています")
	}
	if input.HookEventName != HookEventPreToolUse {
		return fmt.Errorf("hook_event_name は %q である必要があります", HookEventPreToolUse)
	}
	if input.ToolName == "" {
		return fmt.Errorf("必須フィールド tool_name が欠けています")
	}
	if len(input.ToolInput) == 0 {
		return fmt.Errorf("必須フィールド tool_input が欠けています")
	}
	if !json.Valid(input.ToolInput) {
		return fmt.Errorf("tool_input は有効な JSON である必要があります")
	}

	var toolInput map[string]interface{}
	if err := json.Unmarshal(input.ToolInput, &toolInput); err != nil {
		return fmt.Errorf("tool_input はオブジェクトである必要があります")
	}
	if toolInput == nil {
		return fmt.Errorf("tool_input はオブジェクトである必要があります")
	}

	if input.PermissionMode != "" {
		if _, ok := validPermissionModes[input.PermissionMode]; !ok {
			return fmt.Errorf("permission_mode の値が不正です: %q", input.PermissionMode)
		}
	}

	return nil
}

// ToJudgeInput converts a validated DevinInput into the internal JudgeInput format,
// filling missing fields with safe defaults.
func ToJudgeInput(input DevinInput) (JudgeInput, error) {
	var toolInput map[string]interface{}
	if err := json.Unmarshal(input.ToolInput, &toolInput); err != nil {
		return JudgeInput{}, fmt.Errorf("tool_input の変換に失敗しました: %w", err)
	}

	cwd := input.CWD
	if cwd == "" {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		} else {
			cwd = "."
		}
	}

	permissionMode := input.PermissionMode
	if permissionMode == "" {
		permissionMode = DefaultPermissionMode
	}

	return JudgeInput{
		SessionID:      input.SessionID,
		TranscriptPath: input.TranscriptPath,
		CWD:            cwd,
		PermissionMode: permissionMode,
		HookEventName:  input.HookEventName,
		ToolName:       input.ToolName,
		ToolInput:      toolInput,
	}, nil
}
