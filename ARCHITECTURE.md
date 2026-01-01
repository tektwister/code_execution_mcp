# Architecture

The `code_execution_mcp` project follows the **Hexagonal Architecture** (also known as Ports and Adapters) pattern. This ensures a clean separation between the core domain logic, external interfaces (like the MCP protocol), and infrastructure (like the actual code execution mechanisms).

## High-Level Overview

The architecture allows the core business logic to remain independent of external technologies. This makes the application easier to test, maintain, and extend.

## Layers

The project is organized into three main layers:

### 1. Core (`internal/core`)
This is the heart of the application. It contains the business logic and domain entities. It has **no external dependencies**.

- **Domain (`internal/core/domain`)**: Defines pure data structures like `ExecutionRequest` and `ExecutionResult`.
- **Ports (`internal/core/ports`)**: Defines the interfaces (contracts) that the Core uses to interact with the outside world. For example, the `CodeExecutor` interface defines how code should be executed, without specifying *how* it is done.

### 2. Adapters (`internal/adapters`)
This layer connects the Core to specific technologies.

- **Primary Adapter (Driving)**: **MCP** (`internal/adapters/mcp`)
    - Handles incoming requests from the Model Context Protocol.
    - Converts MCP requests into domain objects.
    - Calls the Core Ports to perform actions.

- **Secondary Adapter (Driven)**: **Executors** (`internal/adapters/executor`)
    - Implements the interfaces defined in the Ports layer.
    - **ShellExecutor**: Executes Bash/Zsh scripts.
    - **PythonExecutor**: Executes Python code.
    - **GolangExecutor**: Executes Go code.

### 3. Wiring (`cmd/server`)
The `main.go` file acts as the **Composition Root**. It is responsible for:
1.  Initializing the specific Adapters (Executors).
2.  Injecting them into the Primary Adapter (MCP Handler).
3.  Starting the server.

## Data Flow

1.  **Request**: An MCP client sends a `execute_command` request.
2.  **MCP Adapter**: The `ToolHandler` receives the request and converts it into a `domain.ExecutionRequest`.
3.  **Core Interface**: The handler calls the `Execute` method on the injected `CodeExecutor`.
4.  **Executor Adapter**: The specific implementation (e.g., `ShellExecutor`) runs the actual command on the OS.
5.  **Response**: The result is wrapped in a `domain.ExecutionResult` and returned up the chain to the client.
