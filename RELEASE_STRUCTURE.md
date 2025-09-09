# Release Asset Organization

## Overview

The project has evolved to use an **embedded binary strategy** that combines all CLI, TUI, and server functionality into a single unified binary. This provides the best user experience with maximum simplicity and minimal distribution complexity.

## Current Release Structure

### **Primary Binary Strategy**
```
release/
├── go-toolbox-linux-amd64            # Unified binary (Linux x64)
├── go-toolbox-linux-arm64            # Unified binary (Linux ARM64)
├── go-toolbox-windows-amd64.exe      # Unified binary (Windows x64)
├── go-toolbox-darwin-amd64           # Unified binary (macOS Intel)
├── go-toolbox-darwin-arm64           # Unified binary (macOS Apple Silicon)
├── go-toolbox-v1.0.0.zip             # All binaries archive (ZIP)
└── go-toolbox-v1.0.0.tar.gz          # All binaries archive (TAR.GZ)
```

### **Legacy Binaries (Optional)**
For users who prefer separate binaries, the traditional structure is also available:
```
release/
├── legacy/
│   ├── cli/                          # Individual CLI binaries
│   │   ├── main-linux-amd64         
│   │   ├── serve-linux-amd64        
│   │   └── ...
│   └── tui/                          # Individual TUI binaries
│       ├── main-linux-amd64         
│       └── ...
```

### **Unified Binary Features**

The embedded binary (`go-toolbox`) includes all functionality in one executable:

#### **Automatic Mode Detection**
- **CLI Mode**: When run with command-line arguments
- **TUI Mode**: When run without arguments (interactive)
- **Server Mode**: When run with `serve` or `server` command

#### **Full Feature Set**
- ✅ **File Operations**: Copy, move, delete, organize files
- ✅ **Network Utilities**: Port scanning, connectivity tests
- ✅ **System Tools**: Process monitoring, disk usage
- ✅ **Code Generation**: Template-based tool generation
- ✅ **Configuration**: YAML-based configuration with validation
- ✅ **Logging**: Structured logging with multiple levels
- ✅ **Server**: HTTP/HTTPS server with admin interface

## Benefits

### **For Users**
- ✅ **Single binary**: Download and run one file
- ✅ **Consistent experience**: Same functionality across CLI/TUI modes
- ✅ **Small footprint**: Optimized binary size with shared code
- ✅ **Easy deployment**: Copy one file to any system
- ✅ **Auto-detection**: Switches modes based on usage context

### **For Developers**
- ✅ **Simplified builds**: Single compilation target
- ✅ **Maximum code reuse**: Shared logic across all modes
- ✅ **Unified testing**: Test all functionality in one binary
- ✅ **Easier maintenance**: Single codebase for all features

## Build Process Evolution

### **Legacy Approach (Multiple Binaries)**
```bash
# Old build process - multiple separate binaries
bin/
├── cli-main-linux-amd64
├── cli-serve-linux-amd64  
├── tui-main-linux-amd64
├── cli-main-windows-amd64.exe
├── cli-serve-windows-amd64.exe
└── tui-main-windows-amd64.exe
```

### **Current Approach (Embedded Binary)**
```bash
# New build process - single unified binary
bin/
├── go-toolbox-linux-amd64
├── go-toolbox-linux-arm64
├── go-toolbox-windows-amd64.exe
├── go-toolbox-darwin-amd64
└── go-toolbox-darwin-arm64
```

### **Build Commands**
```bash
# Build embedded binary for current platform
make build-embedded

# Build for all platforms
make build-all-embedded

# Build legacy binaries (if needed)
make build-legacy
```

## Usage Examples

### **Download and Installation**
```bash
# Download the binary for your platform
wget https://github.com/Velvet-Labs-LLC/go-toolbox/releases/download/v1.0.0/go-toolbox-linux-amd64

# Make it executable
chmod +x go-toolbox-linux-amd64

# Install to system PATH
sudo mv go-toolbox-linux-amd64 /usr/local/bin/go-toolbox
```

### **Usage Modes**

#### **CLI Mode (with arguments)**
```bash
# File operations
go-toolbox copy file1.txt file2.txt
go-toolbox move *.log /var/log/
go-toolbox organize /downloads --by-extension

# Network utilities
go-toolbox port-scan 192.168.1.1 1-1000
go-toolbox ping-test google.com

# Code generation
go-toolbox generate tool --name myapp --type cli
```

#### **TUI Mode (interactive)**
```bash
# Run without arguments for interactive TUI
go-toolbox

# Launches beautiful terminal interface with:
# - File browser and operations
# - Network tools with real-time output
# - System monitoring dashboards
# - Configuration management
```

#### **Server Mode**
```bash
# Start HTTP server
go-toolbox serve --port 8080

# Start HTTPS server with auto-generated certificates
go-toolbox serve --port 8443 --tls

# Server provides web interface for all tool functionality
```

