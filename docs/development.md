# Development Guide

This guide covers development practices and patterns used in the Go Toolbox project.

## Project Architecture

### Directory Structure

- `cmd/`: Main applications (executable entry points)
  - `cli/`: Individual CLI tools
  - `tui/`: Individual TUI applications  
  - `unified/`: Alternative unified binary
  - `embedded/`: **Recommended unified binary** (maximum code reuse)
- `internal/`: Private application code (cannot be imported by other projects)
- `pkg/`: Public library code (can be imported by other projects)
- `configs/`: Configuration files
- `scripts/`: Build and deployment scripts
- `test/`: Additional test files and utilities
- `docs/`: Documentation

### Binary Strategy

The project supports multiple binary packaging strategies:

#### 1. Embedded Binary (Recommended) 🎯
- **Location**: `cmd/embedded/main.go`
- **Binary**: `bin/embedded`
- **Benefits**: Maximum code reuse, smallest size, easiest maintenance
- **Usage**: 
  ```bash
  ./bin/embedded --help     # CLI mode (default)
  ./bin/embedded tui        # TUI mode
  ./bin/embedded serve ./   # Server mode
  ```

#### 2. Individual Binaries
- **Binaries**: `bin/cli-main`, `bin/tui-main`, `bin/cli-serve`
- **Benefits**: Component isolation, debugging ease
- **Usage**: Development and testing

#### 3. Unified Binary (Alternative)
- **Location**: `cmd/unified/main.go`  
- **Binary**: `bin/cmd-unified`
- **Benefits**: Alternative implementation approach

### Dependency Management

We use Go modules for dependency management. Key principles:

1. **Minimal dependencies**: Only add dependencies when necessary
2. **Well-maintained packages**: Choose packages with active maintenance
3. **Standard library first**: Prefer standard library when possible
4. **Version pinning**: Pin to specific versions for stability

### Code Organization

#### Package Design

- Each package should have a single, well-defined purpose
- Use interfaces to define contracts between packages
- Keep public APIs minimal and stable
- Document all public functions and types

#### Error Handling

Follow Go's idiomatic error handling:

```go
// Good: Explicit error handling
result, err := someOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Bad: Ignoring errors
result, _ := someOperation()
```

#### Logging

Use structured logging throughout the application:

```go
// Good: Structured logging
logger.Info("Processing file", 
    "filename", filename,
    "size", fileSize,
    "duration", processingTime)

// Bad: Unstructured logging
logger.Info(fmt.Sprintf("Processing file %s with size %d", filename, fileSize))
```

## Development Workflow

### Setting Up Development Environment

1. **Install Go**: Version 1.21 or later
2. **Install development tools**:

   ```bash
   make dev-setup
   ```

3. **Configure your editor** with Go support
4. **Run initial build**:

   ```bash
   make build
   ```

### Adding New Features

1. **Plan the feature**: Consider interfaces and public APIs
2. **Write tests first**: Follow TDD when possible
3. **Implement the feature**: Keep functions small and focused
4. **Update documentation**: Document any new public APIs
5. **Add examples**: Provide usage examples

### Testing Strategy

#### Unit Tests

- Test all public functions
- Use table-driven tests for multiple test cases
- Include edge cases and error conditions
- Aim for high test coverage

```go
func TestStringReverse(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"empty string", "", ""},
        {"single char", "a", "a"},
        {"normal string", "hello", "olleh"},
        {"unicode", "🙂🙃", "🙃🙂"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := String().Reverse(tt.input)
            if result != tt.expected {
                t.Errorf("Reverse(%q) = %q, want %q", tt.input, result, tt.expected)
            }
        })
    }
}
```

#### Integration Tests

- Test complete workflows
- Use real external dependencies when possible
- Provide fallbacks for unavailable services

#### Benchmarks

- Benchmark performance-critical code
- Use `go test -bench=.` to run benchmarks
- Include benchmarks in CI for regression detection

### Code Quality

#### Formatting

- Use `gofmt` for code formatting
- Run `go vet` to catch common issues
- Use `golangci-lint` for comprehensive linting

#### Documentation

- Document all public APIs with godoc comments
- Include examples in documentation
- Keep README files up to date

#### Performance

- Profile applications using `go tool pprof`
- Optimize hot paths identified by profiling
- Consider memory allocations in performance-critical code

