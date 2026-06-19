package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/builtin"
	"gopkg.in/yaml.v3"
)

const (
	DefaultModel   = "default"
	DefaultTimeout = 120 * time.Second
)

// Config holds YAML validator settings.
type Config struct {
	Prompt  string        `yaml:"prompt"`
	Model   string        `yaml:"model,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
}

// LoadFromFile loads and validates a YAML config file from disk.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	return parseConfig(data, path)
}

// LoadBuiltin loads a named builtin config from the embedded filesystem.
func LoadBuiltin(name string) (*Config, error) {
	if strings.Contains(name, "/") || strings.Contains(name, "..") {
		return nil, fmt.Errorf("invalid builtin name: %q", name)
	}

	path := filepath.Join("configs", name+".yaml")
	data, err := builtin.FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("builtin config %q not found", name)
	}
	return parseConfig(data, name)
}

func parseConfig(data []byte, source string) (*Config, error) {
	var raw struct {
		Prompt  string `yaml:"prompt"`
		Model   string `yaml:"model,omitempty"`
		Timeout string `yaml:"timeout,omitempty"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", source, err)
	}

	cfg := &Config{
		Prompt:  strings.TrimSpace(raw.Prompt),
		Model:   raw.Model,
		Timeout: DefaultTimeout,
	}

	if cfg.Prompt == "" {
		return nil, fmt.Errorf("config %q: prompt is required", source)
	}

	if raw.Model == "" {
		cfg.Model = DefaultModel
	}

	if raw.Timeout != "" {
		timeout, err := time.ParseDuration(raw.Timeout)
		if err != nil {
			return nil, fmt.Errorf("config %q: invalid timeout: %w", source, err)
		}
		if timeout <= 0 {
			return nil, fmt.Errorf("config %q: timeout must be positive", source)
		}
		cfg.Timeout = timeout
	}

	return cfg, nil
}
