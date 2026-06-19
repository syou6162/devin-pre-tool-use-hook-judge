package app

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/judge"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func TestRun_IntegrationWithMockApprove(t *testing.T) {
	stdin := strings.NewReader(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)
	var stdout bytes.Buffer

	engine := &judge.MockEngine{
		Result: &schema.JudgeResult{
			Decision: schema.DecisionApprove,
			Reason:   "参照系コマンドのため安全です",
		},
	}

	exitCode := Run(context.Background(), stdin, &stdout, Options{BuiltinName: "validate_find"}, engine)
	if exitCode != schema.ExitApprove {
		t.Fatalf("exit code = %d, want %d", exitCode, schema.ExitApprove)
	}

	var output schema.DevinOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.Decision != schema.DecisionApprove {
		t.Fatalf("decision = %q", output.Decision)
	}
	if engine.Calls != 1 {
		t.Fatalf("engine calls = %d", engine.Calls)
	}
}

func TestRun_IntegrationWithMockDeny(t *testing.T) {
	stdin := strings.NewReader(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "rm -rf /"}
	}`)
	var stdout bytes.Buffer

	engine := &judge.MockEngine{
		Result: &schema.JudgeResult{
			Decision: schema.DecisionDeny,
			Reason:   "破壊的操作のため拒否します",
		},
	}

	exitCode := Run(context.Background(), stdin, &stdout, Options{BuiltinName: "validate_xargs"}, engine)
	if exitCode != schema.ExitBlock {
		t.Fatalf("exit code = %d, want %d", exitCode, schema.ExitBlock)
	}

	var output schema.DevinOutput
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("decode output: %v", err)
	}
	if output.Decision != schema.DecisionDeny {
		t.Fatalf("decision = %q", output.Decision)
	}
}

func TestRun_BlocksOnInvalidInput(t *testing.T) {
	stdin := strings.NewReader(`{"tool_name":"exec"}`)
	var stdout bytes.Buffer
	engine := &judge.MockEngine{Result: &schema.JudgeResult{Decision: schema.DecisionApprove, Reason: "ok"}}

	exitCode := Run(context.Background(), stdin, &stdout, Options{BuiltinName: "validate_find"}, engine)
	if exitCode != schema.ExitBlock {
		t.Fatalf("exit code = %d", exitCode)
	}
	if engine.Calls != 0 {
		t.Fatal("engine should not be called for invalid input")
	}
}

func TestRun_BlocksWhenConfigMissing(t *testing.T) {
	stdin := strings.NewReader(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)
	var stdout bytes.Buffer
	engine := &judge.MockEngine{Result: &schema.JudgeResult{Decision: schema.DecisionApprove, Reason: "ok"}}

	exitCode := Run(context.Background(), stdin, &stdout, Options{}, engine)
	if exitCode != schema.ExitBlock {
		t.Fatalf("exit code = %d", exitCode)
	}
}

func TestRun_BlocksWhenBothConfigFlagsSet(t *testing.T) {
	stdin := strings.NewReader(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)
	var stdout bytes.Buffer
	engine := &judge.MockEngine{Result: &schema.JudgeResult{Decision: schema.DecisionApprove, Reason: "ok"}}

	exitCode := Run(context.Background(), stdin, &stdout, Options{
		ConfigPath:  "/tmp/config.yaml",
		BuiltinName: "validate_find",
	}, engine)
	if exitCode != schema.ExitBlock {
		t.Fatalf("exit code = %d", exitCode)
	}
}

func TestRun_BlocksOnJudgeError(t *testing.T) {
	stdin := strings.NewReader(`{
		"hook_event_name": "PreToolUse",
		"tool_name": "exec",
		"tool_input": {"command": "ls"}
	}`)
	var stdout bytes.Buffer
	engine := &judge.MockEngine{Err: context.DeadlineExceeded}

	exitCode := Run(context.Background(), stdin, &stdout, Options{BuiltinName: "validate_find"}, engine)
	if exitCode != schema.ExitBlock {
		t.Fatalf("exit code = %d", exitCode)
	}
}
