package cli

import "flag"

// Options holds CLI flags for the hook judge binary.
type Options struct {
	Config  string
	Builtin string
	Version bool
}

// Parse parses command-line arguments into Options.
func Parse(args []string) (*Options, error) {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	opts := &Options{}
	fs.BoolVar(&opts.Version, "version", false, "print version and exit")
	fs.StringVar(&opts.Config, "config", "", "path to config YAML file")
	fs.StringVar(&opts.Builtin, "builtin", "", "builtin validator name")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return opts, nil
}
