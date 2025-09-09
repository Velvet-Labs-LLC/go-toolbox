# Single Binary Packaging Strategies

## Overview

This document outlines different approaches to package all CLI and TUI tools into a single binary, providing users with multiple deployment options.

## 🎯 Approach 1: Unified Binary (Implemented)

**Location**: `cmd/unified/main.go`

### **How It Works**
- Single binary that detects execution mode
- Supports subcommands and symlink detection
- All functionality in one executable

### **Usage Examples**
```bash
# Default CLI mode
./go-toolbox --help

# Explicit mode selection
./go-toolbox tui
./go-toolbox serve --addr :9090
./go-toolbox generate template myapp

# Symlink support
ln -s go-toolbox toolbox-tui
./toolbox-tui  # Automatically runs in TUI mode

ln -s go-toolbox toolbox-serve  
./toolbox-serve --addr :8080  # Automatically runs server
```

### **Architecture**
```go
main()
├── detectMode() // Binary name or first argument
├── runCLIMode() // Default Cobra CLI with subcommands
├── runTUIMode() // TUI application
└── runServerMode() // HTTP server
```

### **Benefits**
- ✅ **Single download**: Users get everything in one file
- ✅ **Symlink support**: Create tool-specific shortcuts  
- ✅ **Backward compatible**: Same CLI interface
- ✅ **Smaller total size**: Shared dependencies
- ✅ **Easy deployment**: Just one binary to manage

### **Current Release Structure**
```
release/
├── unified/                           # Single all-in-one binary
│   ├── go-toolbox-linux-amd64        # 🎯 RECOMMENDED
│   ├── go-toolbox-windows-amd64.exe
│   └── go-toolbox-darwin-amd64
├── cli/                               # Separate CLI binaries
│   ├── main-linux-amd64
│   └── serve-linux-amd64
├── tui/                               # Separate TUI binaries
│   └── main-linux-amd64
├── go-toolbox-unified-v1.0.0.tar.gz  # 🎯 RECOMMENDED
├── go-toolbox-separate-v1.0.0.tar.gz # Individual binaries
└── go-toolbox-complete-v1.0.0.tar.gz # Everything
```

## 🔧 Approach 2: Plugin Architecture (Future)

For ultimate flexibility, consider a plugin-based approach:

```go
// Main binary that loads plugins dynamically
type Plugin interface {
    Name() string
    Description() string
    Execute(args []string) error
}

// Plugins are loaded from:
// - Built-in (compiled in)
// - External files (.so, .dll)
// - WebAssembly modules
```

## 🚀 Approach 3: Container-Based

Package everything in a container image:

```dockerfile
FROM alpine:latest
COPY go-toolbox /usr/local/bin/
RUN ln -s go-toolbox /usr/local/bin/toolbox-cli && \
    ln -s go-toolbox /usr/local/bin/toolbox-tui && \
    ln -s go-toolbox /usr/local/bin/toolbox-serve

# Usage: docker run go-toolbox tui
# Usage: docker run go-toolbox serve --addr :8080
```

## 📦 Implementation Details

### **Code Organization**
```
cmd/
├── unified/           # 🎯 Single binary (IMPLEMENTED)
│   └── main.go       # Mode detection and routing
├── cli/              # Separate CLI binaries
│   ├── main/
│   └── serve/
└── tui/              # Separate TUI binaries
    └── main/
```

### **To Complete the Implementation**

1. **Move Business Logic to Packages**:
```bash
# Move logic from cmd/*/main.go to internal packages
internal/
├── cli/
│   ├── main.go     # CLI application logic
│   └── serve.go    # Server application logic
├── tui/
│   └── main.go     # TUI application logic
└── unified/
    └── router.go   # Mode detection and routing
```

2. **Update cmd/unified/main.go**:
```go
// Import and call the actual implementations
func runTUIMode(args []string) {
    tui.Main(args) // Call from internal/tui
}

func runServerMode(args []string) {
    server.Main(args) // Call from internal/cli
}
```

3. **Add Build Tags** (Optional):
```go
//go:build !minimal
// +build !minimal

// Include all features in full build

//go:build minimal
// +build minimal

// Include only CLI in minimal build
```

