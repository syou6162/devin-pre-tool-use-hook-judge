package jsonutil

import (
	"strings"
	"testing"
)

func TestExtractJSON_PlainObject(t *testing.T) {
	got, err := ExtractJSON(`{"decision":"approve","reason":"ok"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != `{"decision":"approve","reason":"ok"}` {
		t.Fatalf("got %q", got)
	}
}

func TestExtractJSON_WithMarkdownFence(t *testing.T) {
	input := "```json\n{\"decision\":\"deny\",\"reason\":\"unsafe\"}\n```"
	got, err := ExtractJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, `"decision":"deny"`) {
		t.Fatalf("got %q", got)
	}
}

func TestExtractJSON_WithLeadingText(t *testing.T) {
	input := "Here is the result:\n\n```json\n{\"decision\":\"approve\",\"reason\":\"safe\"}\n```\n"
	got, err := ExtractJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, `"decision":"approve"`) {
		t.Fatalf("got %q", got)
	}
}

func TestExtractJSON_WithNestedBracesInString(t *testing.T) {
	input := `{"decision":"deny","reason":"contains { braces }"}`
	got, err := ExtractJSON(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != input {
		t.Fatalf("got %q", got)
	}
}

func TestExtractJSON_Array(t *testing.T) {
	got, err := ExtractJSON(`[{"a":1},{"b":2}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(got, "[") {
		t.Fatalf("got %q", got)
	}
}

func TestExtractJSON_EmptyInput(t *testing.T) {
	_, err := ExtractJSON("   ")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestExtractJSON_NoJSON(t *testing.T) {
	_, err := ExtractJSON("no json here")
	if err == nil {
		t.Fatal("expected error when JSON is missing")
	}
}

func TestExtractJSON_UnbalancedJSON(t *testing.T) {
	_, err := ExtractJSON(`{"decision":"approve"`)
	if err == nil {
		t.Fatal("expected error for unbalanced JSON")
	}
}
