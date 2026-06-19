package judge

import (
	"fmt"
	"strings"
)

// extractJSON removes code fences and extracts the first JSON object from raw text.
func extractJSON(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", fmt.Errorf("empty response")
	}

	s = stripCodeFences(s)
	s = strings.TrimSpace(s)

	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return "", fmt.Errorf("no JSON object found in response")
	}

	return strings.TrimSpace(s[start : end+1]), nil
}

func stripCodeFences(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}

	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return strings.TrimPrefix(strings.TrimSuffix(s, "```"), "```")
	}

	content := strings.Join(lines[1:], "\n")
	if idx := strings.LastIndex(content, "```"); idx != -1 {
		content = content[:idx]
	}
	return content
}
