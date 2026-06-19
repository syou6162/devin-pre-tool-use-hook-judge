package main

import (
	"bytes"
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/hook"
)

//go:embed builtin_configs/*.yaml
var testBuiltinConfigs embed.FS

func TestRunBlocksWhenNoConfigSpecified(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	code := run(nil, strings.NewReader("{}"), &stdout, io.Discard, testBuiltinConfigs)
	if code != 2 {
		t.Fatalf("run() = %d, want 2", code)
	}
	if !strings.Contains(stdout.String(), hook.DecisionBlock) {
		t.Fatalf("stdout = %q, want block decision", stdout.String())
	}
}

func TestRunBlocksWhenBothFlagsSpecified(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	code := run(
		[]string{"--config", "cfg.yaml", "--builtin", "example"},
		strings.NewReader("{}"),
		&stdout,
		io.Discard,
		testBuiltinConfigs,
	)
	if code != 2 {
		t.Fatalf("run() = %d, want 2", code)
	}
	if !strings.Contains(stdout.String(), "--configと--builtin") {
		t.Fatalf("stdout = %q, want mutual exclusion error", stdout.String())
	}
}

func TestRunLoadsBuiltinConfig(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	code := run(
		[]string{"--builtin", "example"},
		strings.NewReader(`{"hook_event_name":"PreToolUse"}`),
		&stdout,
		io.Discard,
		testBuiltinConfigs,
	)
	if code != 2 {
		t.Fatalf("run() = %d, want 2 until judge is implemented", code)
	}
}

func TestRunLoadsExternalConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte(`
prompt: "External prompt"
model: swe-1-6-fast
timeout: 30
`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	code := run(
		[]string{"--config", path},
		strings.NewReader(`{"hook_event_name":"PreToolUse"}`),
		&stdout,
		io.Discard,
		testBuiltinConfigs,
	)
	if code != 2 {
		t.Fatalf("run() = %d, want 2 until judge is implemented", code)
	}
}

func TestRunBlocksOnInvalidConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yaml")
	if err := os.WriteFile(path, []byte(`
model: swe-1-6-fast
`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout bytes.Buffer
	code := run(
		[]string{"--config", path},
		strings.NewReader("{}"),
		&stdout,
		io.Discard,
		testBuiltinConfigs,
	)
	if code != 2 {
		t.Fatalf("run() = %d, want 2", code)
	}
	if !strings.Contains(stdout.String(), "設定ファイル読み込みエラー") {
		t.Fatalf("stdout = %q, want config load error", stdout.String())
	}
}

func TestLoadFromFlagsAppliesDefaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.yaml")
	if err := os.WriteFile(path, []byte(`
prompt: "Minimal prompt"
`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := config.LoadFromFlags(path, "", testBuiltinConfigs)
	if err != nil {
		t.Fatalf("LoadFromFlags() error = %v", err)
	}
	if cfg.Model != config.DefaultModel {
		t.Errorf("Model = %q, want %q", cfg.Model, config.DefaultModel)
	}
	if cfg.Timeout != config.DefaultTimeout {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, config.DefaultTimeout)
	}
}
