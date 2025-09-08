# PowerShell build script for the toolbox project

param(
    [switch]$Clean,
    [switch]$Test,
    [switch]$Install
)

# Colors
$Green = "`e[32m"
$Yellow = "`e[33m"
$Red = "`e[31m"
$Reset = "`e[0m"

Write-Host "${Green}Building Toolbox Project...${Reset}"

# Check if Go is installed
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "${Red}Go is not installed. Please install Go first.${Reset}"
    exit 1
}

# Get Go version
$goVersion = go version
Write-Host "${Yellow}Using $goVersion${Reset}"

# Clean if requested
if ($Clean) {
    Write-Host "${Yellow}Cleaning build artifacts...${Reset}"
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
    }
    go clean
}

# Create bin directory if it doesn't exist
if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Run tests if requested
if ($Test) {
    Write-Host "${Yellow}Running tests...${Reset}"
    go test ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Host "${Red}Tests failed!${Reset}"
        exit 1
    }
}

# Build applications
Write-Host "${Green}Building CLI application...${Reset}"
go build -ldflags="-s -w" -o "bin/toolbox.exe" "./cmd/cli/main"

Write-Host "${Green}Building TUI application...${Reset}"
go build -ldflags="-s -w" -o "bin/toolbox-tui.exe" "./cmd/tui/main"

# Check if builds were successful
if ((Test-Path "bin/toolbox.exe") -and (Test-Path "bin/toolbox-tui.exe")) {
    Write-Host "${Green}Build completed successfully!${Reset}"
    Write-Host "${Yellow}Binaries created:${Reset}"
    Write-Host "  - bin/toolbox.exe (CLI)"
    Write-Host "  - bin/toolbox-tui.exe (TUI)"
} else {
    Write-Host "${Red}Build failed!${Reset}"
    exit 1
}

# Show binary sizes
Write-Host "${Yellow}Binary sizes:${Reset}"
Get-ChildItem "bin/*.exe" | Format-Table Name, @{Name="Size";Expression={"{0:N2} MB" -f ($_.Length / 1MB)}}

# Install if requested
if ($Install) {
    Write-Host "${Yellow}Installing applications...${Reset}"
    go install "./cmd/cli/main"
    go install "./cmd/tui/main"
    Write-Host "${Green}Applications installed to GOPATH/bin${Reset}"
}

Write-Host "${Green}Build script completed!${Reset}"
