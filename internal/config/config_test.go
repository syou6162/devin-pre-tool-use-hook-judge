package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadBuiltin_ValidateGitPush(t *testing.T) {
	cfg, err := LoadBuiltin("validate_git_push")
	if err != nil {
		t.Fatalf("LoadBuiltin() error = %v", err)
	}
	if cfg.Prompt == "" {
		t.Fatal("prompt should not be empty")
	}
	if cfg.Model == "" {
		t.Fatal("model should have default")
	}
	if cfg.Timeout != DefaultTimeout {
		t.Fatalf("timeout = %v, want %v", cfg.Timeout, DefaultTimeout)
	}
}

func TestLoadBuiltin_NotFound(t *testing.T) {
	_, err := LoadBuiltin("does_not_exist")
	if err == nil {
		t.Fatal("expected error for missing builtin")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := "prompt: |\n  Validate commands\nmodel: sonnet\ntimeout: 30s\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile() error = %v", err)
	}
	if cfg.Prompt != "Validate commands" {
		t.Fatalf("prompt = %q", cfg.Prompt)
	}
	if cfg.Model != "sonnet" {
		t.Fatalf("model = %q", cfg.Model)
	}
	if cfg.Timeout != 30*time.Second {
		t.Fatalf("timeout = %v", cfg.Timeout)
	}
}

func TestLoadFromFile_MissingPrompt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("model: sonnet\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadFromFile(path)
	if err == nil {
		t.Fatal("expected prompt validation error")
	}
}

func TestLoadBuiltin_InvalidName(t *testing.T) {
	_, err := LoadBuiltin("../escape")
	if err == nil {
		t.Fatal("expected invalid name error")
	}
}
