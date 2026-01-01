package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/aravi/code_execution_mcp/internal/core/domain"
	"github.com/aravi/code_execution_mcp/internal/core/ports"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ToolHandler implements the MCP tool handler adapter
type ToolHandler struct {
	shellExecutor  ports.CodeExecutor
	pythonExecutor ports.CodeExecutor
	goExecutor     ports.CodeExecutor
}

// NewToolHandler creates a new tool handler with the given executors
func NewToolHandler(shellExec, pythonExec, goExec ports.CodeExecutor) *ToolHandler {
	return &ToolHandler{
		shellExecutor:  shellExec,
		pythonExecutor: pythonExec,
		goExecutor:     goExec,
	}
}

// BashInput represents input for bash/zsh script execution
type BashInput struct {
	Script     string   `json:"script"`
	Args       []string `json:"args,omitempty"`
	WorkingDir string   `json:"working_dir,omitempty"`
	Timeout    int      `json:"timeout,omitempty"`
}

// PythonInput represents input for Python script execution
type PythonInput struct {
	Code       string   `json:"code"`
	Args       []string `json:"args,omitempty"`
	WorkingDir string   `json:"working_dir,omitempty"`
	Timeout    int      `json:"timeout,omitempty"`
}

// GolangInput represents input for Go code execution
type GolangInput struct {
	Code       string `json:"code"`
	WorkingDir string `json:"working_dir,omitempty"`
	Timeout    int    `json:"timeout,omitempty"`
}

// RegisterTools registers all execution tools with the MCP server
func (h *ToolHandler) RegisterTools(server *sdk.Server) {
	// Tool 1: Execute Bash/Zsh Script
	sdk.AddTool[BashInput, any](server, &sdk.Tool{
		Name:        "execute_bash_script",
		Description: "Execute a bash or zsh shell script. Use this for shell commands, file operations, system administration tasks, or when you need to chain multiple shell commands together. Works on Unix-like systems (Linux, macOS) and Windows with Git Bash or WSL.",
	}, h.executeBashScript)

	// Tool 2: Execute Python Script
	sdk.AddTool[PythonInput, any](server, &sdk.Tool{
		Name:        "execute_python_script",
		Description: "Execute Python code. Ideal for data processing, mathematical computations, machine learning tasks, API interactions, and any task that benefits from Python's extensive library ecosystem. Requires Python 3 to be installed.",
	}, h.executePythonScript)

	// Tool 3: Execute Go Code
	sdk.AddTool[GolangInput, any](server, &sdk.Tool{
		Name:        "execute_golang_code",
		Description: "Execute Go (Golang) code. Best for high-performance tasks, concurrent operations, system programming, and when you need type safety and compiled performance. The code must include 'package main' and 'func main()'. Requires Go to be installed.",
	}, h.executeGolangCode)
}

// executeBashScript handles bash/zsh script execution
func (h *ToolHandler) executeBashScript(ctx context.Context, _ *sdk.CallToolRequest, input BashInput) (*sdk.CallToolResult, any, error) {
	req := domain.ExecutionRequest{
		Language:   "bash",
		Script:     input.Script,
		Args:       input.Args,
		WorkingDir: input.WorkingDir,
		Timeout:    input.Timeout,
	}

	result, err := h.shellExecutor.Execute(ctx, req)
	if err != nil {
		return &sdk.CallToolResult{
			IsError: true,
			Content: []sdk.Content{
				&sdk.TextContent{Text: fmt.Sprintf("Error executing bash script: %v", err)},
			},
		}, nil, nil
	}

	return formatResult(result, "Bash"), nil, nil
}

// executePythonScript handles Python code execution
func (h *ToolHandler) executePythonScript(ctx context.Context, _ *sdk.CallToolRequest, input PythonInput) (*sdk.CallToolResult, any, error) {
	req := domain.ExecutionRequest{
		Language:   "python",
		Code:       input.Code,
		Args:       input.Args,
		WorkingDir: input.WorkingDir,
		Timeout:    input.Timeout,
	}

	result, err := h.pythonExecutor.Execute(ctx, req)
	if err != nil {
		return &sdk.CallToolResult{
			IsError: true,
			Content: []sdk.Content{
				&sdk.TextContent{Text: fmt.Sprintf("Error executing Python code: %v", err)},
			},
		}, nil, nil
	}

	return formatResult(result, "Python"), nil, nil
}

// executeGolangCode handles Go code execution
func (h *ToolHandler) executeGolangCode(ctx context.Context, _ *sdk.CallToolRequest, input GolangInput) (*sdk.CallToolResult, any, error) {
	req := domain.ExecutionRequest{
		Language:   "go",
		Code:       input.Code,
		WorkingDir: input.WorkingDir,
		Timeout:    input.Timeout,
	}

	result, err := h.goExecutor.Execute(ctx, req)
	if err != nil {
		return &sdk.CallToolResult{
			IsError: true,
			Content: []sdk.Content{
				&sdk.TextContent{Text: fmt.Sprintf("Error executing Go code: %v", err)},
			},
		}, nil, nil
	}

	return formatResult(result, "Go"), nil, nil
}

// formatResult formats the execution result for MCP response
func formatResult(result *domain.ExecutionResult, language string) *sdk.CallToolResult {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("## %s Execution Result\n\n", language))
	summary.WriteString(fmt.Sprintf("**Exit Code:** %d\n", result.ExitCode))
	summary.WriteString(fmt.Sprintf("**Duration:** %s\n\n", result.Duration.String()))

	if result.Stdout != "" {
		summary.WriteString("### Standard Output\n```\n")
		summary.WriteString(result.Stdout)
		if !strings.HasSuffix(result.Stdout, "\n") {
			summary.WriteString("\n")
		}
		summary.WriteString("```\n\n")
	}

	if result.Stderr != "" {
		summary.WriteString("### Standard Error\n```\n")
		summary.WriteString(result.Stderr)
		if !strings.HasSuffix(result.Stderr, "\n") {
			summary.WriteString("\n")
		}
		summary.WriteString("```\n")
	}

	return &sdk.CallToolResult{
		IsError: result.IsError,
		Content: []sdk.Content{
			&sdk.TextContent{Text: summary.String()},
		},
	}
}
