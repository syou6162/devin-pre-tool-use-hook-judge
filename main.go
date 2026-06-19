package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/app"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/judge"
)

func main() {
	configPath := flag.String("config", "", "Path to external YAML configuration file")
	builtinName := flag.String("builtin", "", "Name of builtin configuration to use")
	flag.Parse()

	opts := app.Options{
		ConfigPath:  *configPath,
		BuiltinName: *builtinName,
	}

	engine := judge.NewDevinEngine()
	exitCode := app.Run(context.Background(), os.Stdin, os.Stdout, opts, engine)
	if err := os.Stdout.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "stdout sync: %v\n", err)
	}
	os.Exit(exitCode)
}
