# Architecture: Repo Bootstrapper

## Overview

Repo Bootstrapper is a Go CLI tool that automates repository creation following best practices and constitutional principles. It uses a modular architecture organized by domain concerns.

## System Design

```
┌─────────────┐
│   User      │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────┐
│         CLI Layer (Cobra)                    │
│  ┌─────┐  ┌────────┐  ┌────────┐           │
│  │ new │  │ doctor │  │version │           │
│  └──┬──┘  └───┬────┘  └───┬────┘           │
│     │         │            │                 │
└─────┼─────────┼────────────┼─────────────────┘
      │         │            │
      ▼         ▼            ▼
┌─────────────────────────────────────────────┐
│      Internal Packages                       │
│                                              │
│  ┌─────────┐  ┌───────────┐  ┌──────────┐  │
│  │ wizard  │  │ workspace │  │   git    │  │
│  │         │  │           │  │          │  │
│  │ prompts │  │ templates │  │  github  │  │
│  │ answers │  │  creator  │  │   init   │  │
│  │validate │  │           │  │          │  │
│  └────┬────┘  └─────┬─────┘  └────┬─────┘  │
│       │             │              │        │
│       │    ┌────────┴──────┐      │        │
│       │    │   speckit     │      │        │
│       │    │ integration   │      │        │
│       │    └───────────────┘      │        │
│       │                            │        │
│       └────────┬───────────────────┘        │
│                │                            │
│         ┌──────┴──────┐                     │
│         │   logger    │                     │
│         │  operations │                     │
│         │   summary   │                     │
│         └─────────────┘                     │
└─────────────────────────────────────────────┘
       │                    │
       ▼                    ▼
┌─────────────┐      ┌────────────┐
│ Filesystem  │      │  External  │
│             │      │   Tools    │
│ - Workspace │      │ - gh CLI   │
│ - Templates │      │ - specify  │
│ - Git repo  │      │ - git      │
└─────────────┘      └────────────┘
```

## Component Responsibilities

### CLI Layer (`cmd/repo-bootstrapper`)
- **main.go**: Entry point, Cobra root command setup
- **Responsibilities**:
  - Parse flags and arguments
  - Route to appropriate subcommand
  - Handle global `--debug` flag
  - Exit code management

### Wizard Package (`internal/wizard`)
- **answers.go**: WizardAnswers struct with validation
- **prompts.go**: Interactive Q&A using promptui (TODO)
- **validation.go**: Input validation rules (TODO)
- **config.go**: JSON config file loading (TODO)
- **Responsibilities**:
  - Collect user input interactively or from file
  - Validate all inputs against rules
  - Return structured WizardAnswers

### Workspace Package (`internal/workspace`)
- **creator.go**: Directory and file creation (TODO)
- **templates.go**: Template rendering (TODO)
- **config.go**: WorkspaceConfig tracking (TODO)
- **Responsibilities**:
  - Create workspace directory (handle collisions)
  - Render documentation from templates
  - Write LICENSE, README, .gitignore
  - Track all created files

### Git Package (`internal/git`)
- **init.go**: Git initialization (TODO)
- **github.go**: GitHub CLI integration (TODO)
- **labels.go**: Issue label management (TODO)
- **Responsibilities**:
  - Initialize git repository
  - Add and commit files
  - Check gh CLI authentication
  - Create remote repository
  - Push to origin
  - Create issue labels

### Spec Kit Package (`internal/speckit`)
- **integration.go**: Spec Kit CLI wrapper (TODO)
- **Responsibilities**:
  - Check if `specify` CLI is installed
  - Execute `specify init` with correct arguments
  - Capture and report output

### Logger Package (`internal/logger`)
- **logger.go**: Structured logging wrapper
- **operations.go**: Operation tracking
- **Responsibilities**:
  - JSON logs when `--debug` enabled
  - Human-readable logs by default
  - Track all operations for dry-run and summary
  - Generate operation summaries

## Data Flow

### Interactive Mode

