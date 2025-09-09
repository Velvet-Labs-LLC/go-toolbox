# Directory File Server

This directory contains the legacy serve command implementation. **This functionality is now integrated into the unified `go-toolbox` binary.**

## Current Usage (Embedded Binary)

Use the unified binary instead:

```bash
# Build the unified binary
make build-embedded

# Serve current directory
./bin/go-toolbox serve

# Serve specific directory with custom port
./bin/go-toolbox serve --dir /path/to/directory --port 8080

# Enable HTTPS with auto-generated certificates
./bin/go-toolbox serve --port 8443 --tls
```

## Available Options

- `--dir`: Directory to serve (default: current directory)
- `--port`: Port to listen on (default: 8080)
- `--tls`: Enable HTTPS with auto-generated certificates
- `--cert`: Path to TLS certificate file (requires --key)
- `--key`: Path to TLS private key file (requires --cert)

## Examples

```bash
# Basic HTTP server
./bin/go-toolbox serve --dir ./public --port 9000

# HTTPS server with auto-generated certificates
./bin/go-toolbox serve --dir ./docs --port 8443 --tls

# HTTPS server with existing certificates
./bin/go-toolbox serve --cert server.crt --key server.key --port 8443
```

## Features

- **HTTP/HTTPS support**: Both protocols with automatic certificate generation
- **Cross-platform**: Works on Linux, macOS, Windows, WSL2, containers
- **Network discovery**: Automatically detects and displays local IP addresses
- **File browsing**: Web interface for browsing and downloading files
- **Admin interface**: Management interface at `/admin` endpoint
- **Logging**: Structured logging with configurable levels
- **Configuration**: YAML-based configuration support

## Network Access

The server will display all available network interfaces when started:

```
Server started successfully!
Access the server at:
  Local:    http://localhost:8080
  Network:  http://192.168.1.100:8080
  Network:  http://10.0.0.50:8080
```

Any device on your local network can access the server using the network URLs.

## Legacy Note

This directory contains the original standalone serve implementation for reference. The code has been integrated into the unified binary with enhanced features and better code reuse.
