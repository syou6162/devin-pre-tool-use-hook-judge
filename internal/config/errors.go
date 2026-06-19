package config

import "fmt"

// Error represents a configuration loading or validation failure.
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func configError(message string) *Error {
	return &Error{Message: message}
}

func configErrorf(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}
