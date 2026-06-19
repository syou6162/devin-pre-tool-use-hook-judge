package cli

import "testing"

func TestParseVersionFlag(t *testing.T) {
	opts, err := Parse([]string{"--version"})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !opts.Version {
		t.Error("expected Version to be true")
	}
}

func TestParseConfigAndBuiltinFlags(t *testing.T) {
	opts, err := Parse([]string{"--config", "rules.yaml", "--builtin", "validate_exec"})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if opts.Config != "rules.yaml" {
		t.Errorf("Config = %q, want %q", opts.Config, "rules.yaml")
	}
	if opts.Builtin != "validate_exec" {
		t.Errorf("Builtin = %q, want %q", opts.Builtin, "validate_exec")
	}
}
