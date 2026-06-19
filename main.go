package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/hook"
)

//go:embed builtin_configs/*.yaml
var builtinConfigs embed.FS

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, builtinConfigs))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer, builtins embed.FS) int {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	fs.SetOutput(stderr)

	configPath := fs.String("config", "", "Path to external YAML configuration file")
	builtinName := fs.String("builtin", "", "Name of builtin configuration to use")

	if err := fs.Parse(args); err != nil {
		writeConfigBlock(stdout, fmt.Sprintf("フラグの解析に失敗しました: %v", err))
		return 2
	}

	cfg, err := config.LoadFromFlags(*configPath, *builtinName, builtins)
	if err != nil {
		writeConfigBlock(stdout, fmt.Sprintf("設定ファイル読み込みエラー: %s", err.Error()))
		return 2
	}

	return runWithConfig(cfg, stdin, stdout)
}

func runWithConfig(cfg config.Config, stdin io.Reader, stdout io.Writer) int {
	if _, err := io.Copy(io.Discard, stdin); err != nil {
		writeConfigBlock(stdout, fmt.Sprintf("入力の読み込みに失敗しました: %v", err))
		return 2
	}

	_ = cfg
	// Judge engine integration is implemented in a follow-up issue.
	if err := hook.WriteBlock(stdout, "判定エンジンは未実装です。安全のため操作を拒否します。"); err != nil {
		return 2
	}
	return 2
}

func writeConfigBlock(stdout io.Writer, reason string) {
	_ = hook.WriteBlock(stdout, reason)
}
