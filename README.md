# Code Execution MCP Server

A Model Context Protocol (MCP) server written in Go that provides code execution capabilities for Bash/Zsh scripts, Python code, and Go code. This server enables LLMs to generate and execute code dynamically.

## Features

### Tools

1. **`execute_bash_script`** - Execute Bash/Zsh shell scripts
   - Best for: File operations, system commands, text processing, automation
   - Supports: Arguments, working directory, timeout configuration

2. **`execute_python_script`** - Execute Python 3 code
   - Best for: Data processing, API interactions, machine learning, complex algorithms
   - Supports: Command line arguments (sys.argv), working directory, timeout

3. **`execute_golang_code`** - Execute Go code
   - Best for: High-performance computing, concurrent operations, type-safe code
   - Requires: Complete Go program with `package main` and `func main()`

### Prompts

- **`code_executor`** - An intelligent prompt that helps LLMs choose the right tool based on the task description. Includes a decision framework and detailed documentation for each tool.

## Architecture

For a detailed explanation of the project's internal structure (Hexagonal Architecture), see [ARCHITECTURE.md](ARCHITECTURE.md).

## Download Pre-built Binaries

[![Build Status](https://github.com/tektwister/code_execution_mcp/actions/workflows/build.yml/badge.svg)](https://github.com/tektwister/code_execution_mcp/actions/workflows/build.yml)

**No need to build from source!** Download the latest pre-built binaries for your platform:

### [📥 Download Latest Release](https://github.com/tektwister/code_execution_mcp/releases/latest)

**Available platforms:**
- 🪟 **Windows** (amd64) - `code_execution_mcp-windows-amd64.exe`
- 🐧 **Linux** (amd64) - `code_execution_mcp-linux-amd64`
- 🐧 **Linux** (ARM64) - `code_execution_mcp-linux-arm64`
- 🍎 **macOS** (Intel) - `code_execution_mcp-darwin-amd64`
- 🍎 **macOS** (Apple Silicon) - `code_execution_mcp-darwin-arm64`

**Automated builds:** Binaries are automatically built and released on every commit to the main branch.

## Installation

### Prerequisites

- Go 1.23 or later
- For tool execution:
  - Bash (or Git Bash on Windows)
  - Python 3
  - Go runtime

### Build from Source

#### Using Makefile (Recommended)

```bash
# Download dependencies and build
make deps
make build
```

#### Manual Build

```bash
# Clone the repository
cd code_execution_mcp

# Download dependencies
go mod tidy

# Build the server
go build -o code-execution-mcp ./cmd/server
```

For more build options, see `make help` or check the [Makefile](Makefile).

## Usage

### Running the Server

The MCP server uses stdio transport, which means it communicates via standard input/output:

```bash
./code-execution-mcp
```

### Configuration with Claude Desktop

Add this to your Claude Desktop configuration file:

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "code-execution": {
      "command": "/path/to/code-execution-mcp"
    }
  }
}
```

### Tool Examples

#### Execute Bash Script

```json
{
  "tool": "execute_bash_script",
  "arguments": {
    "script": "echo 'Hello World' && ls -la",
    "timeout": 30
  }
}
```

#### Execute Python Script

```json
{
  "tool": "execute_python_script",
  "arguments": {
    "code": "import json\ndata = {'message': 'Hello from Python'}\nprint(json.dumps(data, indent=2))",
    "timeout": 30
  }
}
```

#### Execute Go Code

```json
{
  "tool": "execute_golang_code",
  "arguments": {
    "code": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello from Go!\")\n}",
    "timeout": 60
  }
}
```

## Security Considerations

⚠️ **Warning**: This MCP server executes arbitrary code on the host machine. Consider the following:

1. **Sandboxing**: Consider running in a container or VM for production use
2. **Timeouts**: All executions have configurable timeouts (max 300 seconds)
3. **Access Control**: Limit who can connect to this MCP server
4. **Code Review**: LLMs may generate code that has unintended side effects

## Output Format

Each execution returns:

- **Exit Code**: 0 for success, non-zero for failure
- **Duration**: Time taken for execution
- **Standard Output**: Program output
- **Standard Error**: Error messages (if any)

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
