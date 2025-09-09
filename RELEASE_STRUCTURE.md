# Release Asset Organization

## Overview

The CI workflow now organizes build artifacts into a structured directory layout for better user experience and clearer separation between CLI and TUI tools.

## New Release Structure

### **Directory Layout**
```
release/
├── cli/                              # Command-line tools
│   ├── main-linux-amd64             # CLI main application (Linux)
│   ├── main-linux-arm64             # CLI main application (Linux ARM)
│   ├── main-windows-amd64.exe       # CLI main application (Windows)
│   ├── main-darwin-amd64            # CLI main application (macOS Intel)
│   ├── main-darwin-arm64            # CLI main application (macOS Apple Silicon)
│   ├── serve-linux-amd64            # CLI serve application (Linux)
│   ├── serve-linux-arm64            # CLI serve application (Linux ARM)
│   ├── serve-windows-amd64.exe      # CLI serve application (Windows)
│   ├── serve-darwin-amd64           # CLI serve application (macOS Intel)
│   └── serve-darwin-arm64           # CLI serve application (macOS Apple Silicon)
├── tui/                              # Terminal UI tools
│   ├── main-linux-amd64             # TUI main application (Linux)
│   ├── main-linux-arm64             # TUI main application (Linux ARM)
│   ├── main-windows-amd64.exe       # TUI main application (Windows)
│   ├── main-darwin-amd64            # TUI main application (macOS Intel)
│   └── main-darwin-arm64            # TUI main application (macOS Apple Silicon)
├── go-toolbox-cli-v1.0.0.zip        # CLI tools archive (ZIP)
├── go-toolbox-cli-v1.0.0.tar.gz     # CLI tools archive (TAR.GZ)
├── go-toolbox-tui-v1.0.0.zip        # TUI tools archive (ZIP)
├── go-toolbox-tui-v1.0.0.tar.gz     # TUI tools archive (TAR.GZ)
├── go-toolbox-all-v1.0.0.zip        # Combined archive (ZIP)
└── go-toolbox-all-v1.0.0.tar.gz     # Combined archive (TAR.GZ)
```

### **Archive Contents**

#### **CLI Archives** (`go-toolbox-cli-*.{zip,tar.gz}`)
Contains all command-line interface tools:
- `cli/main-*` - Main CLI application for all platforms
- `cli/serve-*` - Server CLI application for all platforms

#### **TUI Archives** (`go-toolbox-tui-*.{zip,tar.gz}`)
Contains all terminal user interface tools:
- `tui/main-*` - Main TUI application for all platforms

#### **Combined Archives** (`go-toolbox-all-*.{zip,tar.gz}`)
Contains both CLI and TUI tools in organized directories.

## Benefits

### **For Users**
- ✅ **Clear separation**: Easy to identify CLI vs TUI tools
- ✅ **Selective downloads**: Can download only CLI or TUI tools if needed
- ✅ **Organized structure**: Tools are logically grouped
- ✅ **Multiple formats**: Both ZIP and TAR.GZ archives available

### **For Developers**
- ✅ **Maintainable**: Clear build organization
- ✅ **Extensible**: Easy to add new tool categories
- ✅ **Debuggable**: Easy to identify which tools built successfully
- ✅ **Professional**: Enterprise-ready release structure

## Build Process Changes

### **Before (Flat Structure)**
```bash
bin/
├── cli-main-linux-amd64
├── cli-serve-linux-amd64  
├── tui-main-linux-amd64
├── cli-main-windows-amd64.exe
├── cli-serve-windows-amd64.exe
└── tui-main-windows-amd64.exe
```

### **After (Organized Structure)**
```bash
bin/
├── cli/
│   ├── main-linux-amd64
│   ├── main-windows-amd64.exe
│   ├── serve-linux-amd64
│   └── serve-windows-amd64.exe
└── tui/
    ├── main-linux-amd64
    └── main-windows-amd64.exe
```

## Usage Examples

### **Download Specific Tool Type**
```bash
# Download only CLI tools
wget https://github.com/Velvet-Labs-LLC/go-toolbox/releases/download/v1.0.0/go-toolbox-cli-v1.0.0.tar.gz

# Download only TUI tools  
wget https://github.com/Velvet-Labs-LLC/go-toolbox/releases/download/v1.0.0/go-toolbox-tui-v1.0.0.tar.gz

# Download everything
wget https://github.com/Velvet-Labs-LLC/go-toolbox/releases/download/v1.0.0/go-toolbox-all-v1.0.0.tar.gz
```

