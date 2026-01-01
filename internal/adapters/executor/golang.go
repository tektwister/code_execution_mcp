package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aravi/code_execution_mcp/internal/core/domain"
	"github.com/aravi/code_execution_mcp/internal/core/ports"
)

// GolangExecutor implements CodeExecutor for Go code
type GolangExecutor struct{}

// NewGolangExecutor creates a new Go executor
func NewGolangExecutor() ports.CodeExecutor {
	return &GolangExecutor{}
}

// Supports checks if this executor supports the given language
func (e *GolangExecutor) Supports(language string) bool {
	return language == "go" || language == "golang"
}

// Execute runs Go code
func (e *GolangExecutor) Execute(ctx context.Context, req domain.ExecutionRequest) (*domain.ExecutionResult, error) {
	if strings.TrimSpace(req.Code) == "" {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.ValidationError,
			Stderr:    "Go code cannot be empty",
		}, nil
	}

	// Validate that the code has required structure
	if !strings.Contains(req.Code, "package main") {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.ValidationError,
			Stderr:    "Go code must include 'package main'",
		}, nil
	}

	if !strings.Contains(req.Code, "func main()") {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.ValidationError,
			Stderr:    "Go code must include 'func main()'",
		}, nil
	}

	timeout := getTimeout(req.Timeout, 60)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create a temporary directory for the Go module
	tmpDir, err := os.MkdirTemp("", "mcp_golang_*")
	if err != nil {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.SystemError,
			Stderr:    fmt.Sprintf("Error creating temp directory: %v", err),
		}, nil
	}
	defer os.RemoveAll(tmpDir)

	// Write the Go file
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte(req.Code), 0644); err != nil {
		return &domain.ExecutionResult{
			IsError:   true,
			ErrorType: domain.SystemError,
			Stderr:    fmt.Sprintf("Error writing Go file: %v", err),
		}, nil
	}

	// Use 'go run' for execution
	cmd := exec.CommandContext(ctx, "go", "run", goFile)
	if req.WorkingDir != "" {
		cmd.Dir = req.WorkingDir
	} else {
		cmd.Dir = tmpDir
	}

	return executeCommand(cmd)
}