## 🎯 Recommended User Journey

### **For Most Users** (Recommended)
```bash
# Download unified binary
wget https://github.com/user/go-toolbox/releases/download/v1.0.0/go-toolbox-unified-v1.0.0.tar.gz
tar -xzf go-toolbox-unified-v1.0.0.tar.gz
chmod +x unified/go-toolbox-linux-amd64

# Create convenient symlinks
ln -s go-toolbox-linux-amd64 go-toolbox
ln -s go-toolbox toolbox-cli
ln -s go-toolbox toolbox-tui  
ln -s go-toolbox toolbox-serve

# Use in any mode
./go-toolbox --help           # CLI mode
./toolbox-tui                # TUI mode  
./toolbox-serve --addr :8080 # Server mode
```

### **For Advanced Users**
```bash
# Download separate binaries if needed
wget https://github.com/user/go-toolbox/releases/download/v1.0.0/go-toolbox-separate-v1.0.0.tar.gz

# Or download everything
wget https://github.com/user/go-toolbox/releases/download/v1.0.0/go-toolbox-complete-v1.0.0.tar.gz
```

## 📊 Size Comparison

| Approach | Binary Count | Total Size | Benefits |
|----------|-------------|------------|----------|
| Separate | 4 binaries | ~40MB | Modular, specific |
| Unified | 1 binary | ~12MB | Convenient, shared deps |
| Container | 1 image | ~15MB | Portable, isolated |

## 🔄 Migration Strategy

### **Phase 1: Dual Release** (Current)
- Provide both unified and separate binaries
- Default recommendation: unified binary
- Keep separate binaries for compatibility

### **Phase 2: Unified Focus**
- Promote unified binary as primary
- Keep separate binaries for special cases
- Add container images

### **Phase 3: Plugin Architecture** (Future)
- Modular plugin system
- Dynamic loading capabilities
- WebAssembly support

## 🧪 Testing Strategy

### **Unified Binary Testing**
```bash
# Test all modes work
./go-toolbox --version          # CLI mode
./go-toolbox tui --help        # TUI mode
./go-toolbox serve --help      # Server mode
./go-toolbox generate --help   # CLI subcommand

# Test symlink detection
ln -s go-toolbox toolbox-tui
./toolbox-tui --help           # Should auto-detect TUI mode

# Test argument passing
./go-toolbox serve --addr :9090 --tls
```

### **CI/CD Integration**
The GitHub Actions workflow now:
1. ✅ Builds unified binary
2. ✅ Tests all modes
3. ✅ Creates organized releases
4. ✅ Provides multiple download options

## 💡 Advanced Features

### **Configuration Management**
```bash
# Unified config for all modes
go-toolbox config set --global theme dark
go-toolbox config set --cli default-output json
go-toolbox config set --tui animations true
go-toolbox config set --server default-port 8080
```

### **Plugin Discovery**
```bash
# List available functionality
go-toolbox list                # All available commands/modes
go-toolbox list --plugins      # Available plugins
go-toolbox list --modes        # Available modes (cli, tui, serve)
```

### **Shell Integration**
```bash
# Bash completion for all modes
go-toolbox completion bash > /etc/bash_completion.d/go-toolbox

# Works with all modes
go-toolbox <TAB>               # cli, tui, serve, generate, config
go-toolbox tui <TAB>           # TUI-specific options
go-toolbox serve <TAB>         # Server-specific options
```

## 🎉 Summary

The unified binary approach provides:

1. **🎯 Single Download**: One file contains everything
2. **🔄 Multiple Interfaces**: CLI, TUI, and Server modes
3. **🔗 Symlink Support**: Create tool-specific shortcuts
4. **📦 Organized Releases**: Multiple packaging options
5. **🚀 Easy Deployment**: Just one binary to manage
6. **⚡ Better Performance**: Shared dependencies reduce size
7. **🛠 Backward Compatible**: Same CLI interface

**Recommendation**: Use the unified binary (`go-toolbox-unified-*.tar.gz`) for most deployments, with separate binaries available for special cases.
