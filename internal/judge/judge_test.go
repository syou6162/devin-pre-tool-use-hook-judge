package judge

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func TestMockEngine_Judge(t *testing.T) {
	engine := &MockEngine{
		Result: &schema.JudgeResult{Decision: schema.DecisionApprove, Reason: "safe"},
	}
	cfg := &config.Config{Prompt: "test", Model: "default", Timeout: config.DefaultTimeout}
	input := &schema.JudgeInput{
		HookEventName: schema.HookEventNamePreToolUse,
		ToolName:      "exec",
		ToolInput:     json.RawMessage(`{"command":"ls"}`),
	}

	result, err := engine.Judge(context.Background(), input, cfg)
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if result.Decision != schema.DecisionApprove {
		t.Fatalf("decision = %q", result.Decision)
	}
	if engine.Calls != 1 {
		t.Fatalf("calls = %d", engine.Calls)
	}
}

func TestRenderPrompt(t *testing.T) {
	prompt, err := renderPrompt("rule one", `{"tool_name":"exec"}`)
	if err != nil {
		t.Fatalf("renderPrompt() error = %v", err)
	}
	if !strings.Contains(prompt, "rule one") {
		t.Fatal("prompt should include custom rules")
	}
	if !strings.Contains(prompt, `"tool_name":"exec"`) {
		t.Fatal("prompt should include input JSON")
	}
}
