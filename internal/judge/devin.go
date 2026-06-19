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
	"time"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

const maxRetries = 3

type CommandRunner func(ctx context.Context, name string, args ...string) *exec.Cmd

type DevinEngine struct {
	Binary  string
	Runner  CommandRunner
	NowFunc func() time.Time
}

func NewDevinEngine(binary string) *DevinEngine {
	if binary == "" {
		binary = "devin"
	}
	return &DevinEngine{
		Binary: binary,
		Runner: exec.CommandContext,
		NowFunc: time.Now,
	}
}

func (e *DevinEngine) Judge(ctx context.Context, input schema.JudgeInput, cfg config.Config) (schema.JudgeOutput, error) {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		output, err := e.runOnce(ctx, input, cfg)
		if err == nil {
			return output, nil
		}
		lastErr = err
	}
	return schema.JudgeOutput{}, fmt.Errorf("judgment failed after %d attempts: %w", maxRetries, lastErr)
}

func (e *DevinEngine) runOnce(ctx context.Context, input schema.JudgeInput, cfg config.Config) (schema.JudgeOutput, error) {
	prompt, err := buildPrompt(input, cfg.Prompt)
	if err != nil {
		return schema.JudgeOutput{}, err
	}

	promptFile, err := writePromptFile(prompt)
	if err != nil {
		return schema.JudgeOutput{}, err
	}
	defer os.Remove(promptFile)

	runCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	args := []string{"--print", "--prompt-file", promptFile}
	if cfg.Model != "" {
		args = append([]string{"--model", cfg.Model}, args...)
	}

	cmd := e.Runner(runCtx, e.Binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return schema.JudgeOutput{}, fmt.Errorf("devin command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	jsonText, err := ExtractJSON(stdout.String())
	if err != nil {
		return schema.JudgeOutput{}, err
	}

	var output schema.JudgeOutput
	if err := json.Unmarshal([]byte(jsonText), &output); err != nil {
		return schema.JudgeOutput{}, fmt.Errorf("failed to parse judgment JSON: %w", err)
	}
	if err := schema.ValidateJudgeOutput(output); err != nil {
		return schema.JudgeOutput{}, err
	}
	return output, nil
}

func buildPrompt(input schema.JudgeInput, rules string) (string, error) {
	inputJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal judge input: %w", err)
	}

	var builder strings.Builder
	builder.WriteString("You are a PreToolUse hook validator for Devin CLI.\n\n")
	builder.WriteString("Validate the tool usage below using the custom validation rules.\n")
	builder.WriteString("Respond with JSON only in this exact shape:\n")
	builder.WriteString(`{"permissionDecision":"allow|deny|ask","permissionDecisionReason":"..."}` + "\n\n")
	builder.WriteString("If you cannot determine whether to allow or deny, default to deny for safety.\n\n")
	builder.WriteString("<custom_validation_rules>\n")
	builder.WriteString(strings.TrimSpace(rules))
	builder.WriteString("\n</custom_validation_rules>\n\n")
	builder.WriteString("# Current Tool Usage\n")
	builder.Write(inputJSON)
	builder.WriteByte('\n')
	return builder.String(), nil
}

func writePromptFile(prompt string) (string, error) {
	dir := os.TempDir()
	file, err := os.CreateTemp(dir, "devin-pre-tool-use-hook-judge-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create prompt file: %w", err)
	}

	path, err := filepath.Abs(file.Name())
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return "", err
	}

	if _, err := file.WriteString(prompt); err != nil {
		file.Close()
		os.Remove(path)
		return "", fmt.Errorf("failed to write prompt file: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to close prompt file: %w", err)
	}
	return path, nil
}
