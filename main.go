package main

import (
	"fmt"
	"io"
	"os"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

func main() {
	os.Exit(Run(os.Stdin, os.Stdout))
}

// Run reads hook input from stdin, validates it, and writes the decision to stdout.
func Run(stdin io.Reader, stdout io.Writer) int {
	inputBytes, err := io.ReadAll(stdin)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("入力の読み込みに失敗しました: %v", err))
	}

	devinInput, err := schema.ParseDevinInput(inputBytes)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("入力検証エラー: %v", err))
	}

	judgeInput, err := schema.ToJudgeInput(devinInput)
	if err != nil {
		return writeBlock(stdout, fmt.Sprintf("入力変換エラー: %v", err))
	}

	// Judgment engine integration is implemented in a later issue.
	_ = judgeInput

	return writeApprove(stdout, "入力検証成功")
}

func writeBlock(stdout io.Writer, reason string) int {
	output := schema.BlockOutput(reason)
	if err := schema.WriteOutput(stdout, output); err != nil {
		fmt.Fprintf(os.Stderr, "出力エラー: %v\n", err)
		return schema.ExitCodeBlock
	}
	return schema.ExitCodeForDecision(output.Decision)
}

func writeApprove(stdout io.Writer, reason string) int {
	output := schema.ApproveOutput(reason)
	if err := schema.WriteOutput(stdout, output); err != nil {
		return writeBlock(stdout, fmt.Sprintf("出力エラー: %v", err))
	}
	return schema.ExitCodeForDecision(output.Decision)
}
