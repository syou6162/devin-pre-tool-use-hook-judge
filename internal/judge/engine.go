package judge

import (
	"context"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

// Engine evaluates tool usage requests and returns a judgment result.
type Engine interface {
	Judge(ctx context.Context, input *schema.JudgeInput, cfg *config.Config) (*schema.JudgeResult, error)
}