### **Installation Scripts**
```bash
# Install CLI tools only
tar -xzf go-toolbox-cli-v1.0.0.tar.gz
sudo cp cli/* /usr/local/bin/
chmod +x /usr/local/bin/main /usr/local/bin/serve

# Install TUI tools only
tar -xzf go-toolbox-tui-v1.0.0.tar.gz  
sudo cp tui/* /usr/local/bin/
chmod +x /usr/local/bin/main

# Install all tools
tar -xzf go-toolbox-all-v1.0.0.tar.gz
sudo cp cli/* tui/* /usr/local/bin/
chmod +x /usr/local/bin/*
```

### **Docker Usage**
```dockerfile
# Multi-stage build using organized binaries
FROM alpine:latest as cli-tools
WORKDIR /tools
ADD go-toolbox-cli-v1.0.0.tar.gz .

FROM alpine:latest as tui-tools  
WORKDIR /tools
ADD go-toolbox-tui-v1.0.0.tar.gz .

# Final image with specific tools
FROM alpine:latest
COPY --from=cli-tools /tools/cli/main /usr/local/bin/go-toolbox
COPY --from=cli-tools /tools/cli/serve /usr/local/bin/go-toolbox-serve
RUN chmod +x /usr/local/bin/*
```

## GitHub Release Page

### **Release Assets Display**
The GitHub release page will show:

```
📦 Assets (8)

🗂️ Archive Downloads:
- go-toolbox-all-v1.0.0.tar.gz          (Combined - All tools)
- go-toolbox-all-v1.0.0.zip             (Combined - All tools)  
- go-toolbox-cli-v1.0.0.tar.gz          (CLI tools only)
- go-toolbox-cli-v1.0.0.zip             (CLI tools only)
- go-toolbox-tui-v1.0.0.tar.gz          (TUI tools only)
- go-toolbox-tui-v1.0.0.zip             (TUI tools only)

📁 Individual Binaries:
- [Individual platform binaries as separate uploads]
```

## Backward Compatibility

### **Migration Path**
- ✅ **Old scripts continue working**: Archives contain the same binaries
- ✅ **New organization available**: Users can adopt organized structure gradually
- ✅ **Multiple download options**: Both organized and flat structures available

### **Breaking Changes**
- ❌ **Binary names changed**: `cli-main` → `main` (in cli/ directory)
- ❌ **Path structure changed**: Binaries now in subdirectories
- ✅ **Mitigation**: Provide both individual files and archives

## Monitoring and Validation

### **Build Verification**
The CI workflow now includes verification steps:
```bash
# Verify CLI binaries exist and work
for binary in bin/cli/*linux-amd64*; do
  $binary --help >/dev/null 2>&1 && echo "✅ $binary" || echo "❌ $binary"
done

# Verify TUI binaries exist and work  
for binary in bin/tui/*linux-amd64*; do
  $binary --help >/dev/null 2>&1 && echo "✅ $binary" || echo "❌ $binary"
done
```

### **Release Validation**
```bash
# Verify archive contents
tar -tzf go-toolbox-cli-v1.0.0.tar.gz | head -10
tar -tzf go-toolbox-tui-v1.0.0.tar.gz | head -10
tar -tzf go-toolbox-all-v1.0.0.tar.gz | head -10

# Verify file permissions
tar -tzvf go-toolbox-all-v1.0.0.tar.gz | grep -E "(cli|tui)/"
```

## Future Enhancements

### **Potential Additions**
- 📱 **Mobile builds**: Add `mobile/` directory for mobile-specific tools
- 🐳 **Container images**: Add Docker images to releases
- 📚 **Documentation**: Include man pages or help files in archives
- 🔑 **Checksums**: Add SHA256 checksums for all archives
- 📋 **Package metadata**: Add package.json or similar metadata files

### **Tool Categories**
```bash
release/
├── cli/          # Command-line interface tools
├── tui/          # Terminal user interface tools  
├── web/          # Web-based tools (future)
├── mobile/       # Mobile-specific tools (future)
├── plugins/      # Plugin system tools (future)
└── docs/         # Documentation and man pages (future)
```

This organized structure provides a professional, user-friendly release experience while maintaining flexibility for future expansion.
