# Building Repo Bootstrapper

## Prerequisites

### Required Tools

1. **Go 1.22+**
   ```bash
   # Check version
   go version

   # Install on Linux/WSL
   wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin

   # Install on macOS
   brew install go@1.22

   # Install on Windows
   # Download from: https://go.dev/dl/go1.22.0.windows-amd64.msi
   ```

2. **Git 2.x**
   ```bash
   git --version
   ```

3. **Make** (optional, for convenience)
   ```bash
   make --version
   ```

### Optional Tools

- **golangci-lint**: For linting (make lint)
- **GitHub CLI (gh)**: For testing GitHub integration
- **Spec Kit CLI**: For testing Spec Kit integration

## Build Commands

### Quick Build (Current Platform)

```bash
# Using Make
make build

# Using Go directly
go build -o dist/repo-bootstrapper ./cmd/repo-bootstrapper

# The binary will be in dist/repo-bootstrapper
```

### Build with Version Info

```bash
VERSION=v1.0.0
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT=$(git rev-parse --short HEAD)

go build \
  -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.Commit=$COMMIT" \
  -o dist/repo-bootstrapper \
  ./cmd/repo-bootstrapper
```

### Cross-Compilation (All Platforms)

```bash
# Using Make
make release

# This creates:
# - dist/repo-bootstrapper-windows-amd64.exe
# - dist/repo-bootstrapper-darwin-amd64
# - dist/repo-bootstrapper-darwin-arm64
# - dist/repo-bootstrapper-linux-amd64
# - dist/checksums.txt (SHA256 sums)
```

### Manual Cross-Compilation

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o dist/repo-bootstrapper.exe ./cmd/repo-bootstrapper

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o dist/repo-bootstrapper-mac ./cmd/repo-bootstrapper

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o dist/repo-bootstrapper-mac-arm ./cmd/repo-bootstrapper

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/repo-bootstrapper-linux ./cmd/repo-bootstrapper
```

## Testing

### Run All Tests

```bash
# Using Make
make test

# Using Go directly
go test -v -race -coverprofile=coverage.txt ./...

# View coverage report
go tool cover -html=coverage.txt
```

### Run Specific Tests

```bash
# Unit tests only
go test ./tests/unit/...

# Integration tests only
go test ./tests/integration/...

# Specific package
go test ./internal/wizard/...

# Specific test
go test -run TestWizardAnswers_Validation ./tests/unit/...
```

## Linting

### Run Linter

```bash
# Using Make
make lint

# Using golangci-lint directly
golangci-lint run ./...
```

### Install golangci-lint

```bash
# Binary installation
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.0

# Or via Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Development Workflow

### 1. Make Changes

Edit source files in `cmd/` or `internal/`

### 2. Run Tests

```bash
make test
```

### 3. Build

```bash
make build
```

### 4. Test Locally

```bash
./dist/repo-bootstrapper doctor
./dist/repo-bootstrapper new --dry-run
```

### 5. Lint

```bash
make lint
```

### 6. Commit

```bash
git add .
git commit -m "feat: add new feature"
```

## Clean Build

```bash
# Using Make
make clean

# Manual
rm -rf dist/ coverage.txt
```

## Installation

### Install Locally

```bash
# Using Make
make install

# This copies the binary to $GOPATH/bin/repo-bootstrapper
```

### Install Globally (Linux/macOS)

```bash
sudo cp dist/repo-bootstrapper /usr/local/bin/
```

### Install Globally (Windows)

1. Copy `dist/repo-bootstrapper.exe` to a directory
2. Add that directory to your PATH

## Release Process

### 1. Tag Release

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 2. Build Release Binaries

```bash
make release
```

### 3. Generate Checksums

```bash
cd dist
sha256sum * > checksums.txt
```

### 4. Create GitHub Release

```bash
# Using GitHub CLI
gh release create v1.0.0 \
  dist/repo-bootstrapper-windows-amd64.exe \
  dist/repo-bootstrapper-darwin-amd64 \
  dist/repo-bootstrapper-darwin-arm64 \
  dist/repo-bootstrapper-linux-amd64 \
  dist/checksums.txt \
  --title "v1.0.0" \
  --notes "Release notes here"
```

### 5. Verify Release

```bash
# Download and test each binary
wget https://github.com/<org>/repo-bootstrapper/releases/download/v1.0.0/repo-bootstrapper-linux-amd64
chmod +x repo-bootstrapper-linux-amd64
./repo-bootstrapper-linux-amd64 version
```

## Troubleshooting

### Build Fails: "go: command not found"

**Solution**: Install Go 1.22+ and add to PATH

```bash
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Build Fails: "cannot find package"

**Solution**: Download dependencies

```bash
go mod download
go mod tidy
```

### Tests Fail: "race detector not supported"

**Solution**: Remove `-race` flag for your platform

```bash
go test -v ./...  # Without race detector
```

### Large Binary Size

**Solution**: Strip debug symbols

```bash
go build -ldflags="-s -w" -o dist/repo-bootstrapper ./cmd/repo-bootstrapper
```

This reduces binary size by ~30%.

### Cross-Compilation Issues

**Problem**: CGo dependencies don't cross-compile easily

**Solution**: This project uses pure Go dependencies only (no CGo). If you add CGo dependencies, you'll need platform-specific toolchains.

## CI/CD Integration

### GitHub Actions

See `.github/workflows/ci.yml` for:
- Automated testing on PRs
- Linting checks
- Build verification

See `.github/workflows/release.yml` for:
- Automated release builds
- Asset uploads
- Checksum generation

### Local CI Simulation

```bash
# Simulate CI checks locally
make lint && make test && make build
```

## Performance

- **Build Time**: ~30 seconds (clean build)
- **Incremental Build**: ~5 seconds
- **Test Time**: ~10 seconds (full suite)
- **Binary Size**: ~20MB (uncompressed)

## Dependencies Management

### Update Dependencies

```bash
# Update all to latest minor/patch
go get -u ./...
go mod tidy

# Update specific dependency
go get github.com/spf13/cobra@latest

# View dependency tree
go mod graph
```

### Vendor Dependencies (Optional)

```bash
go mod vendor

# Build using vendor
go build -mod=vendor -o dist/repo-bootstrapper ./cmd/repo-bootstrapper
```

## Static Analysis

### Run All Static Checks

```bash
go vet ./...
staticcheck ./...
golangci-lint run ./...
```

### Security Scanning

```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Scan for security issues
gosec ./...
```

## Additional Resources

- Go Build Documentation: https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
- Cross Compilation: https://golang.org/doc/install/source#environment
- golangci-lint Docs: https://golangci-lint.run/
