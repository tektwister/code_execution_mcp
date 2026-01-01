package main

import (
	"context"
	"log"

	"github.com/aravi/code_execution_mcp/internal/adapters/executor"
	mcpadapter "github.com/aravi/code_execution_mcp/internal/adapters/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Create MCP server with implementation info
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "code-execution-mcp",
		Version: "v1.0.0",
	}, nil)

	// Initialize executors (secondary/outbound adapters)
	shellExecutor := executor.NewShellExecutor()
	pythonExecutor := executor.NewPythonExecutor()
	goExecutor := executor.NewGolangExecutor()

	// Initialize MCP adapters (primary/inbound adapters) with dependencies
	toolHandler := mcpadapter.NewToolHandler(shellExecutor, pythonExecutor, goExecutor)
	promptHandler := mcpadapter.NewPromptHandler()

	// Register tools and prompts
	toolHandler.RegisterTools(server)
	promptHandler.RegisterPrompts(server)

	log.Println("Starting Code Execution MCP Server...")

	// Run the server over stdin/stdout (stdio transport)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