```
User runs: repo-bootstrapper new

1. Wizard collects inputs via promptui
   → RepoOwner, RepoName, License, etc.
   → Validates each input immediately
   → Returns WizardAnswers

2. Workspace creator initializes
   → Creates ~/Documents/<repo-name>/
   → Checks for existing directory
   → Prompts user if collision

3. Template engine renders docs
   → Uses WizardAnswers as variables
   → Generates constitution.md, spec.md, etc.
   → Writes to workspace/docs/

4. Git operations execute
   → git init
   → git add .
   → git commit -m "Initial commit"

5. GitHub integration (if authenticated)
   → gh repo create owner/name
   → git remote add origin
   → git push -u origin main
   → gh label create (if configured)

6. Spec Kit integration
   → specify init <name> --ai <agent>

7. Summary printed
   → "Created X files, pushed to Y"
   → Next steps displayed
```

### Non-Interactive Mode

```
User runs: repo-bootstrapper new --from-file config.json

1. Load and parse JSON config
   → Unmarshal into WizardAnswers
   → Validate entire struct

2. [Same as steps 2-7 above]
```

### Dry-Run Mode

```
User runs: repo-bootstrapper new --dry-run

1. [All validation and planning occurs]

2. Operations logged but NOT executed
   → "[DRY RUN] Would create directory: ~/Documents/myrepo"
   → "[DRY RUN] Would run: git init"
   → "[DRY RUN] Would run: gh repo create..."

3. Summary shows what WOULD happen
   → "Would create 8 files, initialize git, create GitHub repo"
```

## Error Handling Strategy

- **Fail Fast**: Validate all inputs before any filesystem operations
- **Atomic Operations**: Use temp files, then rename for critical writes
- **Clear Messages**: User-facing errors are actionable (not stack traces)
- **Rollback**: On critical failures, attempt to clean up partial state
- **Exit Codes**:
  - 0: Success
  - 1: General error (invalid input, I/O failure)
  - 2: Missing prerequisite (gh, git, specify)
  - 3: Authentication required (gh not authenticated)
  - 4: User cancelled

## Security Model

- **No Secret Persistence**: Tokens never written to disk or logs
- **Environment Only**: GITHUB_TOKEN read from env at runtime
- **Delegation**: GitHub operations via `gh` CLI (manages keyring)
- **Input Sanitization**: All paths cleaned with filepath.Clean
- **Allowlists**: Licenses, agents, app kinds validated against allowlists

## Constitutional Compliance

This architecture implements all 10 constitutional principles:

1. **Simplicity**: Minimal abstractions, clear package boundaries
2. **Least Privilege**: Creates files in user-specified dirs only
3. **Reproducible Builds**: Go modules with lock files
4. **Clear Naming**: Self-documenting package and function names
5. **Cross-Platform**: Go's portability, static binary output
6. **Deterministic Toolchain**: Go 1.22+ required and documented
7. **Security**: Env-only secrets, input validation, safe errors
8. **Observability**: Structured logs, debug flag, operation summaries
9. **Documentation-First**: All docs in /docs, this file exists
10. **Test Discipline**: Tests before implementation (TDD)

## Extension Points

Future enhancements can be added without major refactoring:

- **New AI Agents**: Add to validation allowlist and agent-file generation
- **New Licenses**: Add template to embedded FS, update allowlist
- **New Integrations**: Create new package under `internal/`
- **Custom Templates**: Replace embedded templates with user-provided
- **Plugins**: Add plugin loader if complexity justified

## Performance Characteristics

- **Wizard**: <100ms per prompt response
- **Full Bootstrap**: <10 seconds (local operations only)
- **With GitHub**: +5-10 seconds (network dependent)
- **Memory**: <50MB typical usage
- **Binary Size**: ~20MB (Go + dependencies)

## Dependencies

- **github.com/spf13/cobra**: CLI framework
- **github.com/manifoldco/promptui**: Interactive prompts
- **github.com/go-git/go-git/v5**: Pure Go git implementation
- **github.com/stretchr/testify**: Test assertions
- **github.com/go-playground/validator/v10**: Struct validation
- **github.com/spf13/afero**: Virtual filesystem for testing

All dependencies are battle-tested, actively maintained, and align with constitutional principles.
