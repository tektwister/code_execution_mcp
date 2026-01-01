package ports

import (
	"context"

	"github.com/aravi/code_execution_mcp/internal/core/domain"
)

// CodeExecutor defines the interface for executing code
type CodeExecutor interface {
	// Execute runs code and returns the execution result
	Execute(ctx context.Context, req domain.ExecutionRequest) (*domain.ExecutionResult, error)
	
	// Supports checks if this executor supports the given language
	Supports(language string) bool
}
