<!--
Sync Impact Report:
Version: 0.0.0 → 1.0.0 (Initial constitution)

Modified principles: N/A (new constitution)
Added sections:
  - Core Principles (10 principles)
  - Security & Trust
  - Development Workflow
  - Governance

Templates requiring updates:
  ✅ .specify/templates/plan-template.md (reviewed - compatible)
  ✅ .specify/templates/spec-template.md (reviewed - compatible)
  ✅ .specify/templates/tasks-template.md (reviewed - compatible)

Follow-up TODOs: None
-->

# Repo Bootstrapper Constitution

## Core Principles

### I. Simplicity First (KISS/DRY)
**MUST** prefer simple, maintainable solutions over clever abstractions. Every abstraction MUST justify its existence by eliminating significant duplication or complexity. **MUST** apply DRY only when duplication creates maintenance burden, not reflexively.

**Rationale**: Simple code is debuggable, auditable, and maintainable by others. Premature abstraction creates more problems than it solves.

### II. Least Privilege
**MUST** operate with minimum required permissions. **MUST** never request elevated privileges unless absolutely necessary and documented. **MUST** fail safe - prefer operation failure over security compromise.

**Rationale**: Security through constraint. Limiting scope of access limits scope of damage.

### III. Reproducible Builds
**MUST** pin all dependencies with exact versions (lock files). **MUST** document toolchain versions. **MUST** ensure builds are deterministic - same inputs produce identical outputs.

**Rationale**: Reproducibility enables debugging, security auditing, and reliable distribution.

### IV. Clear Naming
**MUST** use descriptive, unambiguous names for functions, variables, files, and modules. Avoid abbreviations unless universally recognized. Names **MUST** reveal intent.

**Rationale**: Code is read far more than written. Clear names reduce cognitive load and prevent errors.

### V. Cross-Platform Compatibility
**MUST** support Windows, macOS, and Linux. **MUST** use portable APIs and avoid platform-specific assumptions. Primary deployment target is a single self-contained Windows .exe that runs without external dependencies.

**Rationale**: Tool must work everywhere developers work. Windows .exe simplifies distribution and reduces friction.

### VI. Deterministic Toolchain
**MUST** specify exact versions of compilers, build tools, and runtime environments. **MUST** document installation steps. **MUST** provide version check commands.

**Rationale**: "Works on my machine" is not acceptable. Explicit toolchain enables collaboration and debugging.

### VII. Security Hardening
**MUST** never persist secrets to disk, logs, or environment files. **MUST** read GitHub tokens and credentials exclusively from environment variables at runtime. **MUST** validate all external inputs (user input, API responses, file contents). **MUST** handle errors safely - no sensitive data in error messages.

**Rationale**: Secrets on disk are compromised secrets. Defense in depth requires input validation and safe failure modes.

### VIII. Minimal Observability
**MUST** implement structured logging (JSON or key=value format). **MUST** support a `--debug` CLI flag that enables verbose output. Default output **MUST** be minimal and actionable. **MUST NOT** log sensitive data even in debug mode.

**Rationale**: Too much logging is noise; too little is blindness. Structured logs enable parsing and analysis.

### IX. Documentation-First Culture
**MUST** maintain all documentation in `/docs` directory. **MUST** provide a single-page quickstart guide. **MUST** ensure all documentation is printable (avoid infinite scrolls, prefer markdown). **MUST** document decisions, not just outcomes.

**Rationale**: Undocumented features don't exist. Discoverable docs in one place reduce cognitive load.

### X. Test Discipline
**MUST** write unit tests for core business logic. **MUST** provide at least one smoke test that validates the CLI happy path (install → run → verify output). **MUST** run tests in CI before merge. Test coverage **SHOULD** be measured but **MUST NOT** be gamed with meaningless tests.

**Rationale**: Tests document expected behavior and catch regressions. Smoke tests verify distribution integrity.

## Security & Trust

### Secrets Management
- **MUST** read all credentials from environment variables (e.g., `GITHUB_TOKEN`, `GH_TOKEN`)
- **MUST NOT** accept secrets via CLI arguments (visible in process lists)
- **MUST NOT** persist secrets to configuration files
- **MUST** clear sensitive data from memory when no longer needed
- **MUST** fail gracefully with clear error messages when credentials are missing

### Input Validation
- **MUST** validate and sanitize all user inputs before processing
- **MUST** validate all external API responses before trusting
- **MUST** use allowlists over denylists for validation
- **MUST** reject unexpected input rather than coercing it

### Error Handling
- **MUST NOT** expose stack traces or internal paths to end users
- **MUST** log detailed errors internally (when `--debug` enabled)
- **MUST** provide actionable error messages (what went wrong, how to fix)
- **MUST** fail fast and loudly rather than silently corrupting state

## Development Workflow

### Agent Collaboration Guidelines
- **MUST** prefer battle-tested, simple libraries over novel or complex ones
- **MUST** avoid overengineering - solve the problem at hand, not imagined future problems
- **MUST** keep AI prompts and agent instructions in `/prompts` directory
- **MUST** version prompts alongside code when prompt changes affect behavior
- **SHOULD** document prompt engineering decisions in `/docs/prompts.md`

### Auditability & Transparency
- **MUST** log all file creation, modification, and deletion operations with paths
- **MUST** document the "why" for each operation (reference issue, user request, or design decision)
- **MUST** support `--dry-run` mode that shows what would be done without doing it
- **MUST** generate operation summaries (e.g., "Created 5 files, modified 2, deleted 0")

### Testing Requirements
- **MUST** write unit tests for:
  - Input validation logic
  - Business rule enforcement
  - Data transformations
  - Error handling paths
- **MUST** provide a smoke test that:
  - Installs/runs the tool
  - Executes a representative command
  - Validates expected output or side effects
- **MUST** run all tests in CI on pull requests
- **SHOULD** measure test coverage but **MUST NOT** enforce arbitrary thresholds

### Build & Distribution
- Primary target: Single-file Windows .exe with no external dependencies
- Secondary targets: macOS binary, Linux binary
- **MUST** provide checksums (SHA256) for all distributed binaries
- **MUST** document build process in `/docs/building.md`
- **MUST** automate release builds in CI

## Governance

### Amendment Process
1. Propose amendment via pull request to `constitution.md`
2. Document rationale: what problem does the change solve?
3. Identify affected code/docs and plan migration if backward-incompatible
4. Increment version (see Versioning below)
5. Merge only after review and approval

### Versioning
**CONSTITUTION_VERSION** follows semantic versioning:
- **MAJOR**: Backward-incompatible governance changes (removed or redefined principles)
- **MINOR**: New principles, new sections, or material expansions of guidance
- **PATCH**: Clarifications, typo fixes, wording improvements without semantic change

### Compliance & Review
- All pull requests **MUST** pass constitution compliance check
- Reviewers **MUST** verify adherence to principles
- Complexity that violates principles **MUST** be justified in writing or refactored
- Constitution supersedes all other documentation in case of conflict

### Living Document
This constitution is a living document. When reality conflicts with principle, we must:
1. Understand why the conflict exists
2. Determine if the principle is wrong or the approach is wrong
3. Update the constitution or change the approach accordingly
4. Document the decision for future reference

**Version**: 1.0.0 | **Ratified**: 2025-10-01 | **Last Amended**: 2025-10-01
