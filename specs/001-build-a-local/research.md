# Research: Repo Bootstrapper CLI

## Technology Decisions

### 1. CLI Framework Choice

**Decision**: Use Cobra

**Rationale**:
- Industry standard for Go CLI tools (kubectl, gh, hugo all use Cobra)
- Built-in support for subcommands, flags, help generation
- Well-documented with active maintenance
- Integrates well with viper for configuration (if needed later)

**Alternatives Considered**:
- **urfave/cli**: Simpler API, less opinionated, but less feature-rich
- **flag (stdlib)**: Too low-level, would require significant boilerplate
- Cobra chosen for best balance of power and community support

### 2. Interactive Prompt Library

**Decision**: Use promptui

**Rationale**:
- Clean, simple API for building interactive prompts
- Supports selection menus, text input, validation
- Cross-platform terminal handling
- Lightweight with minimal dependencies

**Alternatives Considered**:
- **survey**: More feature-rich but heavier dependency footprint
- **bubbletea**: Modern TUI framework but overkill for simple wizard
- promptui chosen for simplicity and constitutional compliance

### 3. Git Integration Approach

**Decision**: Use go-git library + shell out to `gh` CLI

**Rationale**:
- go-git: Pure Go implementation for basic git operations (init, add, commit)
- `gh` CLI: Official GitHub tool for repo creation, handles auth automatically
- Hybrid approach: programmatic git ops locally, delegate GitHub API to gh
- No need to manage GitHub tokens directly - gh handles auth

**Alternatives Considered**:
- **git2go (libgit2 bindings)**: C dependency breaks "single static binary" goal
- **google/go-github**: Would require manual token management, increases complexity
- Hybrid approach chosen for security and simplicity

### 4. Template Engine

**Decision**: Use text/template (Go stdlib)

**Rationale**:
- Built into Go standard library - no external dependency
- Sufficient power for document generation (conditionals, loops, variables)
- Well-tested and documented
- Aligns with "simple, battle-tested" constitutional principle

**Alternatives Considered**:
- **html/template**: Unnecessary escaping for non-HTML content
- **Handlebars-go**: External dependency not justified for basic templating
- stdlib text/template chosen for zero dependencies

### 5. Logging Strategy

**Decision**: Custom logger wrapping log/slog (Go 1.21+)

**Rationale**:
- slog: Structured logging in Go stdlib since 1.21
- JSON output for --debug, human-readable for default
- Constitutional requirement for structured logging met without dependencies
- Simple facade can wrap slog for custom formatting

**Alternatives Considered**:
- **logrus**: Popular but adds dependency
- **zap**: High performance but overkill for CLI tool
- **zerolog**: Great but stdlib solution preferred per constitution
- slog chosen for zero dependencies and modern design

### 6. Configuration File Format

**Decision**: JSON for non-interactive mode input

**Rationale**:
- encoding/json in stdlib
- Widely supported and easy to generate programmatically
- Schema validation straightforward
- Human-readable enough for small config files

**Alternatives Considered**:
- **YAML**: Requires external dependency, spec quirks
- **TOML**: Requires external dependency, less common
- JSON chosen for stdlib support and ubiquity

### 7. Build System

**Decision**: Makefile + optional GoReleaser

**Rationale**:
- Makefile: Universal, simple, transparent build commands
- GoReleaser: Automates cross-compilation and artifact generation
- Make required, GoReleaser optional (documented for maintainers)
- Aligns with "reproducible builds" principle

**Build Targets**:
- `make build`: Build for current OS
- `make test`: Run all tests
- `make lint`: Run golangci-lint
- `make release`: Cross-compile for Windows/macOS/Linux (via GoReleaser)
- `make clean`: Remove build artifacts

### 8. Testing Strategy

**Decision**: Go testing package + testify/assert + afero (virtual filesystem)

