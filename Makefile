# Go Toolbox Makefile

# Variables

BINARY_NAME=toolbox
BUILD_DIR=bin
GO_FILES=$(shell find . -name "*.go" -type f)
# Find all directories under cmd/ containing a main.go (recursively)
MAIN_DIRS=$(shell find cmd -type f -name "main.go" -exec dirname {} \; | sort)

# Default target
.DEFAULT_GOAL := help

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: help build build-all clean test test-verbose test-coverage lint fmt vet run install deps tidy check-deps security update-deps

## help: Show this help message
help:
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## build: Build all applications

build: clean
	@echo "$(GREEN)Building all applications...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for dir in $(MAIN_DIRS); do \
		parent=$$(basename $$(dirname $$dir)); \
		app=$$(basename $$dir); \
		bin_name=$$parent-$$app; \
		echo "Building $$dir as $$bin_name..."; \
		go build -ldflags="-s -w" -o $(BUILD_DIR)/$$bin_name ./$$dir; \
	done
	@echo "$(GREEN)Build complete!$(NC)"

## build-all: Build for multiple platforms
build-all: clean
	@echo "$(GREEN)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for dir in $(MAIN_DIRS); do \
		app=$$(basename $$dir); \
		echo "Building $$app for multiple platforms..."; \
		GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$$app-windows-amd64.exe ./$$dir; \
		GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$$app-linux-amd64 ./$$dir; \
		GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$$app-darwin-amd64 ./$$dir; \
		GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$$app-darwin-arm64 ./$$dir; \
	done
	@echo "$(GREEN)Multi-platform build complete!$(NC)"

## clean: Remove build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	@go test ./...

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(GREEN)Running tests (verbose)...$(NC)"
	@go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linters
lint:
	@echo "$(GREEN)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found, install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w -local github.com/nate3d/toolbox .; \
	else \
		echo "$(YELLOW)goimports not found, install with: go install golang.org/x/tools/cmd/goimports@latest$(NC)"; \
	fi

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...

## run: Run the main CLI application (if exists)
run:
	@if [ -d "cmd/cli/main" ]; then \
		go run ./cmd/cli/main; \
	else \
		echo "$(RED)No main CLI application found in cmd/cli/main$(NC)"; \
	fi

## install: Install applications to GOPATH/bin
install:
	@echo "$(GREEN)Installing applications...$(NC)"
	@for dir in $(MAIN_DIRS); do \
		echo "Installing $$dir..."; \
		go install ./$$dir; \
	done

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@go mod download

## tidy: Tidy up dependencies
tidy:
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	@go mod tidy

## check-deps: Check for outdated dependencies
check-deps:
	@echo "$(GREEN)Checking for outdated dependencies...$(NC)"
	@go list -u -m all

## security: Run security checks
security:
	@echo "$(GREEN)Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)gosec not found, install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi

## update-deps: Update all dependencies
update-deps:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy


## dev-setup: Set up development environment (installs all required dev tools)
dev-setup:
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)Development tools installed!$(NC)"

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "$(GREEN)All checks passed!$(NC)"
