package schema

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestParseDevinInputValid(t *testing.T) {
	t.Parallel()

	input := `{
		"hook_event_name": "PreToolUse",
		"tool_name": "bash",
		"tool_input": {"command": "ls"}
	}`

	devin, err := ParseDevinInput(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseDevinInput() error = %v", err)
	}

	if devin.HookEventName != "PreToolUse" {
		t.Errorf("HookEventName = %q, want %q", devin.HookEventName, "PreToolUse")
	}
	if devin.ToolName != "bash" {
		t.Errorf("ToolName = %q, want %q", devin.ToolName, "bash")
	}
	if devin.ToolInput["command"] != "ls" {
		t.Errorf("ToolInput[command] = %v, want %q", devin.ToolInput["command"], "ls")
	}
}

func TestParseDevinInputErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty input",
			input: "",
			want:  "input is empty",
		},
		{
			name:  "invalid JSON",
			input: "{",
			want:  "parse JSON",
		},
		{
			name:  "missing hook_event_name",
			input: `{"tool_name":"bash","tool_input":{}}`,
			want:  "missing required field: hook_event_name",
		},
		{
			name:  "missing tool_name",
			input: `{"hook_event_name":"PreToolUse","tool_input":{}}`,
			want:  "missing required field: tool_name",
		},
		{
			name:  "missing tool_input",
			input: `{"hook_event_name":"PreToolUse","tool_name":"bash"}`,
			want:  "missing required field: tool_input",
		},
		{
			name:  "tool_input not object",
			input: `{"hook_event_name":"PreToolUse","tool_name":"bash","tool_input":"invalid"}`,
			want:  "invalid tool_input",
		},
		{
			name:  "tool_input null",
			input: `{"hook_event_name":"PreToolUse","tool_name":"bash","tool_input":null}`,
			want:  "invalid tool_input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseDevinInput(strings.NewReader(tt.input))
			if err == nil {
				t.Fatal("ParseDevinInput() expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestBlockOutput(t *testing.T) {
	t.Parallel()

	output := BlockOutput("test reason")
	if output.Decision != DecisionBlock {
		t.Errorf("Decision = %q, want %q", output.Decision, DecisionBlock)
	}
	if output.Reason != "test reason" {
		t.Errorf("Reason = %q, want %q", output.Reason, "test reason")
	}
}

func TestApproveOutput(t *testing.T) {
	t.Parallel()

	output := ApproveOutput()
	if output.Decision != DecisionApprove {
		t.Errorf("Decision = %q, want %q", output.Decision, DecisionApprove)
	}
	if output.Reason != "" {
		t.Errorf("Reason = %q, want empty string", output.Reason)
	}
}

func TestWriteOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteOutput(&buf, ApproveOutput()); err != nil {
		t.Fatalf("WriteOutput() error = %v", err)
	}

	want := `{"decision":"approve"}` + "\n"
	if buf.String() != want {
		t.Errorf("output = %q, want %q", buf.String(), want)
	}
}

func TestWriteOutputBlock(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteOutput(&buf, BlockOutput("invalid input")); err != nil {
		t.Fatalf("WriteOutput() error = %v", err)
	}

	want := `{"decision":"block","reason":"invalid input"}` + "\n"
	if buf.String() != want {
		t.Errorf("output = %q, want %q", buf.String(), want)
	}
}

func TestParseDevinInputExceedsMaxSize(t *testing.T) {
	t.Parallel()

	input := strings.Repeat("a", MaxInputSize+1)
	_, err := ParseDevinInput(strings.NewReader(input))
	if err == nil {
		t.Fatal("ParseDevinInput() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "input exceeds maximum size") {
		t.Errorf("error = %q, want substring %q", err.Error(), "input exceeds maximum size")
	}
}

func TestParseDevinInputAtMaxSize(t *testing.T) {
	t.Parallel()

	base := `{"hook_event_name":"PreToolUse","tool_name":"bash","tool_input":{"data":"`
	suffix := `"}}`
	paddingLen := MaxInputSize - len(base) - len(suffix)
	if paddingLen < 0 {
		t.Fatalf("base input exceeds MaxInputSize")
	}
	input := base + strings.Repeat("x", paddingLen) + suffix

	if len(input) != MaxInputSize {
		t.Fatalf("test input length = %d, want %d", len(input), MaxInputSize)
	}

	_, err := ParseDevinInput(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseDevinInput() error = %v", err)
	}
}

func TestParseDevinInputReadError(t *testing.T) {
	t.Parallel()

	_, err := ParseDevinInput(&errorReader{})
	if err == nil {
		t.Fatal("ParseDevinInput() expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read input") {
		t.Errorf("error = %q, want substring %q", err.Error(), "read input")
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, os.ErrInvalid
}
