package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
)

//go:embed builtin_configs/*.yaml
var builtinConfigs embed.FS

func main() {
	configPath := flag.String("config", "", "path to YAML configuration file")
	builtinName := flag.String("builtin", "", "builtin configuration name (without .yaml)")
	flag.Parse()

	if *configPath != "" && *builtinName != "" {
		fmt.Fprintln(os.Stderr, "cannot specify both --config and --builtin")
		os.Exit(2)
	}

	switch {
	case *configPath != "":
		if _, err := config.LoadFromFile(*configPath); err != nil {
			fmt.Fprintf(os.Stderr, "config error: %v\n", err)
			os.Exit(2)
		}
	case *builtinName != "":
		if _, err := config.LoadBuiltin(*builtinName, builtinConfigs); err != nil {
			fmt.Fprintf(os.Stderr, "config error: %v\n", err)
			os.Exit(2)
		}
	default:
		fmt.Println("devin-pre-tool-use-hook-judge")
	}
}
