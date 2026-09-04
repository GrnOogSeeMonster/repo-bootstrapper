.PHONY: build test lint clean release

# Build variables
BINARY_NAME=repo-bootstrapper
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.Commit=$(COMMIT)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) for current platform..."
	@mkdir -p dist
	$(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME) ./cmd/repo-bootstrapper

# Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "Coverage report:"
	$(GOCMD) tool cover -func=coverage.txt

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/ coverage.txt coverage.html

# Cross-compile for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/repo-bootstrapper
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/repo-bootstrapper
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/repo-bootstrapper
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/repo-bootstrapper
	@echo "Generating checksums..."
	@cd dist && sha256sum * > checksums.txt
	@echo "Release binaries built successfully!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp dist/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
