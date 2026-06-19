package schema

import "encoding/json"

const (
	HookEventNamePreToolUse = "PreToolUse"

	DecisionApprove = "approve"
	DecisionBlock   = "block"
	DecisionDeny    = "deny"

	ExitApprove = 0
	ExitBlock   = 2

	DefaultSessionID       = "unknown"
	DefaultTranscriptPath  = ""
	DefaultCWD             = "."
	DefaultPermissionMode  = "default"
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
}

// DevinOutput is the JSON payload returned to Devin CLI command hooks.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// JudgeInput is the normalized input passed to the judgment engine.
type JudgeInput struct {
	SessionID      string          `json:"session_id"`
	TranscriptPath string          `json:"transcript_path"`
	CWD            string          `json:"cwd"`
	PermissionMode string          `json:"permission_mode"`
	HookEventName  string          `json:"hook_event_name"`
	ToolName       string          `json:"tool_name"`
	ToolInput      json.RawMessage `json:"tool_input"`
}

// JudgeResult is the JSON schema expected from the judgment engine.
type JudgeResult struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}
