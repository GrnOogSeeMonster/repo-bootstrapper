# Repo Bootstrapper

A Go CLI that takes a new project from nothing to a pushed GitHub repository: an
interactive wizard collects the project's shape, then it creates the workspace, renders
a constitution/spec/plan/tasks document set, initialises git, wires up Spec Kit for your
AI agent of choice, and creates and pushes the remote.

---

## Development Status

**Complete end to end.** `repo-bootstrapper new` runs the whole pipeline — wizard →
workspace → templates → git init → Spec Kit → GitHub create and push — in both
interactive and `--from-file` modes, with `--dry-run` short-circuiting every side effect.
Last worked on 1 October 2025.

| | |
|---|---|
| Size | ~1,750 lines Go across 5 internal packages |
| Working | Interactive wizard (promptui), config-file mode, validation, workspace creation, six document templates, git init with conventional commits, Spec Kit integration, GitHub repo creation and push via `gh`, structured logging, dry-run |
| Incomplete | `doctor` only reports that Go is present — it prints `"Full prerequisite checking not yet implemented"` and checks nothing else |
| Tests | One unit test file covering the wizard (135 lines). Nothing else is covered |
| CI | GitHub Actions for lint (golangci-lint), cross-platform test matrix, and release |
| Verified | Builds to a single 20 MB static binary (`dist/repo-bootstrapper`, built 1 October 2025) |

### Package layout

| Package | Lines | |
|---|---|---|
| `internal/workspace` | 461 | Directory creation and rendering of the six templates below |
| `internal/git` | 343 | `init.go` — repo init, conventional first commit; `github.go` — `gh repo create` and push |
| `internal/wizard` | 372 | Prompts, answers, config, validation |
| `internal/logger` | 134 | Structured logging plus the operation log that drives dry-run output |
| `internal/speckit` | 78 | Shells out to `specify init` |
| `cmd/repo-bootstrapper` | 227 | Cobra CLI — `new`, `doctor`, `version` |

### Generated documents

`constitution.tmpl` · `spec.tmpl` · `plan.tmpl` · `tasks.tmpl` · `primer.tmpl` ·
`README.tmpl`

### Next

Finish `doctor` — it should verify `git`, `gh` (and its auth state) and `specify` before
`new` fails partway through, which is the only rough edge left in the user-facing flow.

---

## Quick Start (90 seconds)

### Installation

```bash
# From source
git clone https://github.com/<your-org>/repo-bootstrapper.git
cd repo-bootstrapper
make build
sudo cp dist/repo-bootstrapper /usr/local/bin/
```

### Usage

```bash
# Check prerequisites
repo-bootstrapper doctor

# Create a new repository (interactive wizard)
repo-bootstrapper new

# Preview operations without executing
repo-bootstrapper new --dry-run

# Non-interactive mode with config file
repo-bootstrapper new --from-file config.json
```

## Features

- **Interactive Q&A Wizard**: Collects project metadata through guided prompts
- **Git Integration**: Initializes repository with conventional commits
- **Spec Kit Setup**: Runs `specify init` for AI agent integration
- **Documentation Generation**: Creates constitution, spec, plan, and primer files
- **GitHub Integration**: Creates remote repository via `gh` CLI and pushes content
- **Cross-Platform**: Single static binary for Windows, macOS, and Linux
- **Dry-Run Mode**: Preview all operations before execution
- **Secure**: Never persists secrets, reads from environment only

## Prerequisites

- Go 1.22+ (for building from source)
- Git 2.x
- GitHub CLI (`gh`) - authenticated
- Spec Kit CLI (`npm install -g @specify/cli`)

## Commands

### `new`
Create a new bootstrapped repository

```bash
repo-bootstrapper new                    # Interactive wizard
repo-bootstrapper new --from-file config.json  # Non-interactive
repo-bootstrapper new --dry-run          # Preview only
```

### `doctor`
Check prerequisites and environment

```bash
repo-bootstrapper doctor
```

### `version`
Show version information

```bash
repo-bootstrapper version
```

## Configuration File Format

For non-interactive mode, create a JSON file:

```json
{
  "repo_owner": "myorg",
  "repo_name": "my-project",
  "visibility": "public",
  "license": "MIT",
  "language": "Go",
  "target_os": ["Windows", "macOS", "Linux"],
  "agent_choice": "claude",
  "package_manager": "go mod",
  "test_framework": "testing",
  "app_kind": "CLI",
  "issue_labels": "standard",
  "workspace_path": "~/Documents/my-project"
}
```

## Supported Licenses

- MIT
- Apache-2.0
- GPL-3.0
- BSD-3-Clause
- MPL-2.0
- ISC
- Unlicense
- None (skip license)

## Supported AI Agents

- Claude Code (`claude`)
- Cursor (`cursor`)
- GitHub Copilot (`copilot`)
- Gemini (`gemini`)

## Building

See [docs/building.md](docs/building.md) for detailed build instructions.

## Architecture

See [docs/architecture.md](docs/architecture.md) for system design.

## License

MIT - See LICENSE file

## Contributing

This tool follows the constitution defined in `.specify/memory/constitution.md`.
All contributions must adhere to the constitutional principles.

## Support

- Report issues: https://github.com/<your-org>/repo-bootstrapper/issues
- Documentation: [docs/](docs/)
- Spec Kit: https://specify.dev
