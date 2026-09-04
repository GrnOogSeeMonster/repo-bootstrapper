# Quickstart: Repo Bootstrapper CLI

## Prerequisites
- Go 1.22+ installed
- Git installed and configured
- GitHub CLI (`gh`) installed and authenticated
- Spec Kit CLI installed (`npm install -g @specify/cli`)

## Installation

```bash
# From source
git clone https://github.com/<owner>/repo-bootstrapper.git
cd repo-bootstrapper
make build
sudo cp dist/repo-bootstrapper /usr/local/bin/

# Or download binary from releases
# Linux/macOS
curl -LO https://github.com/<owner>/repo-bootstrapper/releases/latest/download/repo-bootstrapper-linux-amd64
chmod +x repo-bootstrapper-linux-amd64
sudo mv repo-bootstrapper-linux-amd64 /usr/local/bin/repo-bootstrapper

# Windows
# Download repo-bootstrapper.exe from releases and add to PATH
```

## Usage Scenario 1: Interactive Wizard

```bash
# Check prerequisites
repo-bootstrapper doctor

# Start interactive wizard
repo-bootstrapper new
```

**Wizard Flow**:
1. Enter repository owner/organization: `myorg`
2. Enter repository name: `my-awesome-app`
3. Select visibility: `public`
4. Select license: `MIT`
5. Enter primary language: `Go`
6. Enter framework (optional): `(press Enter to skip)`
7. Select target OS (space to select, enter to confirm): `Windows`, `macOS`, `Linux`
8. Select AI agent: `claude`
9. Enter package manager: `go mod`
10. Enter test framework: `testing`
11. Select application kind: `CLI`
12. Select GitHub issue labels: `standard`
13. Enter workspace path (default ~/Documents/my-awesome-app): `(press Enter for default)`
14. Enter initial features (comma-separated, optional): `(press Enter to skip)`
15. Confirm and create? `y`

**Result**:
```
✓ Workspace created at ~/Documents/my-awesome-app
✓ Git initialized
✓ Spec Kit initialized
✓ Generated docs/constitution.md
✓ Generated docs/spec.md
✓ Generated docs/plan.md
✓ Generated docs/tasks.md
✓ Generated prompts/primer.md
✓ Generated README.md
✓ Generated LICENSE
✓ Created GitHub repo: myorg/my-awesome-app
✓ Pushed to origin/main

Next steps:
  1. cd ~/Documents/my-awesome-app
  2. Open in Claude Code
  3. Paste contents of prompts/primer.md
  4. Run /implement to build your app
```

## Usage Scenario 2: Non-Interactive with Config File

Create `config.json`:
```json
{
  "repo_owner": "myorg",
  "repo_name": "my-awesome-app",
  "visibility": "public",
  "license": "MIT",
  "language": "Go",
  "target_os": ["Windows", "macOS", "Linux"],
  "agent_choice": "claude",
  "package_manager": "go mod",
  "test_framework": "testing",
  "app_kind": "CLI",
  "issue_labels": "standard",
  "workspace_path": "~/Documents/my-awesome-app"
}
```

Run:
```bash
repo-bootstrapper new --from-file config.json
```

## Usage Scenario 3: Dry Run (Preview Operations)

```bash
# Interactive dry-run
repo-bootstrapper new --dry-run

# Non-interactive dry-run
repo-bootstrapper new --from-file config.json --dry-run
```

Output shows all operations that would be performed without executing them.

## Usage Scenario 4: Debugging Issues

Enable debug logging:
```bash
repo-bootstrapper new --debug 2> debug.log
```

Review structured JSON logs in `debug.log`.

## Testing the Generated Workspace

```bash
cd ~/Documents/my-awesome-app

# Review generated files
ls -la
cat prompts/primer.md

# Open in Claude Code
code .

# Follow primer instructions to implement the app
```

## Validation Tests

Test input validation:
```bash
# Invalid repo name
echo '{"repo_name":"invalid name!"}' | repo-bootstrapper new --from-file /dev/stdin
# Expected: Error with validation message

# Missing required field
echo '{"repo_owner":"me"}' | repo-bootstrapper new --from-file /dev/stdin
# Expected: Error listing missing fields
```

## Expected Generated Structure

After running `repo-bootstrapper new`, the workspace contains:

```
my-awesome-app/
├── .git/                     # Git repository
├── .specify/                 # Spec Kit configuration
│   ├── memory/
│   │   └── constitution.md
│   └── templates/
├── docs/
│   ├── constitution.md       # Target app principles
│   ├── spec.md               # Target app specification
│   ├── plan.md               # Target app implementation plan
│   └── tasks.md              # Placeholder for agent
├── prompts/
│   └── primer.md             # Agent workflow guide
├── README.md                 # Project introduction
├── LICENSE                   # License file (if selected)
└── .gitignore                # Standard ignores
```

## Success Criteria

1. **Workspace created**: Directory exists at specified path
2. **Git initialized**: `.git` directory present, initial commit made
3. **Spec Kit ready**: `.specify` directory present with templates
4. **Docs generated**: All 5 documentation files exist and are non-empty
5. **GitHub repo created**: Remote repository exists and is accessible
6. **Pushed to origin**: Local commits appear in remote repository
7. **Labels created**: Issue labels present in GitHub repo (if configured)

## Validation Checklist

- [ ] Run `repo-bootstrapper doctor` - all checks pass
- [ ] Run interactive wizard - completes without errors
- [ ] Generated workspace contains all expected files
- [ ] `git log` shows initial commit
- [ ] `git remote -v` shows correct GitHub URL
- [ ] GitHub repo exists and contains pushed content
- [ ] `docs/constitution.md` contains project-specific principles
- [ ] `prompts/primer.md` contains valid slash-command sequence
- [ ] LICENSE file matches selected license type
- [ ] README contains project name and structure

## Troubleshooting

**Issue**: "GitHub CLI not authenticated"
```bash
gh auth login
# Follow prompts to authenticate
```

**Issue**: "Spec Kit CLI not found"
```bash
npm install -g @specify/cli
```

**Issue**: "Directory already exists"
- Choose different workspace path, or
- Delete existing directory, or
- Respond 'yes' when prompted to overwrite/merge

**Issue**: "Invalid repository name"
- Use only letters, numbers, dots, underscores, hyphens
- No spaces or special characters

**Issue**: Network timeout during GitHub operations
- Check internet connection
- Try again (gh CLI will retry automatically)
- Use `--dry-run` to test locally first

## Performance Expectations

- Interactive wizard: <100ms per prompt response
- Full bootstrap (local only): <5 seconds
- Full bootstrap (with GitHub): <15 seconds (network dependent)
- Dry-run mode: <2 seconds

## Next Steps After Bootstrap

1. Navigate to workspace: `cd ~/Documents/my-awesome-app`
2. Open in your chosen AI coding agent
3. Read `prompts/primer.md`
4. Run slash-commands in order:
   - `/constitution` (paste docs/constitution.md)
   - `/specify` (paste docs/spec.md)
   - `/clarify`
   - `/plan` (paste docs/plan.md)
   - `/tasks`
   - `/implement`
5. Your app is built!
