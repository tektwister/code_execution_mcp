# Code Execution MCP Server - Makefile

# Binary name
BINARY_NAME=code_execution_mcp

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Detect OS
ifeq ($(OS),Windows_NT)
	BINARY_EXT=.exe
	RM=del /Q
	RMDIR=rmdir /S /Q
else
	BINARY_EXT=
	RM=rm -f
	RMDIR=rm -rf
endif

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)$(BINARY_EXT) ./cmd/server

# Build without optimizations (for debugging)
.PHONY: build-debug
build-debug:
	$(GOBUILD) -o $(BINARY_NAME)$(BINARY_EXT) ./cmd/server

# Run the application
.PHONY: run
run:
	$(GOCMD) run ./cmd/server

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	$(RM) $(BINARY_NAME)$(BINARY_EXT)
	$(RM) $(BINARY_NAME)-*

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
.PHONY: update-deps
update-deps:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Cross-compilation targets
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 ./cmd/server

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 ./cmd/server

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe ./cmd/server

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 ./cmd/server

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 ./cmd/server

# Build all platforms
.PHONY: build-all
build-all: build-linux build-linux-arm64 build-windows build-darwin build-darwin-arm64

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Vet code
.PHONY: vet
vet:
	$(GOCMD) vet ./...

# Install the binary to GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install .

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary for current OS"
	@echo "  make build-debug   - Build without optimizations"
	@echo "  make run           - Run the application with go run"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make test          - Run tests"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make update-deps   - Update dependencies to latest"
	@echo "  make build-linux   - Cross-compile for Linux (amd64)"
	@echo "  make build-windows - Cross-compile for Windows (amd64)"
	@echo "  make build-darwin  - Cross-compile for macOS (amd64)"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make install       - Install binary to GOPATH/bin"
	@echo "  make help          - Show this help message"

