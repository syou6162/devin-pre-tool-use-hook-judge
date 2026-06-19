package judge

import (
	"context"
	"errors"
	"fmt"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

const maxRetries = 3

// JudgeOptions configures a judgment request.
type JudgeOptions struct {
	CustomPrompt string
	Model        string
}

// JudgeEngine validates tool usage requests.
type JudgeEngine interface {
	Judge(ctx context.Context, input *schema.JudgeInput, opts JudgeOptions) (*schema.HookOutput, error)
}

// JudgeError indicates the judgment engine failed after retries.
type JudgeError struct {
	Attempts int
	Cause    error
}

func (e *JudgeError) Error() string {
	return fmt.Sprintf("judge failed after %d attempts: %v", e.Attempts, e.Cause)
}

func (e *JudgeError) Unwrap() error {
	return e.Cause
}

// BlockOnError returns a safe deny output when judgment fails.
func BlockOnError(err error) *schema.HookOutput {
	reason := "判定システムが正しいスキーマ形式で応答できませんでした。安全のため操作を拒否します。"
	var judgeErr *JudgeError
	if errors.As(err, &judgeErr) && judgeErr.Cause != nil {
		reason = fmt.Sprintf("%s (%v)", reason, judgeErr.Cause)
	}
	return schema.BlockHookOutput(reason)
}
