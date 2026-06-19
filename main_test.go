package main

import (
	"testing"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
)

func TestMainBuiltinConfigsEmbedded(t *testing.T) {
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

			cfg, err := config.LoadBuiltin(name, builtinConfigs)
			if err != nil {
				t.Fatalf("LoadBuiltin(%q) from embed error: %v", name, err)
			}
			if cfg.Prompt == "" {
				t.Fatalf("LoadBuiltin(%q) prompt is empty", name)
			}
		})
	}
}
