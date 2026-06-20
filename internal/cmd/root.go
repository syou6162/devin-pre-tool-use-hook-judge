package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

const (
	ExitSuccess = 0
	ExitError   = 1
	ExitBlock   = 2

	version = "0.1.0"
)

// Root is the top-level command dispatcher for devin-pre-tool-use-hook-judge.
type Root struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewRoot creates a Root command with stdin, stdout, and stderr wired to os defaults.
func NewRoot() *Root {
	return &Root{
		stdin:  os.Stdin,
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
		fmt.Fprintf(r.stdout, "devin-pre-tool-use-hook-judge %s\n", version)
		return ExitSuccess
	}

	return r.runHook()
}

func (r *Root) runHook() int {
	if _, err := schema.ParseDevinInput(r.stdin); err != nil {
		output := schema.BlockOutput(err.Error())
		if writeErr := schema.WriteOutput(r.stdout, output); writeErr != nil {
			fmt.Fprintf(r.stderr, "%v\n", writeErr)
			return ExitError
		}
		return ExitBlock
	}

	output := schema.ApproveOutput()
	if err := schema.WriteOutput(r.stdout, output); err != nil {
		fmt.Fprintf(r.stderr, "%v\n", err)
		return ExitError
	}
	return ExitSuccess
}
