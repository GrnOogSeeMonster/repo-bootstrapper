# CLI Interface Contract

## Command: new (Interactive Mode)

**Invocation**:
```bash
repo-bootstrapper new
```

**Input**: Interactive prompts (user responds to each)

**Output** (stdout):
```
Creating new repository...
✓ Workspace created at ~/Documents/my-project
✓ Git initialized
✓ Spec Kit initialized
✓ Generated docs/constitution.md
✓ Generated docs/spec.md
✓ Generated docs/plan.md
✓ Generated docs/tasks.md
✓ Generated prompts/primer.md
✓ Generated README.md
✓ Generated LICENSE
✓ Created GitHub repo: owner/my-project
✓ Pushed to origin/main

Summary: Created 8 files, initialized git, created GitHub repo, pushed to origin.

Next steps:
  1. cd ~/Documents/my-project
  2. Open in your coding agent
  3. Paste contents of prompts/primer.md
  4. Run /implement to build your app
```

**Exit Code**: 0 on success, non-zero on error

**Error Output** (stderr):
```
Error: GitHub CLI not authenticated
Run: gh auth login
```

---

## Command: new --from-file (Non-Interactive Mode)

**Invocation**:
```bash
repo-bootstrapper new --from-file config.json
```

**Input File** (config.json):
```json
{
  "repo_owner": "myorg",
  "repo_name": "my-project",
  "visibility": "public",
  "license": "MIT",
  "language": "Go",
  "framework": "",
  "target_os": ["Windows", "macOS", "Linux"],
  "agent_choice": "claude",
  "package_manager": "go mod",
  "test_framework": "testing",
  "app_kind": "CLI",
  "issue_labels": "standard",
  "workspace_path": "~/Documents/my-project"
}
```

**Output**: Same as interactive mode

**Validation Errors**:
```
Error: Invalid configuration in config.json:
  - repo_name: must match ^[a-zA-Z0-9._-]+$
  - license: must be one of: MIT, Apache-2.0, GPL-3.0, BSD-3-Clause, MPL-2.0, ISC, Unlicense, None
```

---

## Command: new --dry-run

**Invocation**:
```bash
repo-bootstrapper new --dry-run
# or
repo-bootstrapper new --from-file config.json --dry-run
```

**Output**:
```
[DRY RUN] Would create directory: ~/Documents/my-project
[DRY RUN] Would initialize git repository
[DRY RUN] Would run: specify init my-project --ai claude
[DRY RUN] Would create file: docs/constitution.md (2,451 bytes)
[DRY RUN] Would create file: docs/spec.md (1,823 bytes)
[DRY RUN] Would create file: docs/plan.md (3,102 bytes)
[DRY RUN] Would create file: docs/tasks.md (156 bytes)
[DRY RUN] Would create file: prompts/primer.md (1,034 bytes)
[DRY RUN] Would create file: README.md (892 bytes)
[DRY RUN] Would create file: LICENSE (1,063 bytes)
[DRY RUN] Would run: git add .
[DRY RUN] Would run: git commit -m "Initial commit from repo-bootstrapper"
[DRY RUN] Would run: gh repo create myorg/my-project --public --source . --remote origin --push
[DRY RUN] Would run: gh label create bug --color d73a4a --description "Something isn't working"
[DRY RUN] Would run: gh label create enhancement --color a2eeef --description "New feature or request"

Summary: Would create 8 files, initialize git, create GitHub repo, create 4 labels, push to origin.
```

**Exit Code**: 0 (dry-run never fails operations, only validates input)

---

## Command: doctor

**Invocation**:
```bash
repo-bootstrapper doctor
```

**Purpose**: Check prerequisites and environment

**Output**:
```
Checking prerequisites...

✓ Go 1.22.1 installed
✓ Git 2.42.0 installed
✓ GitHub CLI 2.40.1 installed
✓ GitHub CLI authenticated as user: myusername
✓ Spec Kit CLI 0.5.0 installed

All prerequisites satisfied.
```

**Failure Example**:
```
Checking prerequisites...

✓ Go 1.22.1 installed
✓ Git 2.42.0 installed
✗ GitHub CLI not found
  Install: https://cli.github.com/
✗ Spec Kit CLI not found
  Install: npm install -g @specify/cli

Missing 2 required tools.
```

**Exit Code**: 0 if all checks pass, 1 if any fail

---

## Command: version

**Invocation**:
```bash
repo-bootstrapper version
```

**Output**:
```
repo-bootstrapper v1.0.0
Built: 2025-10-01T14:32:00Z
Commit: abc123def456
Go: 1.22.1
Platform: darwin/arm64
```

**Exit Code**: 0

---

## Global Flags

**--debug**: Enable verbose structured logging (JSON to stderr)
```bash
repo-bootstrapper new --debug
```

Output (stderr):
```json
{"level":"info","time":"2025-10-01T14:32:00Z","msg":"Starting wizard"}
{"level":"debug","time":"2025-10-01T14:32:01Z","msg":"Validating repo name","input":"my-project","valid":true}
{"level":"info","time":"2025-10-01T14:32:05Z","msg":"Creating workspace","path":"/Users/me/Documents/my-project"}
```

**--help**: Show help for command
```bash
repo-bootstrapper new --help
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0    | Success |
| 1    | General error (invalid input, file I/O failure) |
| 2    | Prerequisite missing (gh CLI, git, specify not found) |
| 3    | GitHub authentication required |
| 4    | User cancelled operation |

---

## Environment Variables

**GITHUB_TOKEN** or **GH_TOKEN**: GitHub personal access token
- Used by `gh` CLI for authentication
- Never read or stored by repo-bootstrapper
- Only checked indirectly via `gh auth status`

**DEBUG**: Set to "true" or "1" to enable debug logging (alternative to --debug flag)

---

## Input Validation Behavior

**Interactive Mode**:
- Invalid input: Re-prompt with error message
- Example: "Invalid repo name. Use only letters, numbers, dots, underscores, and hyphens."

**Non-Interactive Mode**:
- Invalid JSON: Abort with parse error and line number
- Failed validation: Abort with all validation errors listed
- No partial execution - validate first, then execute

**Dry-Run Mode**:
- All validation still performed
- Operations logged but not executed
- Always exits 0 unless validation fails
