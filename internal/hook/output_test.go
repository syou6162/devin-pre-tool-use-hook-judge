package hook

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteBlock(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteBlock(&buf, "test reason"); err != nil {
		t.Fatalf("WriteBlock() error = %v", err)
	}

	var output DevinOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if output.Decision != DecisionBlock {
		t.Errorf("Decision = %q, want %q", output.Decision, DecisionBlock)
	}
	if output.Reason != "test reason" {
		t.Errorf("Reason = %q, want %q", output.Reason, "test reason")
	}
}
