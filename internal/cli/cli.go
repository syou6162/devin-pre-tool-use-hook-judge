package cli

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

// Run parses flags and executes the CLI.
func Run(args []string) error {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *showVersion {
		fmt.Printf("devin-pre-tool-use-hook-judge %s\n", version)
		return nil
	}

	fmt.Println("hello, world")
	return nil
}