### **Docker Usage**
```dockerfile
# Minimal Docker image with unified binary
FROM alpine:latest
WORKDIR /app

# Copy single binary
COPY go-toolbox-linux-amd64 /usr/local/bin/go-toolbox
RUN chmod +x /usr/local/bin/go-toolbox

# Default to server mode
EXPOSE 8080
CMD ["go-toolbox", "serve", "--port", "8080"]
```

### **CI/CD Integration**
```yaml
# GitHub Actions example
- name: Download go-toolbox
  run: |
    wget https://github.com/Velvet-Labs-LLC/go-toolbox/releases/latest/download/go-toolbox-linux-amd64
    chmod +x go-toolbox-linux-amd64
    ./go-toolbox-linux-amd64 organize ./build --by-extension
```

## GitHub Release Page

### **Release Assets Display**
The GitHub release page shows a clean, simple structure:

```
📦 Assets (7)

� Recommended Downloads:
- go-toolbox-linux-amd64                (Unified binary - Linux x64)
- go-toolbox-linux-arm64                (Unified binary - Linux ARM64)  
- go-toolbox-windows-amd64.exe          (Unified binary - Windows x64)
- go-toolbox-darwin-amd64                (Unified binary - macOS Intel)
- go-toolbox-darwin-arm64                (Unified binary - macOS Apple Silicon)

📦 Archive Downloads:
- go-toolbox-v1.0.0.tar.gz              (All binaries - TAR.GZ)
- go-toolbox-v1.0.0.zip                 (All binaries - ZIP)
```

## Migration Guide

### **From Legacy Binaries**
If you were using the old separate binaries, migration is straightforward:

```bash
# Old usage
./cli-main copy file1.txt file2.txt
./tui-main
./cli-serve --port 8080

# New usage (unified binary)
./go-toolbox copy file1.txt file2.txt
./go-toolbox  # TUI mode
./go-toolbox serve --port 8080
```

### **Configuration Migration**
- ✅ **Config files**: Existing `config.yaml` files work unchanged
- ✅ **Command syntax**: All CLI commands remain the same
- ✅ **Environment variables**: All existing environment variables supported
- ✅ **TUI interface**: Same keyboard shortcuts and navigation

### **Compatibility Notes**
- ❌ **Binary names changed**: `cli-main` → `go-toolbox`, `cli-serve` → `go-toolbox serve`
- ❌ **Separate binaries**: No longer needed - everything in one binary
- ✅ **Feature parity**: All functionality from separate binaries is included

## Monitoring and Validation

### **Build Verification**
The CI workflow includes verification for the embedded binary:
```bash
# Verify unified binary works in all modes
./go-toolbox --help                # CLI help
./go-toolbox version              # Version information
timeout 5 ./go-toolbox || true    # TUI mode (timeout after 5s)
./go-toolbox serve --help         # Server mode help

# Cross-platform build verification
for binary in bin/go-toolbox-*; do
  if [[ "$binary" == *"windows"* ]]; then
    echo "✅ $binary (Windows - skipping execution test)"
  else
    $binary --version >/dev/null 2>&1 && echo "✅ $binary" || echo "❌ $binary"
  fi
done
```

### **Release Validation**
```bash
# Verify archive contents
tar -tzf go-toolbox-v1.0.0.tar.gz
unzip -l go-toolbox-v1.0.0.zip

# Verify binary functionality
chmod +x go-toolbox-linux-amd64
./go-toolbox-linux-amd64 --version
./go-toolbox-linux-amd64 --help

# Verify file permissions and size
ls -la go-toolbox-* | grep -v "^d"
du -h go-toolbox-*
```

## Future Enhancements

### **Binary Optimizations**
- � **Size optimization**: Further reduce binary size with build flags
- ⚡ **Performance tuning**: Optimize startup time and memory usage
- � **Plugin system**: Dynamic loading of additional functionality
- � **Compression**: UPX compression for smaller distribution

### **Distribution Improvements**
- 🍺 **Package managers**: Homebrew, Chocolatey, APT/YUM repositories  
- 🐳 **Container images**: Official Docker images on Docker Hub
- 📱 **Mobile support**: Cross-compilation for mobile platforms
- 🔑 **Code signing**: Signed binaries for enhanced security

### **Feature Expansions**
- 🌐 **Web UI**: Enhanced web interface with real-time updates
- 📊 **Monitoring**: Built-in metrics and health checking
- 🔄 **Auto-updates**: Self-updating binary capability
- 🔌 **Extensions**: Third-party extension support

### **Development Tools**
- 🧪 **Testing**: Automated testing of all binary modes
- 📝 **Documentation**: Auto-generated documentation from code
- 🔍 **Debugging**: Built-in debugging and profiling tools
- 📈 **Analytics**: Usage analytics and performance monitoring

The embedded binary strategy provides a solid foundation for future enhancements while maintaining simplicity and ease of use.
