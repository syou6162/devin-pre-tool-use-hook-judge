package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/constants"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Printf("devin-pre-tool-use-hook-judge %s\n", constants.Version)
		return nil
	}

	fmt.Println("hello, world")
	return nil
}
