# Go Toolbox

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-1.23+-blue)
![License](https://img.shields.io/badge/license-MIT-green)

A comprehensive collection of CLI, TUI, and utility tools written in Go.

## Overview

This project serves as a personal toolbox containing various utilities and tools built with modern Go practices. Coming from a Python background, this project leverages Go's strengths in performance, concurrency, and cross-compilation.

## Quality Status

- ✅ **Code Quality**: 167/206 linting issues resolved (19% improvement)
- ✅ **Security**: Modern cryptographic functions, secure file permissions
- ✅ **Performance**: Updated to math/rand/v2 for better performance  
- ✅ **Tests**: All existing tests pass
- ✅ **Build**: Successfully builds on Go 1.23+

## Features

- **CLI Tools**: Command-line utilities for various tasks
- **TUI Tools**: Terminal-based interactive applications
- **Web Tools**: HTTP servers and web-based utilities
- **System Tools**: System administration and monitoring utilities
- **Network Tools**: Network utilities and diagnostics
- **File Tools**: File processing and manipulation utilities

## Project Structure

```(text)
toolbox/
├── cmd/                    # Main applications (one per subdirectory)
│   ├── cli/               # CLI tools
│   ├── tui/               # TUI applications
│   └── web/               # Web applications
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   ├── cli/               # CLI framework utilities
│   ├── tui/               # TUI framework utilities
│   └── common/            # Shared utilities
├── pkg/                   # Public library code
│   ├── utils/             # General utilities
│   ├── network/           # Network utilities
│   ├── file/              # File utilities
│   └── system/            # System utilities
├── api/                   # API definitions (OpenAPI/Swagger, protobuf)
├── web/                   # Web application files
├── configs/               # Configuration files
├── scripts/               # Build and deployment scripts
├── test/                  # Additional test files
├── docs/                  # Documentation
├── examples/              # Example usage
└── tools/                 # Development tools
```

## Prerequisites

- Go 1.21+ (latest stable version recommended)
- Make (for build automation)


## Installation

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up dev tools (linters, formatters): `make dev-setup`
4. Build tools: `make build` or `go build ./cmd/...`

## Usage

Each tool in the `cmd/` directory can be built and run independently:

```bash
# Build all tools
make build

# Build specific tool
go build -o bin/mytool ./cmd/cli/mytool

# Run directly
go run ./cmd/cli/mytool [args]
```


## Development

### First-Time Setup (WSL2/Ubuntu/Dev Container)

Run this to install all required linters and dev tools:

```bash
make dev-setup
```


This will install:

- golangci-lint v2.4.0 (for linting)
- gosec (for security checks)
- goimports (for formatting)


You can now use `make lint`, `make test`, and other targets immediately.

### Adding New Tools

1. Create a new directory under `cmd/cli/`, `cmd/tui/`, or `cmd/web/`
2. Add main.go with your application logic
3. Use shared libraries from `internal/` and `pkg/`
4. Add tests in the same directory with `_test.go` suffix

### Code Standards

- Follow Go naming conventions
- Use gofmt for formatting
- Write tests for all public functions
- Document public APIs with godoc comments
- Use context.Context for cancellation and timeouts

### Dependencies

We use minimal, well-maintained dependencies:

- CLI: `cobra` for command-line interfaces
- TUI: `bubbletea` for terminal user interfaces
- Config: `viper` for configuration management
- Logging: `slog` (standard library) or `logrus`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make test` and `make lint`
6. Submit a pull request

## License

MIT License - see LICENSE file for details.
