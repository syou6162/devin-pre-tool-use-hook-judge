package schema

const (
	HookEventName = "PreToolUse"

	PermissionAllow = "allow"
	PermissionDeny  = "deny"
	PermissionAsk   = "ask"
	PermissionDefer = "defer"

	DefaultPermissionMode = "default"
	DefaultSessionID      = "unknown"
	DefaultTranscriptPath = ""
	DefaultCWD            = "."
)
