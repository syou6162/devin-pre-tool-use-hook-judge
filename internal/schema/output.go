package schema

import (
	"encoding/json"
	"fmt"
	"io"
)

// BlockOutput creates a DevinOutput that blocks the tool execution.
func BlockOutput(reason string) DevinOutput {
	return DevinOutput{
		Decision: DecisionBlock,
		Reason:   reason,
	}
}

// ApproveOutput creates a DevinOutput that approves the tool execution.
func ApproveOutput(reason string) DevinOutput {
	return DevinOutput{
		Decision: DecisionApprove,
		Reason:   reason,
	}
}

// ExitCodeForDecision returns the process exit code for a Devin CLI hook decision.
func ExitCodeForDecision(decision string) int {
	switch decision {
	case DecisionApprove:
		return ExitCodeApprove
	case DecisionBlock, DecisionDeny:
		return ExitCodeBlock
	default:
		return ExitCodeBlock
	}
}

// WriteOutput writes a DevinOutput as JSON to w.
func WriteOutput(w io.Writer, output DevinOutput) error {
	encoded, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("出力のエンコードに失敗しました: %w", err)
	}
	if _, err := w.Write(encoded); err != nil {
		return fmt.Errorf("出力の書き込みに失敗しました: %w", err)
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return fmt.Errorf("出力の書き込みに失敗しました: %w", err)
	}
	return nil
}
