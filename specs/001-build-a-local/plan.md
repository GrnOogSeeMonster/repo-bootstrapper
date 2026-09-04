
# Implementation Plan: Repo Bootstrapper CLI

**Branch**: `001-build-a-local` | **Date**: 2025-10-01 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-build-a-local/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Build a CLI tool that automates repository bootstrapping through an interactive Q&A wizard. The tool collects project metadata (repo name, owner, language, framework, agent choice, etc.), creates a local workspace with git initialization, runs Spec Kit setup, generates tailored documentation artifacts (constitution, spec, plan, tasks placeholder, primer), and publishes to GitHub via `gh` CLI. Supports dry-run mode, non-interactive JSON config mode, and cross-platform distribution as a single static binary.

## Technical Context
**Language/Version**: Go 1.22+
**Primary Dependencies**: Cobra or urfave/cli (CLI framework), promptui (interactive prompts), go-git (git operations)
**Storage**: Filesystem only - no database. Transient state in memory during execution.
**Testing**: Go testing package with testify for assertions, mocked filesystem and shell commands
**Target Platform**: Cross-platform CLI (Windows, macOS, Linux). Primary: single-file Windows .exe
**Project Type**: Single project (CLI tool)
**Performance Goals**: Interactive wizard <100ms response time, full bootstrap operation <10 seconds (excluding network)
**Constraints**: No external runtime dependencies, telemetry-free, secrets never persisted
**Scale/Scope**: Single-user local tool, processes ~10-20 config questions, generates 5-10 files per run

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**I. Simplicity First (KISS/DRY)**: ✅ PASS
- CLI framework (Cobra/urfave) provides battle-tested abstraction
- Template-based doc generation avoids duplication
- No premature abstractions planned

**II. Least Privilege**: ✅ PASS
- Reads credentials from env vars only, never writes them
- Creates files in user-specified directories only
- No elevated privileges required

**III. Reproducible Builds**: ✅ PASS
- Go modules with go.mod/go.sum for dependency locking
- Makefile documents exact build commands
- Static binary compilation ensures consistency

**IV. Clear Naming**: ✅ PASS
- Subcommands: `new`, `doctor`, `version` are self-explanatory
- Internal packages follow Go conventions (cmd/, internal/)

**V. Cross-Platform Compatibility**: ✅ PASS
- Go's cross-compilation handles Windows/macOS/Linux
- Using filepath.Join for path handling
- Primary target: Windows .exe as specified

**VI. Deterministic Toolchain**: ✅ PASS
- Requires Go 1.22+ (documented in README)
- Version check via `go version` command
- Build instructions in docs/building.md

**VII. Security Hardening**: ✅ PASS
- No secrets persisted to disk
- Reads GITHUB_TOKEN from env only
- Input validation for repo names, paths, license choices
- Errors don't expose sensitive data

**VIII. Minimal Observability**: ✅ PASS
- Structured JSON logging with `--debug` flag
- Default output: minimal, actionable
- Operation summaries: "Created X files, pushed to Y"

**IX. Documentation-First Culture**: ✅ PASS
- README with 90-second quickstart
- All docs in /docs directory
- Generated primer.md guides agent workflow

**X. Test Discipline**: ✅ PASS
- Unit tests for validation, templating, git operations
- End-to-end dry-run test (no remote side-effects)
- CI runs tests on every PR

**Verdict**: All constitutional requirements satisfied. No violations to justify.

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
cmd/
└── repo-bootstrapper/
    └── main.go              # CLI entry point

internal/
├── wizard/
│   ├── prompts.go           # Q&A prompts and validation
│   └── answers.go           # WizardAnswers struct and methods
├── workspace/
│   ├── creator.go           # Directory and file creation
│   └── templates.go         # Document template rendering
├── git/
│   ├── init.go              # Git initialization
│   └── github.go            # GitHub CLI integration
├── speckit/
│   └── integration.go       # Spec Kit CLI integration
└── logger/
    └── logger.go            # Structured logging

tests/
├── integration/
│   └── dryrun_test.go       # End-to-end dry-run test
└── unit/
    ├── wizard_test.go       # Input validation tests
    ├── workspace_test.go    # File creation tests
    └── git_test.go          # Git operations tests (mocked)

docs/
├── architecture.md          # System design overview
├── building.md              # Build and release process
└── templates/
    ├── constitution.tmpl    # Constitution template
    ├── spec.tmpl            # Spec template
    ├── plan.tmpl            # Plan template
    └── primer.tmpl          # Primer template

.github/
└── workflows/
    ├── ci.yml               # Lint, build, test on PR
    └── release.yml          # Build and publish on tag

Makefile                      # Build automation
go.mod                        # Go module definition
go.sum                        # Dependency lock
.gitignore                    # Git ignore patterns
README.md                     # 90-second quickstart
```

**Structure Decision**: Single-project Go CLI tool following standard Go layout. `cmd/` contains the executable entry point, `internal/` contains private application logic organized by domain (wizard, workspace, git, speckit, logger). Tests are separated by type (unit vs integration). Documentation and templates stored in `docs/`. GitHub Actions for CI/CD.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (data-model.md, cli-interface.md, quickstart.md)
- Group by Go packages: cmd/, internal/wizard/, internal/workspace/, internal/git/, internal/speckit/, internal/logger/
- CLI contract → CLI test tasks (one per subcommand)
- Each entity in data-model → struct definition + validation tests [P]
- Integration test from quickstart scenarios
- TDD approach: Tests before implementation

**Ordering Strategy**:
1. **Setup Phase**: Project init, go.mod, directory structure, Makefile
2. **Test Phase (TDD)**:
   - Unit tests for wizard validation [P]
   - Unit tests for workspace operations [P]
   - Unit tests for git/GitHub mocking [P]
   - Integration test (dry-run end-to-end)
3. **Implementation Phase**:
   - Data structures (internal/wizard/answers.go, etc.) [P]
   - Core logic (validation, templating, git ops) [P - different packages]
   - CLI commands (cmd/repo-bootstrapper/main.go + subcommands)
   - Integration (wire everything together)
4. **Polish Phase**:
   - Documentation (README, docs/architecture.md, docs/building.md) [P]
   - CI workflows (.github/workflows/) [P]
   - Makefile targets (build, test, release)

**Estimated Output**: 35-40 numbered, ordered tasks in tasks.md

**Parallelization Opportunities**:
- Multiple package implementations (wizard, workspace, git, logger) can be developed in parallel
- Test files for different packages can be written in parallel
- Documentation files can be written in parallel
- Mark with [P] when tasks touch different files/packages with no dependencies

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none required)

**Artifacts Generated**:
- [x] plan.md (this file)
- [x] research.md
- [x] data-model.md
- [x] contracts/cli-interface.md
- [x] quickstart.md
- [x] CLAUDE.md (agent context)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
