// Package constants defines shared application constants.
package constants

const (
	Version = "0.1.0"

	HookEventPreToolUse = "PreToolUse"

	PermissionAllow = "allow"
	PermissionDeny  = "deny"
	PermissionAsk   = "ask"

	DefaultPermissionDecision = PermissionDeny
	DefaultPermissionReason   = "Invalid or missing permission decision from judgment system"
)
