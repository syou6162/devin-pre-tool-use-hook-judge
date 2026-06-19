package main

import (
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
