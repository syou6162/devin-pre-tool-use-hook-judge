package judge

import "errors"

var (
	ErrEmptyResponse = errors.New("empty response from judgment engine")
	ErrInvalidJSON   = errors.New("could not extract JSON from response")
)
