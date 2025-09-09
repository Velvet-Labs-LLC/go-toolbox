#!/bin/bash

# Install Git hooks for the go-toolbox project
# This script sets up the Git hooks directory and installs the pre-commit hook

set -e

echo "🔧 Installing Git hooks for go-toolbox..."

# Check if we're in a Git repository
if [ ! -d ".git" ]; then
    echo "❌ Error: Not in a Git repository. Please run this from the project root."
    exit 1
fi

# Set the Git hooks directory to .githooks
echo "📁 Setting Git hooks directory to .githooks..."
git config core.hooksPath .githooks

# Verify the hooks directory is set correctly
hooks_path=$(git config core.hooksPath)
if [ "$hooks_path" = ".githooks" ]; then
    echo "✅ Git hooks directory set to: $hooks_path"
else
    echo "❌ Failed to set Git hooks directory"
    exit 1
fi

# Make sure all hook files are executable
echo "🔐 Making hook files executable..."
find .githooks -name "*" -type f -exec chmod +x {} \;

# List installed hooks
echo ""
echo "📋 Installed hooks:"
ls -la .githooks/

echo ""
echo "🎉 Git hooks installed successfully!"
echo ""
echo "The following hooks are now active:"
echo "  • pre-commit: Runs 'make fmt', 'make lint', and 'make test'"
echo ""
echo "To test the pre-commit hook:"
echo "  git commit -m \"test commit\""
echo ""
echo "To bypass hooks (not recommended):"
echo "  git commit --no-verify -m \"message\""
