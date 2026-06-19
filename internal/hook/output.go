package hook

import (
	"encoding/json"
	"io"
)

const (
	DecisionApprove = "approve"
	DecisionBlock   = "block"
)

// DevinOutput is the JSON response format for Devin CLI command hooks.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

// WriteBlock writes a block decision to w and returns exit code 2.
func WriteBlock(w io.Writer, reason string) error {
	return writeOutput(w, DevinOutput{
		Decision: DecisionBlock,
		Reason:   reason,
	})
}

// WriteApprove writes an approve decision to w.
func WriteApprove(w io.Writer, reason string) error {
	return writeOutput(w, DevinOutput{
		Decision: DecisionApprove,
		Reason:   reason,
	})
}

func writeOutput(w io.Writer, output DevinOutput) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
