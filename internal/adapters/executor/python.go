package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aravi/code_execution_mcp/internal/core/domain"
	"github.com/aravi/code_execution_mcp/internal/core/ports"
)

// PythonExecutor implements CodeExecutor for Python code
type PythonExecutor struct{}

// NewPythonExecutor creates a new Python executor
func NewPythonExecutor() ports.CodeExecutor {
	return &PythonExecutor{}
}

// Supports checks if this executor supports the given language
func (e *PythonExecutor) Supports(language string) bool {
	return language == "python" || language == "python3"
}

// Execute runs Python code
func (e *PythonExecutor) Execute(ctx context.Context, req domain.ExecutionRequest) (*domain.ExecutionResult, error) {
	if strings.TrimSpace(req.Code) == "" {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.ValidationError,
			Stderr:    "Python code cannot be empty",
		}, nil
	}

	timeout := getTimeout(req.Timeout, 30)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create a temporary file for the Python script
	tmpFile, err := os.CreateTemp("", "mcp_python_*.py")
	if err != nil {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.SystemError,
			Stderr:    fmt.Sprintf("Error creating temp file: %v", err),
		}, nil
	}
	defer os.Remove(tmpFile.Name())

	// Write the Python code to the temp file
	if _, err := tmpFile.WriteString(req.Code); err != nil {
		tmpFile.Close()
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.SystemError,
			Stderr:    fmt.Sprintf("Error writing to temp file: %v", err),
		}, nil
	}
	tmpFile.Close()

	// Build command arguments
	args := []string{tmpFile.Name()}
	args = append(args, req.Args...)

	// Try python3 first, then python
	pythonCmd := "python3"
	if runtime.GOOS == "windows" {
		pythonCmd = "python"
	}

	cmd := exec.CommandContext(ctx, pythonCmd, args...)
	if req.WorkingDir != "" {
		cmd.Dir = req.WorkingDir
	}

	return executeCommand(cmd)
}
