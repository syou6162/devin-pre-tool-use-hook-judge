package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultModel   = "haiku"
	DefaultTimeout = 60 * time.Second
)

// Config holds YAML validator settings.
type Config struct {
	Prompt         string        `yaml:"prompt"`
	Model          string        `yaml:"model,omitempty"`
	TimeoutSeconds int           `yaml:"timeout,omitempty"`
	AllowedTools   []string      `yaml:"allowed_tools,omitempty"`
	Timeout        time.Duration `yaml:"-"`
}

// ConfigError indicates configuration loading or validation failed.
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

// LoadFromFile loads and validates a YAML configuration file from disk.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &ConfigError{Message: fmt.Sprintf("config file '%s' not found", path)}
		}
		return nil, &ConfigError{Message: fmt.Sprintf("failed to read config file '%s': %v", path, err)}
	}
	return parseAndValidate(data, fmt.Sprintf("config file '%s'", path))
}

// LoadBuiltin loads a builtin configuration by name from the provided filesystem.
func LoadBuiltin(name string, configs fs.FS) (*Config, error) {
	if err := validateBuiltinName(name); err != nil {
		return nil, err
	}

	path := filepath.Join("builtin_configs", name+".yaml")
	data, err := fs.ReadFile(configs, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, &ConfigError{Message: fmt.Sprintf("builtin config '%s' not found", name)}
		}
		return nil, &ConfigError{Message: fmt.Sprintf("failed to read builtin config '%s': %v", name, err)}
	}

	return parseAndValidate(data, fmt.Sprintf("builtin config '%s'", name))
}

func validateBuiltinName(name string) error {
	if name == "" {
		return &ConfigError{Message: "builtin config name is required"}
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		return &ConfigError{Message: fmt.Sprintf("invalid builtin config name '%s'", name)}
	}
	return nil
}

func parseAndValidate(data []byte, source string) (*Config, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, &ConfigError{Message: fmt.Sprintf("failed to parse %s: %v", source, err)}
	}
	if raw == nil {
		return nil, &ConfigError{Message: fmt.Sprintf("%s is empty", source)}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &ConfigError{Message: fmt.Sprintf("failed to parse %s: %v", source, err)}
	}

	if err := Validate(&cfg); err != nil {
		return nil, &ConfigError{Message: fmt.Sprintf("validation failed for %s: %v", source, err)}
	}

	ApplyDefaults(&cfg)
	return &cfg, nil
}

// Validate checks required fields and rejects unknown YAML keys.
func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(cfg.Prompt) == "" {
		return fmt.Errorf("prompt is required")
	}
	return nil
}

// ApplyDefaults fills in model and timeout when omitted.
func ApplyDefaults(cfg *Config) {
	if cfg.Model == "" {
		cfg.Model = DefaultModel
	}
	switch {
	case cfg.TimeoutSeconds > 0:
		cfg.Timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	default:
		cfg.Timeout = DefaultTimeout
	}
}
