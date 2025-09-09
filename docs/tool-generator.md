# 🛠️ Tool Generator

The Go Toolbox includes a powerful tool generator that allows you to quickly create new CLI and TUI applications with pre-configured templates and best practices.

## Features

- **Interactive TUI Generator**: Navigate through an intuitive text-based interface
- **Multiple Tool Types**: Generate CLI and TUI applications
- **Template-Based**: Uses industry-standard templates with proper structure
- **VS Code Integration**: Quick-create tasks available in the command palette
- **Automatic Build Integration**: Generated tools are automatically recognized by the Makefile

## Usage Methods

### 1. TUI Generator (Recommended)

Run the TUI application and select "🛠️ Tool Generator":

```bash
./bin/tui-main
# or using the embedded binary
./bin/embedded tui
```

Then follow the interactive prompts to:

1. Select tool type (CLI or TUI)
2. Enter tool name (lowercase, no spaces)
3. Enter tool description
4. Generate the tool with all necessary files

### 2. VS Code Tasks (Quick Generation)

Use VS Code's Command Palette (`Ctrl+Shift+P`) and search for:

- **🛠️ Generate CLI Tool**: Launch the TUI generator
- **🚀 Quick New CLI Tool**: Create a CLI tool with prompts
- **🖥️ Quick New TUI Tool**: Create a TUI tool with prompts

### 3. Manual Generation

You can also create tools manually using the templates in `internal/generator/templates.go`.

## Generated Tool Structure

### CLI Tools

```(text)
cmd/cli/your-tool/
└── main.go          # CLI application with flag parsing
```

**Features**:
- Cobra command-line framework integration
- Flag parsing and validation
- Help text generation
- Integration with toolbox utilities

### TUI Tools

```(text)
cmd/tui/your-tool/
└── main.go          # TUI application with Bubble Tea
```

**Features**:
- Bubble Tea framework integration
- Interactive terminal interface
- Keyboard navigation
- Styled components using Lip Gloss

## Build Integration

Generated tools are automatically integrated into the build system:

- CLI tools: `bin/cli-your-tool`
- TUI tools: `bin/tui-your-tool`

Build with:

```bash
make build        # Builds all tools including generated ones
make build-all    # Cross-platform builds
```

## Tool Templates

All generated tools include:

- ✅ **Configuration Management**: Using the internal config package
- ✅ **Structured Logging**: Using the internal logger package
- ✅ **Error Handling**: Proper error handling patterns
- ✅ **Help Systems**: Built-in help and usage information
- ✅ **Build Integration**: Automatic Makefile recognition

### CLI Tool Features

- Flag parsing with `flag` package
- Version and help commands
- Verbose logging option
- Configuration and logger initialization

### TUI Tool Features

- Bubble Tea framework integration
- Keyboard navigation (j/k, arrows)
- Styled interface with Lipgloss
- Proper cleanup and quit handling

## Building Generated Tools

After generating a tool, build it using:

```bash
make build
```

Your new tool will be available in the `bin/` directory as:

- CLI tools: `bin/cli-your-tool`
- TUI tools: `bin/tui-your-tool`

## Development Workflow

1. **Generate**: Use the TUI generator or VS Code tasks
2. **Implement**: Add your business logic to the generated template
3. **Test**: Run `make test` to ensure quality
4. **Build**: Run `make build` to compile
5. **Deploy**: Your tool is ready for distribution

## Examples

### Generated CLI Tool Example

```bash
# Generate a file hash calculator
./bin/tui-main
# or
./bin/embedded tui
# Select: CLI Tool -> "filehasher" -> "Calculate file hashes"

# Use the generated tool
./bin/cli-filehasher --help
./bin/cli-filehasher /path/to/file
```

### Generated TUI Tool Example

```bash
# Generate a system monitor
./bin/tui-main
# or
./bin/embedded tui
# Select: TUI Tool -> "sysmonitor" -> "Monitor system resources"

# Use the generated tool
./bin/tui-sysmonitor
```

## Customization

### Modifying Templates

Edit the templates in `internal/generator/templates.go`:

- `cliTemplate`: CLI application template
- `tuiTemplate`: TUI application template

### Adding New Tool Types

To add new tool types:

1. Add a new `ToolType` constant
2. Create a new template
3. Update the generator UI choices
4. Add generation logic in `generateMainFile()`

## Best Practices

### Naming Conventions

- Use lowercase names with hyphens: `file-hasher`, `network-ping`
- Avoid spaces, special characters, or uppercase
- Keep names descriptive but concise

### Project Structure

- CLI tools go in `cmd/cli/`
- TUI tools go in `cmd/tui/`
- Shared code goes in `internal/` or `pkg/`

### Development Tips

- Use the built-in config and logger packages
- Follow Go naming conventions
- Add tests in `*_test.go` files
- Document your tools with good help text
- Consider adding examples in your tool's help output

## Troubleshooting

### Common Issues

### "Tool name already exists"

- Choose a different name or remove the existing tool directory

### "Build fails after generation"

- Run `go mod tidy` to ensure dependencies are properly managed
- Check that all imports are valid

### Getting Help

- Check the generated tool's `--help` output
- Review the template source in `internal/generator/templates.go`
- Run `make lint` to check for common issues
- Use VS Code's integrated debugging features
