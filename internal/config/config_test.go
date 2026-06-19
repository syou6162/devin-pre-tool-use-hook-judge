package config

import (
	"strings"
	"testing"
)

func TestLoadBuiltin(t *testing.T) {
	t.Parallel()

	cfg, err := LoadBuiltin("validate_git_push")
	if err != nil {
		t.Fatalf("LoadBuiltin() error = %v", err)
	}
	if !strings.Contains(cfg.Prompt, "git push") {
		t.Fatalf("prompt should mention git push")
	}
	if cfg.Model == "" {
		t.Fatal("expected default model to be applied")
	}
	if cfg.Timeout == 0 {
		t.Fatal("expected default timeout to be applied")
	}
}

func TestLoadFile(t *testing.T) {
	t.Parallel()

	readFile = func(path string) ([]byte, error) {
		return []byte("prompt: test prompt\nmodel: haiku\n"), nil
	}
	t.Cleanup(func() { readFile = defaultReadFile })

	cfg, err := LoadFile("custom.yaml")
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if cfg.Prompt != "test prompt" {
		t.Fatalf("prompt = %q, want %q", cfg.Prompt, "test prompt")
	}
	if cfg.Model != "haiku" {
		t.Fatalf("model = %q, want %q", cfg.Model, "haiku")
	}
}

func TestLoadFileRequiresPrompt(t *testing.T) {
	t.Parallel()

	readFile = func(path string) ([]byte, error) {
		return []byte("model: haiku\n"), nil
	}
	t.Cleanup(func() { readFile = defaultReadFile })

	if _, err := LoadFile("invalid.yaml"); err == nil {
		t.Fatal("expected error for missing prompt")
	}
}
