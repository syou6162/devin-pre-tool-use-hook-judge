package judge

import (
	"context"
	"fmt"

	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/config"
	"github.com/syou6162/devin-pre-tool-use-hook-judge/internal/schema"
)

// MockEngine is a test double for Engine.
type MockEngine struct {
	Result *schema.JudgeResult
	Err    error
	Calls  int
}

// Judge returns the configured result or error.
func (m *MockEngine) Judge(ctx context.Context, input *schema.JudgeInput, cfg *config.Config) (*schema.JudgeResult, error) {
	m.Calls++
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Result == nil {
		return nil, fmt.Errorf("mock engine has no result configured")
	}
	return m.Result, nil
}
