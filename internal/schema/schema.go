package schema

import "encoding/json"

const (
	HookEventPreToolUse = "PreToolUse"

	DecisionApprove = "approve"
	DecisionBlock   = "block"
	DecisionDeny    = "deny"

	ExitCodeApprove = 0
	ExitCodeBlock   = 2

	DefaultPermissionMode = "default"
)

var validPermissionModes = map[string]struct{}{
	"default":            {},
	"plan":               {},
	"acceptEdits":        {},
	"auto":               {},
	"dontAsk":            {},
	"bypassPermissions":  {},
}

// DevinInput is the JSON payload received from Devin CLI command hooks.
type DevinInput struct {
	HookEventName  string          `json:"hook_event_name"`
	ToolName       string          `json:"tool_name"`
	ToolInput      json.RawMessage `json:"tool_input"`
	SessionID      string          `json:"session_id,omitempty"`
	TranscriptPath string          `json:"transcript_path,omitempty"`
	CWD            string          `json:"cwd,omitempty"`
	PermissionMode string          `json:"permission_mode,omitempty"`
	ToolUseID      string          `json:"tool_use_id,omitempty"`
	Timestamp      string          `json:"timestamp,omitempty"`
}

// DevinOutput is the JSON payload returned to Devin CLI command hooks.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// JudgeInput is the internal representation used by the judgment engine.
// It mirrors the Claude Code PreToolUse hook input format.
type JudgeInput struct {
	SessionID      string                 `json:"session_id"`
	TranscriptPath string                 `json:"transcript_path"`
	CWD            string                 `json:"cwd"`
	PermissionMode string                 `json:"permission_mode"`
	HookEventName  string                 `json:"hook_event_name"`
	ToolName       string                 `json:"tool_name"`
	ToolInput      map[string]interface{} `json:"tool_input"`
}
