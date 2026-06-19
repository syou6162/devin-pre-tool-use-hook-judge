package judge

import (
	"encoding/json"
	"fmt"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

const inputJSONSchema = `{
  "type": "object",
  "required": [
    "session_id",
    "transcript_path",
    "cwd",
    "permission_mode",
    "hook_event_name",
    "tool_name",
    "tool_input"
  ],
  "properties": {
    "session_id": {"type": "string"},
    "transcript_path": {"type": "string"},
    "cwd": {"type": "string"},
    "permission_mode": {"type": "string"},
    "hook_event_name": {"type": "string", "const": "PreToolUse"},
    "tool_name": {"type": "string"},
    "tool_input": {"type": "object"}
  }
}`

const outputJSONSchema = `{
  "type": "object",
  "required": ["hookSpecificOutput"],
  "properties": {
    "hookSpecificOutput": {
      "type": "object",
      "required": [
        "hookEventName",
        "permissionDecision",
        "permissionDecisionReason"
      ],
      "properties": {
        "hookEventName": {"type": "string", "const": "PreToolUse"},
        "permissionDecision": {
          "type": "string",
          "enum": ["allow", "deny", "ask", "defer"]
        },
        "permissionDecisionReason": {"type": "string"},
        "updatedInput": {"type": "object"},
        "additionalContext": {"type": "string"}
      }
    }
  }
}`

const systemPromptTemplate = `You are a PreToolUse hook validator for Devin CLI.

Your task is to validate tool usage and return a decision based on the validation rules provided in <custom_validation_rules>.

The input structure is:
<input_json_schema>
%s
</input_json_schema>

You MUST respond with ONLY a valid JSON object matching this output schema (no markdown, no code fences, no extra text):
<output_json_schema>
%s
</output_json_schema>

If you cannot determine whether to allow or deny based on the provided rules, default to DENY for safety.`

// buildPrompt generates the full prompt for devin --print.
func buildPrompt(input *schema.JudgeInput, customPrompt string) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input is nil")
	}
	if customPrompt == "" {
		return "", fmt.Errorf("custom prompt is required")
	}

	inputJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal input: %w", err)
	}

	systemPrompt := fmt.Sprintf(systemPromptTemplate, inputJSONSchema, outputJSONSchema)

	return fmt.Sprintf(`<system_instructions>
%s
</system_instructions>

<custom_validation_rules>
%s
</custom_validation_rules>

# Current Tool Usage
%s
`, systemPrompt, customPrompt, string(inputJSON)), nil
}
