package domain

import "time"

// ExecutionRequest represents a request to execute code
type ExecutionRequest struct {
	Language   string
	Code       string
	Script     string
	Args       []string
	WorkingDir string
	Timeout    int
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	ExitCode  int
	Stdout    string
	Stderr    string
	Duration  time.Duration
	IsError   bool
	ErrorType ExecutionErrorType
}

// ExecutionErrorType categorizes execution errors
type ExecutionErrorType int

const (
	NoError ExecutionErrorType = iota
	ValidationError
	TimeoutError
	RuntimeError
	SystemError
)

// String returns the string representation of ExecutionErrorType
func (e ExecutionErrorType) String() string {
	switch e {
	case NoError:
		return "NoError"
	case ValidationError:
		return "ValidationError"
	case TimeoutError:
		return "TimeoutError"
	case RuntimeError:
		return "RuntimeError"
	case SystemError:
		return "SystemError"
	default:
		return "UnknownError"
	}
}
