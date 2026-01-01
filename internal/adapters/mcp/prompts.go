package mcp

import (
	"context"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// PromptHandler implements the MCP prompt handler adapter
type PromptHandler struct{}

// NewPromptHandler creates a new prompt handler
func NewPromptHandler() *PromptHandler {
	return &PromptHandler{}
}

// RegisterPrompts registers all prompts with the MCP server
func (h *PromptHandler) RegisterPrompts(server *sdk.Server) {
	prompt := &sdk.Prompt{
		Name:        "code_executor",
		Description: "Helps you choose the right programming language and execute code based on your task. This prompt analyzes your requirements and suggests whether to use Bash/Zsh (for shell operations), Python (for data processing and scripting), or Go (for high-performance tasks).",
		Arguments: []*sdk.PromptArgument{
			{
				Name:        "task",
				Description: "The programming task or problem you want to solve",
				Required:    true,
			},
			{
				Name:        "preferences",
				Description: "Optional: Specify language preference (bash, python, go) or any specific requirements",
				Required:    false,
			},
		},
	}

	server.AddPrompt(prompt, h.handleCodeExecutorPrompt)
}

// handleCodeExecutorPrompt generates the prompt content
func (h *PromptHandler) handleCodeExecutorPrompt(ctx context.Context, req *sdk.GetPromptRequest) (*sdk.GetPromptResult, error) {
	task := ""
	preferences := ""

	if args := req.Params.Arguments; args != nil {
		if t, ok := args["task"]; ok {
			task = t
		}
		if p, ok := args["preferences"]; ok {
			preferences = p
		}
	}

	promptText := generateCodeExecutorPrompt(task, preferences)

	return &sdk.GetPromptResult{
		Description: "Code Executor - Intelligent Language Selection",
		Messages: []*sdk.PromptMessage{
			{
				Role:    "user",
				Content: &sdk.TextContent{Text: promptText},
			},
		},
	}, nil
}

// generateCodeExecutorPrompt creates the prompt text for code execution guidance
func generateCodeExecutorPrompt(task, preferences string) string {
	prompt := `# Code Execution Assistant

You are a helpful coding assistant with access to three code execution tools. Your job is to help the user accomplish their programming task by choosing the most appropriate language and writing executable code.

## Available Tools

### 1. execute_bash_script
**Best for:**
- File system operations (creating, moving, deleting files/directories)
- Text processing with sed, awk, grep
- System administration tasks
- Chaining multiple shell commands
- Environment variable operations
- Quick one-liners and automation scripts
- Package management commands
- Git operations

**Input parameters:**
- ` + "`script`" + ` (required): The bash script to execute
- ` + "`args`" + ` (optional): Command line arguments
- ` + "`working_dir`" + ` (optional): Working directory
- ` + "`timeout`" + ` (optional): Timeout in seconds (default: 30, max: 300)

### 2. execute_python_script
**Best for:**
- Data processing and analysis
- Mathematical computations
- Working with APIs and JSON/XML data
- Machine learning and data science tasks
- String manipulation and regex
- File parsing (CSV, JSON, YAML, etc.)
- Web scraping
- Complex algorithms
- Rapid prototyping

**Input parameters:**
- ` + "`code`" + ` (required): Python 3 code to execute
- ` + "`args`" + ` (optional): Arguments accessible via sys.argv
- ` + "`working_dir`" + ` (optional): Working directory
- ` + "`timeout`" + ` (optional): Timeout in seconds (default: 30, max: 300)

### 3. execute_golang_code
**Best for:**
- High-performance computing
- Concurrent/parallel processing
- Network programming
- System programming
- When you need type safety
- Memory-efficient operations
- Building and testing Go snippets

**Input parameters:**
- ` + "`code`" + ` (required): Go code with ` + "`package main`" + ` and ` + "`func main()`" + `
- ` + "`working_dir`" + ` (optional): Working directory
- ` + "`timeout`" + ` (optional): Timeout in seconds (default: 60, max: 300)

## Decision Framework

Use this decision tree to select the right tool:

1. **Is it a shell/system command or file operation?** → Use ` + "`execute_bash_script`" + `
2. **Does it involve data processing, APIs, or needs Python libraries?** → Use ` + "`execute_python_script`" + `
3. **Does it need high performance, concurrency, or type safety?** → Use ` + "`execute_golang_code`" + `
4. **Is it a simple script or automation?** → Use ` + "`execute_bash_script`" + ` or ` + "`execute_python_script`" + `

## User's Task

`

	if task != "" {
		prompt += "**Task:** " + task + "\n\n"
	} else {
		prompt += "**Task:** [Please describe what you want to accomplish]\n\n"
	}

	if preferences != "" {
		prompt += "**User Preferences:** " + preferences + "\n\n"
	}

	prompt += `## Your Response Should Include:

1. **Language Recommendation**: Which tool to use and why
2. **Code Implementation**: Complete, executable code that solves the task
3. **Execution**: Call the appropriate tool with the code
4. **Interpretation**: Explain the output and any next steps

Remember:
- Always provide complete, runnable code
- Handle potential errors gracefully
- If the task is ambiguous, ask clarifying questions before executing
- If the output indicates an error, help debug and provide a corrected solution
`

	return prompt
}
