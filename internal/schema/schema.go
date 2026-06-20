package schema

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	DecisionApprove = "approve"
	DecisionBlock   = "block"

	// MaxInputSize is the maximum allowed stdin payload size in bytes.
	MaxInputSize = 1 << 20 // 1 MiB
)

// DevinInput is the JSON payload received from Devin CLI via stdin.
type DevinInput struct {
	HookEventName string                 `json:"hook_event_name"`
	ToolName      string                 `json:"tool_name"`
	ToolInput     map[string]interface{} `json:"tool_input"`
}

// DevinOutput is the JSON payload written to stdout for Devin CLI.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

// ParseDevinInput reads JSON from r, validates the schema, and returns DevinInput.
func ParseDevinInput(r io.Reader) (DevinInput, error) {
	limited := io.LimitReader(r, MaxInputSize+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return DevinInput{}, fmt.Errorf("read input: %w", err)
	}
	if len(data) > MaxInputSize {
		return DevinInput{}, fmt.Errorf("input exceeds maximum size of %d bytes", MaxInputSize)
	}

	if len(data) == 0 {
		return DevinInput{}, fmt.Errorf("input is empty")
	}

	var raw struct {
		HookEventName string          `json:"hook_event_name"`
		ToolName      string          `json:"tool_name"`
		ToolInput     json.RawMessage `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return DevinInput{}, fmt.Errorf("parse JSON: %w", err)
	}

	input := DevinInput{
		HookEventName: raw.HookEventName,
		ToolName:      raw.ToolName,
	}

	if err := validateDevinInput(&input, raw.ToolInput); err != nil {
		return DevinInput{}, err
	}

	return input, nil
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
