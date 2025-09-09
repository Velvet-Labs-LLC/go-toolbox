# Single Binary Packaging Strategies

## Overview

This document outlines different approaches to package all CLI and TUI tools into a single binary, providing users with multiple deployment options. We've implemented multiple strategies to give maximum flexibility.

## 🎯 Approach 1: Embedded Binary (Recommended) ✨

**Location**: `cmd/embedded/main.go`

### **How It Works**
- Single binary that maximally reuses existing code from CLI and TUI components
- Detects execution mode based on binary name or first argument
- Supports subcommands and symlink detection
- All functionality embedded in one optimized executable

### **Usage Examples**
```bash
# Default CLI mode with all subcommands
./embedded --help
./embedded file hash README.md
./embedded network ping google.com
./embedded utils string reverse "hello world"

# TUI mode - interactive terminal interface
./embedded tui

# Server mode - HTTP file server
./embedded serve ./docs --port 8080 --tls

# Symlink support for convenience
ln -s embedded toolbox-tui
./toolbox-tui  # Automatically runs in TUI mode

ln -s embedded toolbox-serve  
./toolbox-serve ./docs --port 8080  # Automatically runs server
```

### **Architecture**
```go
main()
├── detectMode() // Binary name or first argument detection
├── runCLIMode() // Full CLI with reused command implementations
├── runTUIMode() // TUI with reused existing TUI models  
└── runServerMode() // HTTP server with reused server logic
```

### **Code Reuse Benefits**
- ✅ **Maximum Reuse**: All CLI command functions directly reused from `cmd/cli/main/main.go`
- ✅ **TUI Integration**: Reuses existing TUI models and generator components
- ✅ **Shared Logic**: Configuration, logging, and utilities shared across modes
- ✅ **Maintainability**: Updates to base components automatically benefit embedded version

## 🔄 Approach 2: Unified Binary (Alternative)

**Location**: `cmd/unified/main.go`

### **How It Works**
- Single binary that detects execution mode
- Supports subcommands and symlink detection  
- Alternative implementation approach

### **Usage Examples**
```bash
# Default CLI mode
./cmd-unified --help

# Explicit mode selection
./cmd-unified tui
./cmd-unified serve --addr :9090
./cmd-unified generate template myapp

# Symlink support
ln -s cmd-unified toolbox-tui
./toolbox-tui  # Automatically runs in TUI mode
```

### **Benefits**
- ✅ **Single download**: Users get everything in one file
- ✅ **Symlink support**: Create tool-specific shortcuts  
- ✅ **Backward compatible**: Same CLI interface
- ✅ **Alternative approach**: Different implementation strategy

## 📊 Comparison: Embedded vs Unified vs Individual

| Feature | Embedded | Unified | Individual |
|---------|----------|---------|------------|
| **Binary Count** | 1 | 1 | 4+ |
| **Code Reuse** | Maximum | Moderate | Minimal |
| **Maintenance** | Easiest | Moderate | Most complex |
| **File Size** | Smallest | Medium | Largest total |
| **Deployment** | Simplest | Simple | Complex |
| **Mode Switching** | Fast | Fast | Process spawn |
| **Memory Usage** | Most efficient | Efficient | Higher overhead |

## 🎯 Recommendation

**Use the Embedded Binary (`./bin/embedded`)** for:
- ✅ **Production deployments**
- ✅ **Container images** 
- ✅ **End-user distribution**
- ✅ **CI/CD pipelines**
- ✅ **Maximum efficiency**

**Use Individual Binaries** for:
- 🔧 **Development and testing**
- 🔧 **Debugging specific components**
- 🔧 **Legacy compatibility**
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
