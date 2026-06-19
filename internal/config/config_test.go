package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestLoadBuiltinAllConfigs(t *testing.T) {
	t.Parallel()

	names := []string{
		"validate_bq_query",
		"validate_codex_mcp",
		"validate_git_push",
		"validate_find",
		"validate_xargs",
	}

	for _, name := range names {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := LoadBuiltin(name, repoBuiltinFS(t))
			if err != nil {
				t.Fatalf("LoadBuiltin(%q) error: %v", name, err)
			}
			if cfg.Prompt == "" {
				t.Fatalf("LoadBuiltin(%q) prompt is empty", name)
			}
			if cfg.Model != DefaultModel {
				t.Fatalf("LoadBuiltin(%q) model = %q, want %q", name, cfg.Model, DefaultModel)
			}
			if cfg.Timeout != DefaultTimeout {
				t.Fatalf("LoadBuiltin(%q) timeout = %v, want %v", name, cfg.Timeout, DefaultTimeout)
			}
		})
	}
}

func TestLoadBuiltinEmbeddedFS(t *testing.T) {
	t.Parallel()

	cfg, err := LoadBuiltin("validate_bq_query", repoBuiltinFS(t))
	if err != nil {
		t.Fatalf("LoadBuiltin error: %v", err)
	}
	if cfg.Prompt == "" {
		t.Fatal("expected non-empty prompt")
	}
}

func TestLoadBuiltinNotFound(t *testing.T) {
	t.Parallel()

	_, err := LoadBuiltin("missing_config", fstest.MapFS{})
	if err == nil {
		t.Fatal("expected error for missing builtin config")
	}
}

func TestLoadBuiltinInvalidName(t *testing.T) {
	t.Parallel()

	_, err := LoadBuiltin("../escape", fstest.MapFS{})
	if err == nil {
		t.Fatal("expected error for invalid builtin name")
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.yaml")
	content := "prompt: |\n  test prompt\nmodel: sonnet\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile error: %v", err)
	}
	if cfg.Prompt != "test prompt\n" {
		t.Fatalf("prompt = %q, want %q", cfg.Prompt, "test prompt\n")
	}
	if cfg.Model != "sonnet" {
		t.Fatalf("model = %q, want sonnet", cfg.Model)
	}
}

func TestValidateRequiresPrompt(t *testing.T) {
	t.Parallel()

	err := Validate(&Config{})
	if err == nil {
		t.Fatal("expected validation error for empty prompt")
	}
}

func repoBuiltinFS(t *testing.T) fs.FS {
	t.Helper()
	return os.DirFS(repoRoot(t))
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "builtin_configs", "validate_bq_query.yaml")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate repo root with builtin_configs")
		}
		dir = parent
	}
}