### Adding CLI Commands

To add a new CLI command:

1. **Create command function**:

   ```go
   func createMyCommand() *cobra.Command {
       baseCmd := cli.NewBaseCommand("mycommand", "Description of my command")
       
       baseCmd.RunE = func(cmd *cobra.Command, args []string) error {
           return runMyCommand(baseCmd, args)
       }
       
       return baseCmd.Command
   }
   ```

2. **Implement command logic**:

   ```go
   func runMyCommand(cmd *cli.BaseCommand, args []string) error {
       cmd.PrintHeader("My Command")
       
       // Implementation here
       
       cmd.PrintSuccess("Command completed successfully")
       return nil
   }
   ```

3. **Add to root command**:

   ```go
   rootCmd.AddCommand(createMyCommand())
   ```

### Adding TUI Screens

To add a new TUI screen:

1. **Define model struct**:

   ```go
   type myScreenModel struct {
       // Screen state
   }
   ```

2. **Implement Bubble Tea interface**:

   ```go
   func (m myScreenModel) Init() tea.Cmd {
       return nil
   }
   
   func (m myScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       // Handle events
       return m, nil
   }
   
   func (m myScreenModel) View() string {
       // Render UI
       return "Screen content"
   }
   ```

3. **Add navigation**:

   ```go
   func NewMyScreenModel() tea.Model {
       return myScreenModel{}
   }
   ```

### Working with Configuration

#### Adding New Config Options

1. **Update config struct**:

   ```go
   type Config struct {
       // Existing fields...
       MyNewOption string `mapstructure:"my_new_option"`
   }
   ```

2. **Set default value**:

   ```go
   func setDefaults() {
       // Existing defaults...
       viper.SetDefault("my_new_option", "default_value")
   }
   ```

3. **Update config file**:

   ```yaml
   # configs/config.yaml
   my_new_option: "value"
   ```

#### Using Configuration

```go
// Get typed value
config := config.Get()
value := config.MyNewOption

// Get directly from viper
value := config.GetString("my_new_option")
```

### Building and Deployment

#### Local Development

```bash
# Run without building
go run ./cmd/cli/main

# Build for current platform
make build

# Run tests
make test

# Run all checks
make check
```

#### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/toolbox-linux ./cmd/cli/main
```

#### Release Process

1. **Update version** in main.go files
2. **Update CHANGELOG.md**
3. **Create git tag**: `git tag v1.0.0`
4. **Build release binaries**: `make build-all`
5. **Create GitHub release** with binaries

### Performance Guidelines

#### Memory Management

- Use `sync.Pool` for frequent allocations
- Avoid unnecessary string concatenations in loops
- Use slices efficiently (preallocate when size is known)

#### Concurrency

- Use goroutines for I/O-bound operations
- Use channels for communication between goroutines
- Consider using `context.Context` for cancellation
- Use `sync.WaitGroup` for waiting on multiple goroutines

#### File Operations

- Use buffered I/O for large files
- Close files and resources properly
- Consider memory-mapped files for large read-only files

### Debugging

- TBD

#### Logging Best Practices

- Use appropriate log levels
- Include context in log messages
- Avoid logging sensitive information

#### Debugging Tools

- Use `dlv` (Delve) for debugging
- Use `go tool trace` for performance analysis
- Use `go tool pprof` for profiling

### Contributing Guidelines

1. **Fork the repository**
2. **Create feature branch**: `git checkout -b feature/my-feature`
3. **Make changes** following the coding standards
4. **Add tests** for new functionality
5. **Update documentation** as needed
6. **Run all checks**: `make check`
7. **Commit changes** with descriptive messages
8. **Push to your fork**
9. **Create pull request**

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] All tests pass
- [ ] Code is properly documented
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed
- [ ] Error handling is appropriate
- [ ] Logging is structured and appropriate
- [ ] Dependencies are justified

## Common Patterns

### Error Wrapping

```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}
```

### Context Usage

```go
func processFile(ctx context.Context, filename string) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process file
    }
}
```

### Interface Design

```go
type FileProcessor interface {
    Process(ctx context.Context, filename string) error
}

type TextProcessor struct{}

func (t TextProcessor) Process(ctx context.Context, filename string) error {
    // Implementation
}
```

This guide should help you understand the development practices and get started contributing to the Go Toolbox project.
