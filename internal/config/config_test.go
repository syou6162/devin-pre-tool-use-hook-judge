package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadBuiltinConfigSuccess(t *testing.T) {
	t.Parallel()

	builtins := testBuiltinFS(t)
	cfg, err := LoadBuiltinConfig("example", builtins)
	if err != nil {
		t.Fatalf("LoadBuiltinConfig() error = %v", err)
	}
	if !strings.Contains(cfg.Prompt, "security validator") {
		t.Fatalf("Prompt = %q, want security validator text", cfg.Prompt)
	}
	if cfg.Model != "swe-1-6-fast" {
		t.Errorf("Model = %q, want swe-1-6-fast", cfg.Model)
	}
	if cfg.Timeout != 45 {
		t.Errorf("Timeout = %d, want 45", cfg.Timeout)
	}
	if len(cfg.AllowedTools) != 1 || cfg.AllowedTools[0] != "exec" {
		t.Errorf("AllowedTools = %v, want [exec]", cfg.AllowedTools)
	}
}

func TestLoadBuiltinConfigNotFound(t *testing.T) {
	t.Parallel()

	_, err := LoadBuiltinConfig("missing", testBuiltinFS(t))
	if err == nil {
		t.Fatal("LoadBuiltinConfig() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %q, want not found", err.Error())
	}
}

func TestLoadConfigSuccess(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "valid.yaml")
	writeFile(t, path, `
prompt: "Test validation prompt"
model: claude-sonnet-4.5
timeout: 120
allowed_tools:
  - exec
  - read
`)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Prompt != "Test validation prompt" {
		t.Errorf("Prompt = %q", cfg.Prompt)
	}
	if cfg.Model != "claude-sonnet-4.5" {
		t.Errorf("Model = %q", cfg.Model)
	}
	if cfg.Timeout != 120 {
		t.Errorf("Timeout = %d", cfg.Timeout)
	}
	if len(cfg.AllowedTools) != 2 {
		t.Errorf("AllowedTools = %v", cfg.AllowedTools)
	}
}

func TestLoadConfigAppliesDefaults(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.yaml")
	writeFile(t, path, `
prompt: "Minimal prompt"
`)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Model != DefaultModel {
		t.Errorf("Model = %q, want %q", cfg.Model, DefaultModel)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d, want %d", cfg.Timeout, DefaultTimeout)
	}
}

func TestLoadConfigMultilinePrompt(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "multiline.yaml")
	writeFile(t, path, `
prompt: |
  Line one
  Line two
model: haiku
`)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if !strings.Contains(cfg.Prompt, "Line one") || !strings.Contains(cfg.Prompt, "Line two") {
		t.Errorf("Prompt = %q", cfg.Prompt)
	}
	if cfg.Model != "haiku" {
		t.Errorf("Model = %q", cfg.Model)
	}
}

func TestLoadConfigErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		path    string
		want    string
	}{
		{
			name: "missing file",
			path: "/nonexistent/path/config.yaml",
			want: "not found",
		},
		{
			name:    "empty file",
			content: "",
			want:    "empty",
		},
		{
			name: "invalid yaml",
			content: `prompt: "broken
model: bad`,
			want: "Failed to parse",
		},
		{
			name: "missing prompt",
			content: `model: swe-1-6-fast`,
			want: "prompt",
		},
		{
			name: "invalid model",
			content: `prompt: "test"
model: invalid-model`,
			want: "Validation failed",
		},
		{
			name: "invalid prompt type",
			content: `prompt: 123`,
			want: "Validation failed",
		},
		{
			name: "invalid allowed_tools type",
			content: `prompt: "test"
allowed_tools: exec`,
			want: "Validation failed",
		},
		{
			name: "invalid timeout",
			content: `prompt: "test"
timeout: 0`,
			want: "Validation failed",
		},
		{
			name: "unknown field",
			content: `prompt: "test"
extra: true`,
			want: "Validation failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.path
			if path == "" {
				dir := t.TempDir()
				path = filepath.Join(dir, "config.yaml")
				writeFile(t, path, tt.content)
			}

			_, err := LoadConfig(path)
			if err == nil {
				t.Fatal("LoadConfig() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tt.want)
			}
		})
	}
}

func TestLoadFromFlags(t *testing.T) {
	t.Parallel()

	builtins := testBuiltinFS(t)

	t.Run("mutual exclusion", func(t *testing.T) {
		t.Parallel()
		_, err := LoadFromFlags("a.yaml", "example", builtins)
		if err == nil || !strings.Contains(err.Error(), "--configと--builtin") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("no flags", func(t *testing.T) {
		t.Parallel()
		_, err := LoadFromFlags("", "", builtins)
		if err == nil || !strings.Contains(err.Error(), "設定ファイルが指定されていません") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestApplyDefaults(t *testing.T) {
	t.Parallel()

	cfg := Config{Prompt: "test"}
	cfg.ApplyDefaults()
	if cfg.Model != DefaultModel {
		t.Errorf("Model = %q", cfg.Model)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout = %d", cfg.Timeout)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func testBuiltinFS(t *testing.T) fs.FS {
	t.Helper()
	return os.DirFS(filepath.Join("testdata"))
}
