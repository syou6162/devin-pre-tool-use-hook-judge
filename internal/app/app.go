package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/judge"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

// Options configure a hook run.
type Options struct {
	ConfigPath  string
	BuiltinName string
}

// Run executes the PreToolUse hook validator.
func Run(ctx context.Context, stdin io.Reader, stdout io.Writer, opts Options, engine judge.Engine) int {
	cfg, err := loadConfig(opts)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("設定ファイル読み込みエラー: %v", err))
	}

	raw, err := io.ReadAll(stdin)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("入力読み込みエラー: %v", err))
	}

	input, err := schema.ParseDevinInput(raw)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("入力検証エラー: %v", err))
	}

	judgeInput := schema.ToJudgeInput(input)
	result, err := engine.Judge(ctx, judgeInput, cfg)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("判定エンジンエラー: %v", err))
	}

	output := schema.ToDevinOutput(result)
	return writeOutput(stdout, output)
}

func loadConfig(opts Options) (*config.Config, error) {
	if opts.ConfigPath != "" && opts.BuiltinName != "" {
		return nil, fmt.Errorf("--configと--builtinの両方を指定することはできません")
	}
	if opts.ConfigPath != "" {
		return config.LoadFromFile(opts.ConfigPath)
	}
	if opts.BuiltinName != "" {
		return config.LoadBuiltin(opts.BuiltinName)
	}
	return nil, fmt.Errorf("設定ファイルが指定されていません")
}

func writeBlock(stdout io.Writer, reason string) int {
	return writeOutput(stdout, schema.BlockOutput(reason))
}

func writeOutput(stdout io.Writer, output schema.DevinOutput) int {
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output)
	return schema.ExitCodeForOutput(output)
}

// ResolveBuiltinRoot returns the directory containing embedded builtin configs for tests.
func ResolveBuiltinRoot() string {
	if dir := os.Getenv("DEVIN_HOOK_BUILTIN_ROOT"); dir != "" {
		return dir
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

// NormalizeReason trims whitespace for stable test comparisons.
func NormalizeReason(reason string) string {
	return strings.TrimSpace(reason)
}
