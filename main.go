package main

import (
	"fmt"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/cli"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/version"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(version.String())
		return
	}

	fmt.Println("devin-pre-tool-use-hook-judge")
}
