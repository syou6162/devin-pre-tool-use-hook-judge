// Package exception defines custom error types for the hook judge.
package exception

// JudgeError is the base error for judgment failures.
type JudgeError struct {
	Message string
}

func (e *JudgeError) Error() string {
	return e.Message
}

// NoResponseError indicates no response from the judgment system.
type NoResponseError struct {
	JudgeError
}

// SchemaValidationError indicates schema validation failed.
type SchemaValidationError struct {
	JudgeError
}
