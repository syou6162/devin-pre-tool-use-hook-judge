package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/constants"
)

const (
	ExitSuccess = 0
	ExitError   = 1
)

// Root is the top-level command dispatcher for devin-pre-tool-use-hook-judge.
type Root struct {
	stdout io.Writer
	stderr io.Writer
}

// NewRoot creates a Root command with stdout and stderr wired to os.Stdout and os.Stderr.
func NewRoot() *Root {
	return &Root{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// Run parses flags and executes the CLI.
func (r *Root) Run(args []string) int {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	fs.SetOutput(r.stderr)

	showVersion := fs.Bool("version", false, "print version and exit")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(r.stderr, "%v\n", err)
		return ExitError
	}

	if *showVersion {
		fmt.Fprintf(r.stdout, "devin-pre-tool-use-hook-judge %s\n", constants.Version)
		return ExitSuccess
	}

	fmt.Fprintln(r.stdout, "hello, world")
	return ExitSuccess
}