**Rationale**:
- testing (stdlib): Industry standard
- testify/assert: Minimal assertion library for readability
- afero: Virtual filesystem for testing file operations without side-effects
- Mocking gh CLI calls via interface abstraction

**Test Coverage Targets**:
- Unit tests: wizard validation, template rendering, git operations
- Integration test: Full dry-run workflow (no remote side-effects)
- No arbitrary coverage threshold (per constitution)

### 9. License Templates Source

**Decision**: Embed license text files in binary using Go embed

**Rationale**:
- Go 1.16+ supports //go:embed directive
- Licenses are static text - embed at compile time
- No external API calls or file dependencies at runtime
- Supports all clarified licenses: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, MPL-2.0, ISC, Unlicense

**Implementation**:
```go
//go:embed templates/licenses/*.txt
var licenseFS embed.FS
```

### 10. Default Workspace Location

**Decision**: Use os.UserHomeDir() + "/Documents"

**Rationale**:
- Cross-platform: Works on Windows, macOS, Linux
- Per clarification: User Documents folder is default
- Fallback to current directory if Documents doesn't exist
- User can override via wizard prompt

**Platform Paths**:
- Windows: C:\Users\<username>\Documents
- macOS: /Users/<username>/Documents
- Linux: /home/<username>/Documents

## Integration Patterns

### Spec Kit Integration

**Pattern**: Shell out to `specify` CLI

**Implementation**:
```go
cmd := exec.Command("specify", "init", projectName, "--ai", agentChoice)
cmd.Dir = workspaceDir
output, err := cmd.CombinedOutput()
```

**Error Handling**:
- Check if `specify` exists via `exec.LookPath("specify")`
- If not found: provide clear error with installation instructions
- If fails: capture output and show to user

### GitHub CLI Integration

**Pattern**: Shell out to `gh` CLI for repo creation

**Implementation**:
```go
// Check auth status first
authCmd := exec.Command("gh", "auth", "status")
if err := authCmd.Run(); err != nil {
    return ErrGitHubNotAuthenticated
}

// Create repo
visibility := "--public" // or "--private"
createCmd := exec.Command("gh", "repo", "create",
    fmt.Sprintf("%s/%s", owner, repoName),
    visibility, "--source", ".", "--remote", "origin", "--push")
```

**Fallback**:
- If gh not installed or not authenticated: print manual instructions
- Never fail the entire operation - local workspace still valid

## Security Considerations

### Input Validation Rules

1. **Repo Name**:
   - Regex: `^[a-zA-Z0-9._-]+$`
   - Max length: 100 characters
   - No leading/trailing hyphens

2. **Repo Owner/Org**:
   - Regex: `^[a-zA-Z0-9-]+$`
   - Max length: 39 characters (GitHub limit)

3. **Directory Paths**:
   - Use filepath.Clean to prevent directory traversal
   - Validate absolute paths
   - Check for invalid characters per OS

4. **License Selection**:
   - Allowlist: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, MPL-2.0, ISC, Unlicense, None
   - Case-insensitive comparison
   - Reject any value not in allowlist

### Token Handling

- Read from environment: `GITHUB_TOKEN` or `GH_TOKEN`
- Never store in files or logs
- Never pass as CLI argument
- Delegate to `gh` CLI which manages keyring

## Performance Optimizations

- Templates compiled once at startup
- File operations buffered
- Parallel git operations not needed (simple init + add + commit)
- Wizard responses collected before any file I/O (fail fast on validation)

## Error Recovery

- Dry-run mode: No actual operations, just print intent
- Operation log: Track each file created for potential rollback
- Atomic operations where possible (create temp, then rename)
- Clear error messages with remediation steps

## Documentation Requirements

**docs/architecture.md**: Component diagram (ASCII), data flow, decision log
**docs/building.md**: Prerequisites, build commands, release process, cross-compilation
**README.md**: Install, quickstart (90 seconds), example usage, troubleshooting

---

**Research Complete**: All NEEDS CLARIFICATION resolved. Ready for Phase 1 design.
