# Go Toolbox

[![CI](https://github.com/Velvet-Labs-LLC/go-toolbox/actions/workflows/ci.yml/badge.svg)](https://github.com/Velvet-Labs-LLC/go-toolbox/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-20.0%25-yellow)](https://github.com/Velvet-Labs-LLC/go-toolbox/actions)
[![Security](https://img.shields.io/badge/gosec-passing-brightgreen)](https://github.com/Velvet-Labs-LLC/go-toolbox/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Velvet-Labs-LLC/go-toolbox)](https://goreportcard.com/report/github.com/Velvet-Labs-LLC/go-toolbox)

A comprehensive collection of CLI, TUI, and utility tools written in Go with modern development practices.

## Overview

This project serves as a modern Go toolbox containing various utilities and applications built with Go 1.24+. It demonstrates enterprise-level CI/CD practices, comprehensive testing, security scanning, and cross-platform builds.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/Velvet-Labs-LLC/go-toolbox.git
cd go-toolbox

# Install dependencies
go mod download

# Build all tools
make build

# Run the main CLI tool
./bin/cli-main --help

# Start a file server (great for local network sharing)
./bin/cli-serve --dir ./docs --port 8080
```

## 📊 Project Status

| Metric | Status |
|--------|--------|
| **Build** | ![Passing](https://img.shields.io/badge/status-passing-brightgreen) |
| **Tests** | ![Coverage 20%](https://img.shields.io/badge/coverage-20.0%25-yellow) |
| **Security** | ![Gosec Clean](https://img.shields.io/badge/gosec-clean-brightgreen) |
| **Linting** | ![golangci-lint v2.4](https://img.shields.io/badge/golangci--lint-v2.4-blue) |
| **Go Version** | ![1.24+](https://img.shields.io/badge/go-1.24+-blue) |
| **Platforms** | ![Multi-platform](https://img.shields.io/badge/platforms-linux%20%7C%20windows%20%7C%20darwin-lightgrey) |

## 🛠️ Available Tools

### CLI Tools
- **`cli-main`** - Main CLI application with multiple utilities and sub-commands
- **`cli-serve`** - HTTP file server for local network sharing ([docs](./cmd/cli/serve/README.md))

### TUI Tools  
- **`tui-main`** - Interactive terminal application with tool generation capabilities

## 🏗️ Architecture

```
go-toolbox/
├── cmd/                    # Applications (builds to bin/)
│   ├── cli/               # Command-line tools
│   │   ├── main/          # → cli-main binary
│   │   └── serve/         # → cli-serve binary  
│   └── tui/               # Terminal UI applications
│       └── main/          # → tui-main binary
├── internal/              # Private application code
│   ├── cli/               # CLI framework and utilities
│   ├── config/            # Configuration management (Viper)
│   ├── generator/         # Code/tool generation utilities
│   └── logger/            # Structured logging (slog)
├── pkg/                   # Public library code
│   └── utils/             # General utilities (crypto-secure)
├── examples/              # Usage examples
├── docs/                  # Documentation
└── scripts/               # Build and development scripts
```

## 🔧 Development

### Prerequisites
- **Go 1.24+** (uses modern Go features)
- **Make** (for build automation)
- **Git** (for version control)

### First-Time Setup

```bash
# Install development tools (linters, formatters, security scanners)
make dev-setup

# Install Git pre-commit hooks (recommended for contributors)
make install-hooks

# Verify installation
make lint      # Run golangci-lint v2.4
make test      # Run all tests with coverage
make security  # Run gosec security scan
```

### Available Make Targets

```bash
make build        # Build all binaries (cross-platform)
make test         # Run tests with coverage
make lint         # Run golangci-lint
make security     # Run gosec security scan
make clean        # Clean build artifacts
make deps         # Download dependencies
make dev-setup    # Install development tools
make install-hooks # Install Git pre-commit hooks
```

### Git Hooks

The project includes a pre-commit hook that automatically runs before every commit:

- **`make fmt`** - Formats code using gofmt
- **`make lint`** - Runs golangci-lint to catch issues
- **`make test`** - Runs all tests to ensure functionality

Install with `make install-hooks`. The hook will:
- ✅ **Format** your code automatically
- ⚠️ **Block commits** if linting or tests fail
- 💡 **Skip hook**: `git commit --no-verify` (not recommended)

### Code Quality Standards

- ✅ **Security**: All crypto operations use `crypto/rand` (not `math/rand`)
- ✅ **Linting**: Passes golangci-lint v2.4 with strict configuration
- ✅ **Testing**: Comprehensive test suite with race detection
- ✅ **Documentation**: Godoc comments for all public APIs
- ✅ **Performance**: Benchmark tracking with regression detection

## 🧪 Testing & Quality

Our CI pipeline ensures high code quality:

```bash
# Run the full test suite
go test -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Security scan
gosec ./...

# Lint check
golangci-lint run
```

### Benchmark Tracking
The CI automatically tracks performance regressions:
- Runs benchmarks on every commit
- Compares against baseline performance
- Fails PR if significant regressions detected

## 🔐 Security Features

- **Crypto-secure randomness**: Uses `crypto/rand` for all random operations
- **Secure HTTP servers**: TLS 1.2+, proper timeouts, secure headers
- **File permissions**: Restrictive permissions (0750) for generated directories
- **Input validation**: Comprehensive path validation to prevent traversal attacks
- **Dependency scanning**: Automated security vulnerability detection

## 🌍 Cross-Platform Builds

The CI automatically builds for multiple platforms:

| OS | Architecture | Status |
|----|--------------|--------|
| Linux | amd64, arm64 | ✅ |
| Windows | amd64 | ✅ |
| macOS | amd64, arm64 | ✅ |

Download binaries from the [Releases page](https://github.com/Velvet-Labs-LLC/go-toolbox/releases).

## 📚 Key Dependencies

```go
// CLI Framework
github.com/spf13/cobra v1.10.1
github.com/spf13/viper v1.20.1

// TUI Framework  
github.com/charmbracelet/bubbletea v1.3.7
github.com/charmbracelet/lipgloss v1.1.0

// Utilities
github.com/fatih/color v1.18.0
github.com/olekukonko/tablewriter v1.0.9
github.com/schollz/progressbar/v3 v3.18.0
```

## 🤝 Contributing

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Make** your changes with tests
4. **Ensure** quality: `make lint && make test && make security`
5. **Commit** your changes: `git commit -m 'Add amazing feature'`
6. **Push** to branch: `git push origin feature/amazing-feature`
7. **Open** a Pull Request

### Development Workflow

```bash
# Start development
git checkout -b feature/my-feature

# Make changes, then verify quality
make lint test security

# Add benchmarks if performance-critical
go test -bench=. ./...

# Commit and push
git add .
git commit -m "feat: add amazing feature"
git push origin feature/my-feature
```

## 📈 Performance

Performance is tracked automatically in CI:
- **Benchmark suite** runs on every commit
- **Regression detection** prevents performance degradation  
- **Memory profiling** included in test coverage

Example benchmark results:
```
BenchmarkRandomString-8    1000000    1234 ns/op    64 B/op    2 allocs/op
BenchmarkHashSHA256-8       500000    2456 ns/op   128 B/op    3 allocs/op
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏷️ Version History

See [Releases](https://github.com/Velvet-Labs-LLC/go-toolbox/releases) for detailed changelog and version history.

---

**Built with ❤️ using Go 1.24 and modern development practices.**
