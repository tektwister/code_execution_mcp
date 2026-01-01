package ports

import "github.com/modelcontextprotocol/go-sdk/mcp"

// ToolRegistry defines the interface for registering MCP tools
type ToolRegistry interface {
	// RegisterTools registers all available tools with the MCP server
	RegisterTools(server *mcp.Server)
}
