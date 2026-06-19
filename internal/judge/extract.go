package judge

import (
	"regexp"
	"strings"
)

var codeFencePattern = regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")

func ExtractJSON(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ErrEmptyResponse
	}

	if match := codeFencePattern.FindStringSubmatch(trimmed); len(match) == 2 {
		return strings.TrimSpace(match[1]), nil
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end > start {
		return trimmed[start : end+1], nil
	}

	return "", ErrInvalidJSON
}
