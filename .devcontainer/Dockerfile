# syntax=docker/dockerfile:1
FROM mcr.microsoft.com/devcontainers/go:1.23-bullseye

# Install additional tools and dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        curl \
        build-essential \
        ca-certificates \
        openssh-client \
        unzip \
        jq \
        bash-completion \
        sudo \
        python3-pip \
        nodejs \
        npm \
    && rm -rf /var/lib/apt/lists/*

# Install golangci-lint, gosec, goimports
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0 \
    && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest \
    && go install golang.org/x/tools/cmd/goimports@latest

# Install latest npm global tools for JS/TS dev
RUN npm install -g npm@latest prettier eslint

# Set up a non-root user (vscode)
ARG USERNAME=vscode
RUN usermod -aG sudo $USERNAME

USER $USERNAME
WORKDIR /workspace
