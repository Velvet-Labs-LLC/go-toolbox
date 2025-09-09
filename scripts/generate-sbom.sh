#!/bin/bash

# Script to generate Software Bill of Materials (SBOM) for the Go project
# This script generates SBOM files in multiple formats

set -e

echo "🔍 Generating SBOM for go-toolbox..."

# Create output directory
mkdir -p ./sbom-output

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Syft if not available
if ! command_exists syft; then
    echo "📦 Installing Syft..."
    curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin
fi

# Generate SBOM in different formats
echo "📋 Generating SPDX JSON format..."
syft . -o spdx-json=./sbom-output/go-toolbox-sbom.spdx.json

echo "📋 Generating CycloneDX JSON format..."
syft . -o cyclonedx-json=./sbom-output/go-toolbox-sbom.cyclonedx.json

echo "📋 Generating Syft JSON format..."
syft . -o syft-json=./sbom-output/go-toolbox-sbom.syft.json

echo "📋 Generating SPDX Tag-Value format..."
syft . -o spdx-tag-value=./sbom-output/go-toolbox-sbom.spdx

echo "📋 Generating CycloneDX XML format..."
syft . -o cyclonedx-xml=./sbom-output/go-toolbox-sbom.cyclonedx.xml

# Generate a human-readable table format
echo "📋 Generating human-readable table format..."
syft . -o table=./sbom-output/go-toolbox-sbom.txt

echo "✅ SBOM generation complete!"
echo "📁 Files generated in ./sbom-output/:"
ls -la ./sbom-output/

echo ""
echo "📄 SBOM Files Generated:"
echo "  • SPDX JSON: ./sbom-output/go-toolbox-sbom.spdx.json"
echo "  • CycloneDX JSON: ./sbom-output/go-toolbox-sbom.cyclonedx.json"
echo "  • Syft JSON: ./sbom-output/go-toolbox-sbom.syft.json"
echo "  • SPDX Tag-Value: ./sbom-output/go-toolbox-sbom.spdx"
echo "  • CycloneDX XML: ./sbom-output/go-toolbox-sbom.cyclonedx.xml"
echo "  • Human-readable: ./sbom-output/go-toolbox-sbom.txt"

echo ""
echo "🔒 To scan for vulnerabilities, run:"
echo "  grype ./sbom-output/go-toolbox-sbom.syft.json"
