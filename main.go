package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/judge"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

type runner struct {
	stdin  io.Reader
	stdout io.Writer
	engine judge.Engine
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, judge.NewDevinEngine("devin")))
}

func run(args []string, stdin io.Reader, stdout io.Writer, engine judge.Engine) int {
	r := &runner{
		stdin:  stdin,
		stdout: stdout,
		engine: engine,
	}
	return r.execute(args)
}

func (r *runner) execute(args []string) int {
	fs := flag.NewFlagSet("devin-pre-tool-use-hook-judge", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var configPath string
	var builtinName string
	fs.StringVar(&configPath, "config", "", "Path to external YAML configuration file")
	fs.StringVar(&builtinName, "builtin", "", "Name of builtin configuration to use")

	if err := fs.Parse(args); err != nil {
		return r.writeBlock(fmt.Sprintf("failed to parse flags: %v", err))
	}

	cfg, err := r.loadConfig(configPath, builtinName)
	if err != nil {
		return r.writeBlock(err.Error())
	}

	raw, err := io.ReadAll(r.stdin)
	if err != nil {
		return r.writeBlock(fmt.Sprintf("failed to read stdin: %v", err))
	}

	devinInput, err := schema.ParseDevinInput(raw)
	if err != nil {
		return r.writeBlock(err.Error())
	}
	if err := schema.ValidateDevinInput(devinInput); err != nil {
		return r.writeBlock(err.Error())
	}

	judgeInput := schema.ToJudgeInput(devinInput)
	judgeOutput, err := r.engine.Judge(context.Background(), judgeInput, cfg)
	if err != nil {
		return r.writeBlock(fmt.Sprintf("judgment failed: %v", err))
	}

	output := schema.ToDevinOutput(judgeOutput)
	return r.writeOutput(output)
}

func (r *runner) loadConfig(configPath, builtinName string) (config.Config, error) {
	switch {
	case configPath != "" && builtinName != "":
		return config.Config{}, fmt.Errorf("--config and --builtin cannot be used together")
	case configPath != "":
		return config.LoadFile(configPath)
	case builtinName != "":
		return config.LoadBuiltin(builtinName)
	default:
		return config.Config{}, fmt.Errorf("no configuration specified; refusing operation for safety")
	}
}

func (r *runner) writeBlock(reason string) int {
	return r.writeOutput(schema.BlockOutput(reason))
}

func (r *runner) writeOutput(output schema.DevinOutput) int {
	encoded, err := json.Marshal(output)
	if err != nil {
		fallback := schema.BlockOutput(fmt.Sprintf("failed to encode output: %v", err))
		encoded, _ = json.Marshal(fallback)
		output = fallback
	}

	if _, err := fmt.Fprintln(r.stdout, string(encoded)); err != nil {
		return schema.ExitCodeForOutput(schema.BlockOutput("failed to write output"))
	}
	return schema.ExitCodeForOutput(output)
}
