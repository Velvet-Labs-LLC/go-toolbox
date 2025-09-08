#!/bin/bash
# Build script for the toolbox project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Toolbox Project...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go first.${NC}"
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | cut -d' ' -f3)
echo -e "${YELLOW}Using Go version: $GO_VERSION${NC}"

# Create bin directory if it doesn't exist
mkdir -p bin

# Build applications
echo -e "${GREEN}Building CLI application...${NC}"
go build -ldflags="-s -w" -o bin/toolbox ./cmd/cli/main

echo -e "${GREEN}Building TUI application...${NC}"
go build -ldflags="-s -w" -o bin/toolbox-tui ./cmd/tui/main

# Check if builds were successful
if [ -f "bin/toolbox" ] && [ -f "bin/toolbox-tui" ]; then
    echo -e "${GREEN}Build completed successfully!${NC}"
    echo -e "${YELLOW}Binaries created:${NC}"
    echo "  - bin/toolbox (CLI)"
    echo "  - bin/toolbox-tui (TUI)"
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi

# Show binary sizes
echo -e "${YELLOW}Binary sizes:${NC}"
ls -lh bin/

echo -e "${GREEN}Build script completed!${NC}"
