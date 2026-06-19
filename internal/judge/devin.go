package judge

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

// DevinRunner executes devin --print with a prompt file.
type DevinRunner interface {
	Run(ctx context.Context, promptFile string, model string) (string, error)
}

// ExecDevinRunner runs the devin CLI as a subprocess.
type ExecDevinRunner struct {
	DevinPath string
	Timeout   time.Duration
}

// Run executes devin --print --prompt-file and returns stdout.
func (r *ExecDevinRunner) Run(ctx context.Context, promptFile string, model string) (string, error) {
	devinPath := r.DevinPath
	if devinPath == "" {
		devinPath = "devin"
	}

	args := []string{"--print", "--prompt-file", promptFile}
	if model != "" {
		args = append(args, "--model", model)
	}

	timeout := r.Timeout
	if timeout == 0 {
		timeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, devinPath, args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("devin exited with code %d: %s", exitErr.ExitCode(), string(exitErr.Stderr))
		}
		return "", fmt.Errorf("run devin: %w", err)
	}

	return string(output), nil
}

// DevinEngine calls devin --print to judge tool usage.
type DevinEngine struct {
	Runner DevinRunner
}

// NewDevinEngine creates a DevinEngine with the default subprocess runner.
func NewDevinEngine(devinPath string, timeout time.Duration) *DevinEngine {
	return &DevinEngine{
		Runner: &ExecDevinRunner{
			DevinPath: devinPath,
			Timeout:   timeout,
		},
	}
}

// Judge validates the input using devin --print with up to maxRetries attempts.
func (e *DevinEngine) Judge(ctx context.Context, input *schema.JudgeInput, opts JudgeOptions) (*schema.HookOutput, error) {
	if e == nil || e.Runner == nil {
		return nil, fmt.Errorf("devin engine is not configured")
	}
	if err := schema.ValidateJudgeInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	prompt, err := buildPrompt(input, opts.CustomPrompt)
	if err != nil {
		return nil, err
	}

	promptFile, err := writeTempPrompt(prompt)
	if err != nil {
		return nil, err
	}
	defer os.Remove(promptFile)

	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		raw, err := e.Runner.Run(ctx, promptFile, opts.Model)
		if err != nil {
			lastErr = err
			continue
		}

		jsonText, err := extractJSON(raw)
		if err != nil {
			lastErr = err
			continue
		}

		output, err := schema.ParseHookOutput([]byte(jsonText))
		if err != nil {
			lastErr = err
			continue
		}

		return output, nil
	}

	return nil, &JudgeError{Attempts: maxRetries, Cause: lastErr}
}

func writeTempPrompt(prompt string) (string, error) {
	dir := os.TempDir()
	file, err := os.CreateTemp(dir, "devin-judge-prompt-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temp prompt file: %w", err)
	}

	name := file.Name()
	if _, err := file.WriteString(prompt); err != nil {
		file.Close()
		os.Remove(name)
		return "", fmt.Errorf("write temp prompt file: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(name)
		return "", fmt.Errorf("close temp prompt file: %w", err)
	}

	return filepath.Clean(name), nil
}
