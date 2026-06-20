package main

import (
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/cmd"
)

func main() {
	os.Exit(cmd.NewRoot().Run(os.Args[1:]))
}
