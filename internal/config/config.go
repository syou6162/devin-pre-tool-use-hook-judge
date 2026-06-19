package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/builtinconfigs"
	"gopkg.in/yaml.v3"
)


const (
	DefaultModel   = "sonnet"
	DefaultTimeout = 60 * time.Second
)

type Config struct {
	Prompt  string        `yaml:"prompt"`
	Model   string        `yaml:"model"`
	Timeout time.Duration `yaml:"timeout"`
}

func (c *Config) ApplyDefaults() {
	if c.Model == "" {
		c.Model = DefaultModel
	}
	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}
}

func LoadFile(path string) (Config, error) {
	data, err := readFile(path)
	if err != nil {
		return Config{}, err
	}
	return parseYAML(data, path)
}

func LoadBuiltin(name string) (Config, error) {
	if strings.Contains(name, "/") || strings.Contains(name, "..") {
		return Config{}, fmt.Errorf("invalid builtin name %q", name)
	}

	path := filepath.Join(name + ".yaml")
	data, err := builtinconfigs.Files.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("builtin config %q not found", name)
	}

	cfg, err := parseYAML(data, name)
	if err != nil {
		return Config{}, fmt.Errorf("builtin config %q: %w", name, err)
	}
	return cfg, nil
}

func parseYAML(data []byte, source string) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config %q: %w", source, err)
	}
	if strings.TrimSpace(cfg.Prompt) == "" {
		return Config{}, fmt.Errorf("config %q: prompt is required", source)
	}
	cfg.ApplyDefaults()
	return cfg, nil
}

var readFile = defaultReadFile

func defaultReadFile(path string) ([]byte, error) {
	return readFileFromOS(path)
}
