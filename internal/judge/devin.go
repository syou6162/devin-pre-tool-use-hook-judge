package judge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/jsonutil"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

const maxRetries = 3

const promptTemplate = `You are a PreToolUse hook validator for Devin CLI.

Evaluate the tool usage request below according to the validation rules.
Respond with JSON only in this exact shape:
{"decision":"approve"|"deny","reason":"<short explanation>"}

If you cannot determine whether to approve or deny, default to deny for safety.

<validation_rules>
{{.Prompt}}
</validation_rules>

<tool_usage_request>
{{.InputJSON}}
</tool_usage_request>
`

// DevinEngine calls `devin --print` to perform judgments.
type DevinEngine struct {
	DevinPath string
}

// NewDevinEngine creates an engine that invokes the Devin CLI.
func NewDevinEngine() *DevinEngine {
	return &DevinEngine{DevinPath: "devin"}
}

// Judge runs the configured prompt against Devin CLI with retries.
func (e *DevinEngine) Judge(ctx context.Context, input *schema.JudgeInput, cfg *config.Config) (*schema.JudgeResult, error) {
	inputJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal judge input: %w", err)
	}

	prompt, err := renderPrompt(cfg.Prompt, string(inputJSON))
	if err != nil {
		return nil, err
	}

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := e.runOnce(ctx, prompt, cfg)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("judge failed after %d attempts: %w", maxRetries, lastErr)
}

func renderPrompt(customPrompt, inputJSON string) (string, error) {
	tmpl, err := template.New("prompt").Parse(promptTemplate)
	if err != nil {
		return "", fmt.Errorf("parse prompt template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"Prompt":    customPrompt,
		"InputJSON": inputJSON,
	}); err != nil {
		return "", fmt.Errorf("render prompt template: %w", err)
	}
	return buf.String(), nil
}

func (e *DevinEngine) runOnce(ctx context.Context, prompt string, cfg *config.Config) (*schema.JudgeResult, error) {
	promptFile, err := os.CreateTemp("", "devin-pre-tool-use-hook-judge-*.txt")
	if err != nil {
		return nil, fmt.Errorf("create prompt file: %w", err)
	}
	defer os.Remove(promptFile.Name())

	if _, err := promptFile.WriteString(prompt); err != nil {
		promptFile.Close()
		return nil, fmt.Errorf("write prompt file: %w", err)
	}
	if err := promptFile.Close(); err != nil {
		return nil, fmt.Errorf("close prompt file: %w", err)
	}

	runCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	args := []string{"--print", "--prompt-file", filepath.Clean(promptFile.Name())}
	if cfg.Model != "" && cfg.Model != config.DefaultModel {
		args = append(args, "--model", cfg.Model)
	}

	cmd := exec.CommandContext(runCtx, e.DevinPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("devin command failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	jsonText, err := jsonutil.ExtractJSON(string(output))
	if err != nil {
		return nil, fmt.Errorf("extract JSON: %w", err)
	}

	var result schema.JudgeResult
	if err := json.Unmarshal([]byte(jsonText), &result); err != nil {
		return nil, fmt.Errorf("parse judge JSON: %w", err)
	}
	if err := schema.ValidateJudgeResult(&result); err != nil {
		return nil, fmt.Errorf("validate judge result: %w", err)
	}
	if strings.TrimSpace(result.Reason) == "" {
		return nil, fmt.Errorf("judge result reason is empty")
	}

	return &result, nil
}
