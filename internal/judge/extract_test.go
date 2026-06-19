package judge

import "testing"

func TestExtractJSONFromCodeFence(t *testing.T) {
	t.Parallel()

	raw := "Here is the result:\n```json\n{\"permissionDecision\":\"allow\",\"permissionDecisionReason\":\"ok\"}\n```"
	got, err := ExtractJSON(raw)
	if err != nil {
		t.Fatalf("ExtractJSON() error = %v", err)
	}
	want := `{"permissionDecision":"allow","permissionDecisionReason":"ok"}`
	if got != want {
		t.Fatalf("ExtractJSON() = %q, want %q", got, want)
	}
}

func TestExtractJSONFromPlainObject(t *testing.T) {
	t.Parallel()

	raw := `Some text before {"permissionDecision":"deny","permissionDecisionReason":"blocked"} trailing text`
	got, err := ExtractJSON(raw)
	if err != nil {
		t.Fatalf("ExtractJSON() error = %v", err)
	}
	want := `{"permissionDecision":"deny","permissionDecisionReason":"blocked"}`
	if got != want {
		t.Fatalf("ExtractJSON() = %q, want %q", got, want)
	}
}

func TestExtractJSONEmpty(t *testing.T) {
	t.Parallel()

	if _, err := ExtractJSON("   "); err != ErrEmptyResponse {
		t.Fatalf("expected ErrEmptyResponse, got %v", err)
	}
}
