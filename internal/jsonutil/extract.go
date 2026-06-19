package jsonutil

import (
	"fmt"
	"strings"
)

// ExtractJSON finds and returns the first balanced JSON object or array in text.
// It strips markdown code fences when present.
func ExtractJSON(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("empty input")
	}

	text = stripCodeFences(text)
	text = strings.TrimSpace(text)

	startObj := strings.Index(text, "{")
	startArr := strings.Index(text, "[")
	start := -1
	isObject := true

	switch {
	case startObj >= 0 && (startArr < 0 || startObj < startArr):
		start = startObj
		isObject = true
	case startArr >= 0:
		start = startArr
		isObject = false
	default:
		return "", fmt.Errorf("no JSON object or array found")
	}

	end, err := findMatchingEnd(text, start, isObject)
	if err != nil {
		return "", err
	}

	return text[start : end+1], nil
}

func stripCodeFences(text string) string {
	if !strings.Contains(text, "```") {
		return text
	}

	lines := strings.Split(text, "\n")
	var cleaned []string
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence || !strings.HasPrefix(trimmed, "```") {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}

func findMatchingEnd(text string, start int, isObject bool) (int, error) {
	open := byte('{')
	close := byte('}')
	if !isObject {
		open = '['
		close = ']'
	}

	depth := 0
	inString := false
	escaped := false

	for i := start; i < len(text); i++ {
		ch := text[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				return i, nil
			}
			if depth < 0 {
				return 0, fmt.Errorf("unbalanced JSON")
			}
		}
	}

	return 0, fmt.Errorf("unbalanced JSON")
}
