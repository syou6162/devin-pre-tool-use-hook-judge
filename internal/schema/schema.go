package schema

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	DecisionApprove = "approve"
	DecisionBlock   = "block"

	defaultPermissionMode = "default"
)

// DevinInput is the JSON payload received from Devin CLI via stdin.
type DevinInput struct {
	HookEventName string                 `json:"hook_event_name"`
	ToolName      string                 `json:"tool_name"`
	ToolInput     map[string]interface{} `json:"tool_input"`
}

// JudgeInput is the internal representation used for hook validation.
type JudgeInput struct {
	SessionID      string                 `json:"session_id"`
	TranscriptPath string                 `json:"transcript_path"`
	Cwd            string                 `json:"cwd"`
	PermissionMode string                 `json:"permission_mode"`
	HookEventName  string                 `json:"hook_event_name"`
	ToolName       string                 `json:"tool_name"`
	ToolInput      map[string]interface{} `json:"tool_input"`
}

// DevinOutput is the JSON payload written to stdout for Devin CLI.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

type rawInput struct {
	HookEventName  string          `json:"hook_event_name"`
	ToolName       string          `json:"tool_name"`
	ToolInput      json.RawMessage `json:"tool_input"`
	SessionID      string          `json:"session_id"`
	TranscriptPath string          `json:"transcript_path"`
	Cwd            string          `json:"cwd"`
	PermissionMode string          `json:"permission_mode"`
}

// ParseDevinInput reads JSON from r, validates the schema, and converts it to JudgeInput.
func ParseDevinInput(r io.Reader) (JudgeInput, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return JudgeInput{}, fmt.Errorf("read input: %w", err)
	}

	if len(data) == 0 {
		return JudgeInput{}, fmt.Errorf("input is empty")
	}

	var raw rawInput
	if err := json.Unmarshal(data, &raw); err != nil {
		return JudgeInput{}, fmt.Errorf("parse JSON: %w", err)
	}

	devin := DevinInput{
		HookEventName: raw.HookEventName,
		ToolName:      raw.ToolName,
	}

	if err := validateDevinInput(&devin, raw.ToolInput); err != nil {
		return JudgeInput{}, err
	}

	return toJudgeInput(devin, raw), nil
}

func validateDevinInput(input *DevinInput, toolInputRaw json.RawMessage) error {
	if input.HookEventName == "" {
		return fmt.Errorf("missing required field: hook_event_name")
	}
	if input.ToolName == "" {
		return fmt.Errorf("missing required field: tool_name")
	}
	if len(toolInputRaw) == 0 {
		return fmt.Errorf("missing required field: tool_input")
	}

	var toolInput map[string]interface{}
	if err := json.Unmarshal(toolInputRaw, &toolInput); err != nil {
		return fmt.Errorf("invalid tool_input: must be a JSON object")
	}
	if toolInput == nil {
		return fmt.Errorf("invalid tool_input: must be a JSON object")
	}

	input.ToolInput = toolInput
	return nil
}

func toJudgeInput(devin DevinInput, raw rawInput) JudgeInput {
	sessionID := raw.SessionID
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	cwd := raw.Cwd
	if cwd == "" {
		cwd = os.Getenv("PWD")
	}

	permissionMode := raw.PermissionMode
	if permissionMode == "" {
		permissionMode = defaultPermissionMode
	}

	return JudgeInput{
		SessionID:      sessionID,
		TranscriptPath: raw.TranscriptPath,
		Cwd:            cwd,
		PermissionMode: permissionMode,
		HookEventName:  devin.HookEventName,
		ToolName:       devin.ToolName,
		ToolInput:      devin.ToolInput,
	}
}

func generateSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// BlockOutput returns a DevinOutput that blocks execution with the given reason.
func BlockOutput(reason string) DevinOutput {
	return DevinOutput{
		Decision: DecisionBlock,
		Reason:   reason,
	}
}

// ApproveOutput returns a DevinOutput that approves execution.
func ApproveOutput() DevinOutput {
	return DevinOutput{
		Decision: DecisionApprove,
	}
}

// WriteOutput encodes output as JSON to w.
func WriteOutput(w io.Writer, output DevinOutput) error {
	data, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}
