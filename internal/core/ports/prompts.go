package ports

import "github.com/modelcontextprotocol/go-sdk/mcp"

// PromptRegistry defines the interface for registering MCP prompts
type PromptRegistry interface {
	// RegisterPrompts registers all available prompts with the MCP server
	RegisterPrompts(server *mcp.Server)
}
