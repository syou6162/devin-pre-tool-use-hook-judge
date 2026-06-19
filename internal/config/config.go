package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DefaultModel   = "swe-1-6-fast"
	DefaultTimeout = 60
)

// Config holds YAML validator settings.
type Config struct {
	Prompt       string   `yaml:"prompt"`
	Model        string   `yaml:"model,omitempty"`
	Timeout      int      `yaml:"timeout,omitempty"`
	AllowedTools []string `yaml:"allowed_tools,omitempty"`
}

// ApplyDefaults fills unspecified model and timeout fields.
func (c *Config) ApplyDefaults() {
	if strings.TrimSpace(c.Model) == "" {
		c.Model = DefaultModel
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
}

// LoadConfig reads and validates an external YAML configuration file.
func LoadConfig(path string) (Config, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, configErrorf("Config file '%s' not found", path)
		}
		return Config{}, configErrorf("Failed to read config file '%s': %v", path, err)
	}
	if info.IsDir() {
		return Config{}, configErrorf("Config file '%s' not found", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, configErrorf("Failed to read config file '%s': %v", path, err)
	}

	return parseAndValidate(data, path)
}

// LoadBuiltinConfig reads and validates an embedded builtin configuration file.
func LoadBuiltinConfig(name string, builtins fs.FS) (Config, error) {
	if strings.TrimSpace(name) == "" {
		return Config{}, configError("Builtin config name is required")
	}
	if strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return Config{}, configErrorf("Builtin config '%s' not found", name)
	}

	path := filepath.Join("builtin_configs", name+".yaml")
	data, err := fs.ReadFile(builtins, path)
	if err != nil {
		if os.IsNotExist(err) || errorsIsNotExist(err) {
			return Config{}, configErrorf("Builtin config '%s' not found", name)
		}
		return Config{}, configErrorf("Failed to read builtin config '%s': %v", name, err)
	}

	return parseAndValidate(data, name)
}

// LoadFromFlags loads configuration based on CLI flags.
// Exactly one of configPath or builtinName must be provided.
func LoadFromFlags(configPath, builtinName string, builtins fs.FS) (Config, error) {
	configPath = strings.TrimSpace(configPath)
	builtinName = strings.TrimSpace(builtinName)

	switch {
	case configPath != "" && builtinName != "":
		return Config{}, configError("--configと--builtinの両方を指定することはできません")
	case configPath != "":
		return LoadConfig(configPath)
	case builtinName != "":
		return LoadBuiltinConfig(builtinName, builtins)
	default:
		return Config{}, configError("設定ファイルが指定されていません。安全のため操作を拒否します。")
	}
}

func parseAndValidate(data []byte, source string) (Config, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return Config{}, configErrorf("Config file '%s' is empty", source)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, configErrorf("Failed to parse config file '%s': %v", source, err)
	}
	if raw == nil {
		return Config{}, configErrorf("Config file '%s' is empty", source)
	}

	if err := validateRaw(raw); err != nil {
		return Config{}, configErrorf("Validation failed for config file '%s': %v", source, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, configErrorf("Failed to parse config file '%s': %v", source, err)
	}

	cfg.ApplyDefaults()
	return cfg, nil
}

func validateRaw(raw map[string]any) error {
	allowedKeys := map[string]struct{}{
		"prompt":         {},
		"model":          {},
		"timeout":        {},
		"allowed_tools":  {},
	}

	for key := range raw {
		if _, ok := allowedKeys[key]; !ok {
			return fmt.Errorf("unknown field %q", key)
		}
	}

	prompt, ok := raw["prompt"]
	if !ok {
		return fmt.Errorf("'prompt' is a required property")
	}
	promptStr, ok := prompt.(string)
	if !ok {
		return fmt.Errorf("'prompt' must be a string")
	}
	if strings.TrimSpace(promptStr) == "" {
		return fmt.Errorf("'prompt' must not be empty")
	}

	if model, ok := raw["model"]; ok {
		modelStr, ok := model.(string)
		if !ok {
			return fmt.Errorf("'model' must be a string")
		}
		if strings.TrimSpace(modelStr) == "" {
			return fmt.Errorf("'model' must not be empty")
		}
		if !isValidModel(modelStr) {
			return fmt.Errorf("'model' has invalid value %q", modelStr)
		}
	}

	if timeout, ok := raw["timeout"]; ok {
		value, ok := positiveInt(timeout)
		if !ok || value <= 0 {
			return fmt.Errorf("'timeout' must be a positive integer")
		}
	}

	if allowedTools, ok := raw["allowed_tools"]; ok {
		items, ok := allowedTools.([]any)
		if !ok {
			return fmt.Errorf("'allowed_tools' must be an array")
		}
		for i, item := range items {
			if _, ok := item.(string); !ok {
				return fmt.Errorf("'allowed_tools' item at index %d must be a string", i)
			}
		}
	}

	return nil
}

func isValidModel(model string) bool {
	_, ok := validModels[model]
	return ok
}

var validModels = map[string]struct{}{
	"swe-1-6-fast":       {},
	"swe-1-6":            {},
	"claude-sonnet-4.5":  {},
	"claude-opus-4.6":    {},
	"gpt-5.2":            {},
	"gemini-3-pro":       {},
	"default":            {},
	"sonnet":             {},
	"opus":               {},
	"haiku":              {},
}

func errorsIsNotExist(err error) bool {
	return os.IsNotExist(err) || err == fs.ErrNotExist
}

func positiveInt(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case uint64:
		return int(v), true
	case float64:
		if v != float64(int(v)) {
			return 0, false
		}
		return int(v), true
	default:
		return 0, false
	}
}
