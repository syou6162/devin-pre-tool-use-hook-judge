package schema

import "encoding/json"

const (
	HookEventPreToolUse = "PreToolUse"

	DecisionApprove = "approve"
	DecisionBlock   = "block"
	DecisionDeny    = "deny"

	PermissionAllow = "allow"
	PermissionDeny  = "deny"
	PermissionAsk   = "ask"

	DefaultSessionID       = "unknown"
	DefaultCWD             = "."
	DefaultPermissionMode  = "default"
	DefaultPermissionReason = "Invalid or missing permission decision from judgment system"
)

// DevinInput is the JSON payload Devin CLI sends to command hooks.
type DevinInput struct {
	HookEventName string          `json:"hook_event_name"`
	ToolName      string          `json:"tool_name"`
	ToolInput     json.RawMessage `json:"tool_input"`
}

// DevinOutput is the JSON payload returned to Devin CLI command hooks.
type DevinOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason,omitempty"`
}

// JudgeInput is the internal representation used by the judgment engine.
type JudgeInput struct {
	SessionID      string          `json:"session_id"`
	TranscriptPath string          `json:"transcript_path,omitempty"`
	CWD            string          `json:"cwd"`
	PermissionMode string          `json:"permission_mode"`
	HookEventName  string          `json:"hook_event_name"`
	ToolName       string          `json:"tool_name"`
	ToolInput      json.RawMessage `json:"tool_input"`
}

// JudgeOutput is the structured response expected from the judgment engine.
type JudgeOutput struct {
	PermissionDecision       string `json:"permissionDecision"`
	PermissionDecisionReason string `json:"permissionDecisionReason"`
}
