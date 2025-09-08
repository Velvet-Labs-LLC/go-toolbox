# Go Toolbox Dev Container

This project is ready to use with [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers) or GitHub Codespaces for a modern, reproducible Go development environment.

## Features

- Go 1.23 (latest)
- Docker-in-Docker support
- Node.js (LTS) and npm
- GitHub CLI
- Preinstalled Go tools: golangci-lint v2, gosec, goimports
- Preinstalled npm tools: prettier, eslint
- VS Code extensions for Go, Docker, GitHub Actions, Prettier, ESLint
- Ports 8080 (HTTP) and 8443 (HTTPS) auto-forwarded
- Non-root `vscode` user for best practices

## Getting Started

1. Open this folder in VS Code and install the "Dev Containers" extension if prompted.
2. Reopen in Container when prompted (or use the Command Palette: "Dev Containers: Reopen in Container").
3. The container will build and run `make dev-setup` automatically.
4. Start developing! All tools and dependencies are ready.

## Customization

- Edit `.devcontainer/devcontainer.json` or `.devcontainer/Dockerfile` to add more tools or change settings.
- See [VS Code Dev Containers documentation](https://code.visualstudio.com/docs/devcontainers/containers) for more info.

---

This setup ensures every developer has a consistent, secure, and modern environment for Go, Node, and cloud-native development.
