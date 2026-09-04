# Tasks: Repo Bootstrapper CLI

**Input**: Design documents from `/specs/001-build-a-local/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `cmd/`, `internal/`, `tests/` at repository root
- Paths shown below assume single project structure per plan.md

## Phase 3.1: Setup
- [x] T001 Initialize Go module with `go mod init github.com/<owner>/repo-bootstrapper` in repository root
- [x] T002 [P] Create project directory structure: cmd/repo-bootstrapper/, internal/{wizard,workspace,git,speckit,logger}/, tests/{unit,integration}/, docs/templates/
- [x] T003 [P] Create .gitignore with Go-specific patterns (vendor/, dist/, *.exe, *.test, .DS_Store)
- [x] T004 [P] Install dependencies: `go get github.com/spf13/cobra github.com/manifoldco/promptui github.com/go-git/go-git/v5 github.com/stretchr/testify`
- [x] T005 [P] Configure golangci-lint with .golangci.yml (enable gofmt, govet, staticcheck, errcheck)
- [x] T006 [P] Create Makefile with targets: build, test, lint, clean, release

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Unit Tests (Parallel - Different Packages)
- [ ] T007 [P] Write test for WizardAnswers validation in tests/unit/wizard_test.go (test invalid repo names, missing required fields, invalid licenses)
- [ ] T008 [P] Write test for workspace directory creation in tests/unit/workspace_test.go (test existing dir prompt, permission errors, path sanitization)
- [ ] T009 [P] Write test for template rendering in tests/unit/workspace_test.go (test variable substitution, missing variables, invalid templates)
- [ ] T010 [P] Write test for git initialization in tests/unit/git_test.go (mock git commands, test init, add, commit flow)
- [ ] T011 [P] Write test for GitHub CLI integration in tests/unit/git_test.go (mock gh commands, test auth check, repo create, label create)
- [ ] T012 [P] Write test for Spec Kit integration in tests/unit/speckit_test.go (mock specify command, test init with different agents)

### CLI Contract Tests (Parallel - Different Commands)
- [ ] T013 [P] Write test for `new` command interactive mode in tests/unit/cli_new_test.go (test prompt flow, validation, error handling)
- [ ] T014 [P] Write test for `new --from-file` command in tests/unit/cli_new_test.go (test JSON parsing, validation errors, file not found)
- [ ] T015 [P] Write test for `new --dry-run` command in tests/unit/cli_new_test.go (verify no filesystem writes, operation logging)
- [ ] T016 [P] Write test for `doctor` command in tests/unit/cli_doctor_test.go (test prerequisite checks, missing tools, version detection)
- [ ] T017 [P] Write test for `version` command in tests/unit/cli_version_test.go (test version output format)

### Integration Test
- [ ] T018 Write end-to-end dry-run test in tests/integration/dryrun_test.go (full wizard flow, verify all operations logged, no side effects)

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Data Models (Parallel - Different Files)
- [ ] T019 [P] Implement WizardAnswers struct with validation tags in internal/wizard/answers.go
- [ ] T020 [P] Implement WorkspaceConfig struct in internal/workspace/config.go
- [ ] T021 [P] Implement DocumentationArtifact struct in internal/workspace/artifact.go
- [ ] T022 [P] Implement OperationLog and Operation types in internal/logger/operations.go
- [ ] T023 [P] Implement Label struct with predefined label sets in internal/git/labels.go

### Validation & Input (Parallel - Different Packages)
- [ ] T024 [P] Implement input validation functions in internal/wizard/validation.go (repo name, owner, path, license allowlist)
- [ ] T025 [P] Implement interactive prompts using promptui in internal/wizard/prompts.go (all Q&A questions with validation)
- [ ] T026 [P] Implement JSON config file loading and validation in internal/wizard/config.go

### Template System (Sequential - Same Package)
- [ ] T027 Create template files in docs/templates/: constitution.tmpl, spec.tmpl, plan.tmpl, primer.tmpl
- [ ] T028 Implement template engine using text/template in internal/workspace/templates.go (compile templates, render with variables)
- [ ] T029 Embed license templates using go:embed in internal/workspace/licenses.go

### Workspace Operations (Sequential - Shared State)
- [ ] T030 Implement workspace directory creation with collision detection in internal/workspace/creator.go
- [ ] T031 Implement document generation from templates in internal/workspace/generator.go
- [ ] T032 Implement file writing with operation logging in internal/workspace/writer.go

### Git Integration (Sequential - Depends on git library)
- [ ] T033 Implement git initialization using go-git in internal/git/init.go (init repo, set default branch to main)
- [ ] T034 Implement git add and commit operations in internal/git/commit.go
- [ ] T035 Implement GitHub CLI wrapper in internal/git/github.go (auth check, repo create, label create, push)

### Spec Kit Integration
- [ ] T036 Implement Spec Kit CLI wrapper in internal/speckit/integration.go (check if installed, run specify init command)

### Logging System
- [ ] T037 Implement structured logger wrapping log/slog in internal/logger/logger.go (JSON mode for --debug, human mode for default)
- [ ] T038 Implement operation summary generation in internal/logger/summary.go (count files created, git operations, etc.)

### CLI Commands (Sequential - Wire Everything Together)
- [ ] T039 Implement root command with Cobra in cmd/repo-bootstrapper/main.go (global --debug and --help flags)
- [ ] T040 Implement `new` command with interactive mode in cmd/repo-bootstrapper/cmd/new.go (orchestrate wizard → workspace → git → github flow)
- [ ] T041 Implement `new --from-file` flag handling in cmd/repo-bootstrapper/cmd/new.go
- [ ] T042 Implement `new --dry-run` flag handling in cmd/repo-bootstrapper/cmd/new.go (log operations without execution)
- [ ] T043 Implement `doctor` command in cmd/repo-bootstrapper/cmd/doctor.go (check go, git, gh, specify versions)
- [ ] T044 Implement `version` command in cmd/repo-bootstrapper/cmd/version.go (display version, build date, commit hash)

## Phase 3.4: Integration & Error Handling
- [ ] T045 Wire all components together in main.go (dependency injection, error propagation)
- [ ] T046 Implement graceful error handling and rollback on failure (delete partially created workspace)
- [ ] T047 Implement exit code handling (0 success, 1 error, 2 missing prereq, 3 auth required, 4 cancelled)
- [ ] T048 Add context and cancellation support for user interrupts (Ctrl+C handling)

## Phase 3.5: Polish

### Documentation (Parallel - Different Files)
- [ ] T049 [P] Write README.md with 90-second quickstart, installation instructions, and examples
- [ ] T050 [P] Write docs/architecture.md with ASCII component diagram, data flow, and decision log
- [ ] T051 [P] Write docs/building.md with build prerequisites, commands, cross-compilation, and release process

### CI/CD (Parallel - Different Workflows)
- [ ] T052 [P] Create .github/workflows/ci.yml (golangci-lint, go test, go build on PR)
- [ ] T053 [P] Create .github/workflows/release.yml (cross-compile, create release, upload binaries on tag)

### Makefile Targets
- [ ] T054 Implement `make build` target (build for current OS, output to dist/)
- [ ] T055 Implement `make test` target (run all tests with coverage)
- [ ] T056 Implement `make lint` target (run golangci-lint)
- [ ] T057 Implement `make release` target (cross-compile for Windows/macOS/Linux, generate SHA256 checksums)
- [ ] T058 Implement `make clean` target (remove dist/, coverage files)

### Final Verification
- [ ] T059 Run `make lint` and fix all linting errors
- [ ] T060 Run `make test` and verify all tests pass with >80% coverage
- [ ] T061 Run integration test with real filesystem (not dry-run)
- [ ] T062 Test manual scenario: Interactive wizard → create workspace → verify GitHub repo → verify generated files
- [ ] T063 Test manual scenario: Non-interactive mode with config.json → verify same result as interactive
- [ ] T064 Test manual scenario: Dry-run mode → verify no filesystem side effects
- [ ] T065 Build release binaries for all platforms and verify they run

## Dependencies

**Critical Path** (must be sequential):
1. Setup (T001-T006) before everything
2. Tests (T007-T018) before implementation
3. Data models (T019-T023) before business logic
4. Core logic (T024-T038) before CLI commands
5. CLI commands (T039-T044) before integration
6. Integration (T045-T048) before polish
7. Polish (T049-T065) last

**Blocking Dependencies**:
- T027 (templates) blocks T028 (template engine)
- T028 blocks T031 (document generation)
- T033 (git init) blocks T034 (git commit)
- T034 blocks T035 (github operations)
- T039 (root command) blocks T040-T044 (subcommands)
- T045 (wiring) blocks T046-T048 (error handling)

## Parallel Execution Examples

### Phase 1: Setup Tasks
```bash
# Launch T002-T006 together (different files, no dependencies):
Task: "Create project directory structure: cmd/repo-bootstrapper/, internal/{wizard,workspace,git,speckit,logger}/, tests/{unit,integration}/, docs/templates/"
Task: "Create .gitignore with Go-specific patterns (vendor/, dist/, *.exe, *.test, .DS_Store)"
Task: "Install dependencies: go get github.com/spf13/cobra github.com/manifoldco/promptui github.com/go-git/go-git/v5 github.com/stretchr/testify"
Task: "Configure golangci-lint with .golangci.yml (enable gofmt, govet, staticcheck, errcheck)"
Task: "Create Makefile with targets: build, test, lint, clean, release"
```

### Phase 2: Unit Tests (Different Packages)
```bash
# Launch T007-T012 together:
Task: "Write test for WizardAnswers validation in tests/unit/wizard_test.go"
Task: "Write test for workspace directory creation in tests/unit/workspace_test.go"
Task: "Write test for template rendering in tests/unit/workspace_test.go"
Task: "Write test for git initialization in tests/unit/git_test.go"
Task: "Write test for GitHub CLI integration in tests/unit/git_test.go"
Task: "Write test for Spec Kit integration in tests/unit/speckit_test.go"
```

### Phase 2: CLI Contract Tests
```bash
# Launch T013-T017 together:
Task: "Write test for new command interactive mode in tests/unit/cli_new_test.go"
Task: "Write test for new --from-file command in tests/unit/cli_new_test.go"
Task: "Write test for new --dry-run command in tests/unit/cli_new_test.go"
Task: "Write test for doctor command in tests/unit/cli_doctor_test.go"
Task: "Write test for version command in tests/unit/cli_version_test.go"
```

### Phase 3: Data Models
```bash
# Launch T019-T023 together (different files):
Task: "Implement WizardAnswers struct with validation tags in internal/wizard/answers.go"
Task: "Implement WorkspaceConfig struct in internal/workspace/config.go"
Task: "Implement DocumentationArtifact struct in internal/workspace/artifact.go"
Task: "Implement OperationLog and Operation types in internal/logger/operations.go"
Task: "Implement Label struct with predefined label sets in internal/git/labels.go"
```

### Phase 3: Validation & Input
```bash
# Launch T024-T026 together (different packages):
Task: "Implement input validation functions in internal/wizard/validation.go"
Task: "Implement interactive prompts using promptui in internal/wizard/prompts.go"
Task: "Implement JSON config file loading and validation in internal/wizard/config.go"
```

### Phase 5: Documentation
```bash
# Launch T049-T051 together:
Task: "Write README.md with 90-second quickstart"
Task: "Write docs/architecture.md with ASCII component diagram"
Task: "Write docs/building.md with build prerequisites and commands"
```

### Phase 5: CI/CD
```bash
# Launch T052-T053 together:
Task: "Create .github/workflows/ci.yml"
Task: "Create .github/workflows/release.yml"
```

## Notes
- [P] tasks = different files, no dependencies → safe to parallelize
- Verify tests fail before implementing (TDD enforcement)
- Commit after each logical group of tasks
- Run `make lint && make test` frequently
- Avoid: vague tasks, same file conflicts, implementation before tests

## Task Generation Rules
*Applied during main() execution*

1. **From CLI Contracts**:
   - Each command → contract test task [P]
   - Each command → implementation task

2. **From Data Model**:
   - Each entity → struct definition task [P]
   - Each entity → validation test [P]

3. **From Quickstart Scenarios**:
   - Each scenario → integration test
   - Each prerequisite check → doctor command check

4. **Ordering**:
   - Setup → Tests → Models → Logic → CLI → Integration → Polish
   - Tests always before corresponding implementation
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [x] All CLI commands have corresponding tests (new, doctor, version)
- [x] All entities have struct definitions (WizardAnswers, WorkspaceConfig, DocumentationArtifact, OperationLog, Label)
- [x] All tests come before implementation (T007-T018 before T019+)
- [x] Parallel tasks truly independent (different files/packages)
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Template files creation included (T027)
- [x] Makefile targets fully specified (T054-T058)
- [x] Manual verification scenarios included (T062-T064)
