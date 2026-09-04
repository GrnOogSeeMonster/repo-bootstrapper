# Data Model: Repo Bootstrapper CLI

## Entity: WizardAnswers

**Purpose**: Captures all user inputs from the interactive Q&A wizard or loaded from JSON config file.

**Fields**:
```go
type WizardAnswers struct {
    // Repository metadata
    RepoOwner      string   `json:"repo_owner" validate:"required,alphanum_hyphen,max=39"`
    RepoName       string   `json:"repo_name" validate:"required,repo_name,max=100"`
    Visibility     string   `json:"visibility" validate:"required,oneof=public private"`

    // Project configuration
    License        string   `json:"license" validate:"required,oneof=MIT Apache-2.0 GPL-3.0 BSD-3-Clause MPL-2.0 ISC Unlicense None"`
    Language       string   `json:"language" validate:"required"`
    Framework      string   `json:"framework"`
    TargetOS       []string `json:"target_os" validate:"required,min=1"`

    // Tooling
    AgentChoice    string   `json:"agent_choice" validate:"required,oneof=cursor claude copilot gemini"`
    PackageManager string   `json:"package_manager"`
    TestFramework  string   `json:"test_framework"`

    // Application type
    AppKind        string   `json:"app_kind" validate:"required,oneof=CLI web service"`

    // GitHub features
    IssueLabels    string   `json:"issue_labels" validate:"oneof=standard extended custom none"`
    CustomLabels   []Label  `json:"custom_labels,omitempty"`

    // Workspace
    WorkspacePath  string   `json:"workspace_path" validate:"required,dirpath"`

    // Optional
    InitialFeatures []string `json:"initial_features,omitempty"`
}
```

**Validation Rules**:
- `repo_owner`: Alphanumeric + hyphens only, max 39 chars (GitHub limit)
- `repo_name`: Alphanumeric + dots, underscores, hyphens, max 100 chars
- `visibility`: Exactly "public" or "private"
- `license`: Must be in allowlist (see clarifications)
- `target_os`: At least one of: Windows, macOS, Linux
- `agent_choice`: One of: cursor, claude, copilot, gemini
- `app_kind`: One of: CLI, web, service
- `workspace_path`: Valid directory path (cleaned with filepath.Clean)

**Relationships**: None (value object)

**State Transitions**: Immutable once collected. No transitions.

---

## Entity: WorkspaceConfig

**Purpose**: Represents the configuration and state of the bootstrapped workspace.

**Fields**:
```go
type WorkspaceConfig struct {
    RootPath       string
    GitInitialized bool
    SpecKitReady   bool
    DocsGenerated  []string  // Paths to generated docs
    CreatedFiles   []string  // All created file paths
    GitRemoteURL   string
}
```

**Relationships**:
- Derived from `WizardAnswers`
- Used by `workspace.Creator` to track progress

**State Transitions**:
1. **Initial**: Empty struct
2. **DirectoryCreated**: RootPath set
3. **GitInitialized**: GitInitialized = true
4. **SpecKitReady**: SpecKitReady = true
5. **DocsGenerated**: DocsGenerated populated
6. **RemoteCreated**: GitRemoteURL set

---

## Entity: DocumentationArtifact

**Purpose**: Represents a generated documentation file with its content and metadata.

**Fields**:
```go
type DocumentationArtifact struct {
    FilePath     string
    TemplateName string
    Content      string
    Variables    map[string]interface{}
    Generated    time.Time
}
```

**Types**:
- `docs/constitution.md`: Principles tailored to target app
- `docs/spec.md`: What/why for target app
- `docs/plan.md`: Tech choices and architecture
- `docs/tasks.md`: Placeholder header
- `prompts/primer.md`: Agent workflow instructions
- `README.md`: Project introduction
- `LICENSE`: License text (if not None)

**Relationships**:
- Rendered from templates using `WizardAnswers`
- Stored in `WorkspaceConfig.DocsGenerated`

---

## Entity: OperationLog

**Purpose**: Records all filesystem and GitHub operations for audit trail and dry-run display.

**Fields**:
```go
type OperationLog struct {
    Operations []Operation
}

type Operation struct {
    Type      OperationType  // CREATE_DIR, CREATE_FILE, GIT_INIT, GH_CREATE, etc.
    Target    string         // Path or resource name
    Details   string         // Human-readable description
    DryRun    bool           // Was this a dry-run?
    Timestamp time.Time
    Error     error          // Nil if successful
}

type OperationType int

const (
    OpCreateDir OperationType = iota
    OpCreateFile
    OpWriteFile
    OpGitInit
    OpGitAdd
    OpGitCommit
    OpGitRemoteAdd
    OpGitPush
    OpGHRepoCreate
    OpGHLabelCreate
    OpSpecKitInit
)
```

**Relationships**:
- Populated by all workspace, git, and GitHub operations
- Consumed by logger for summary output

**State Transitions**: Append-only log. Operations never modified after recording.

---

## Entity: Label

**Purpose**: Represents a GitHub issue label configuration.

**Fields**:
```go
type Label struct {
    Name        string `json:"name" validate:"required,max=50"`
    Description string `json:"description" validate:"max=100"`
    Color       string `json:"color" validate:"required,hexcolor"`
}
```

**Predefined Label Sets**:

**Standard Set**:
```json
[
  {"name": "bug", "description": "Something isn't working", "color": "d73a4a"},
  {"name": "enhancement", "description": "New feature or request", "color": "a2eeef"},
  {"name": "documentation", "description": "Improvements or additions to documentation", "color": "0075ca"},
  {"name": "question", "description": "Further information is requested", "color": "d876e3"}
]
```

**Extended Set** (includes standard + app-kind specific):
- CLI apps add: `cli`, `user-experience`
- Web apps add: `frontend`, `backend`, `api`
- Service apps add: `performance`, `reliability`

**Validation**:
- Name: Required, max 50 chars, no special chars except hyphen/underscore
- Color: 6-digit hex color (without #)
- Description: Max 100 chars

---

## Value Objects

### GitConfig
```go
type GitConfig struct {
    DefaultBranch string  // "main"
    UserName      string  // From git config or prompt
    UserEmail     string  // From git config or prompt
}
```

### TemplateVars
```go
type TemplateVars struct {
    ProjectName    string
    RepoOwner      string
    License        string
    Language       string
    AgentChoice    string
    AppKind        string
    CurrentDate    string
    BranchStrategy string  // "Git Flow: main, develop, feature/*, release/*, hotfix/*"
}
```

---

## Validation Summary

All entities with user input implement validation using `go-playground/validator`:

**Custom Validators**:
- `repo_name`: Regex `^[a-zA-Z0-9._-]+$`, no leading/trailing hyphen
- `alphanum_hyphen`: Regex `^[a-zA-Z0-9-]+$`
- `dirpath`: Cleaned absolute path, no directory traversal
- `hexcolor`: Regex `^[0-9A-Fa-f]{6}$`

**Validation Timing**:
- Interactive mode: Validate each answer before accepting
- Non-interactive mode: Validate entire WizardAnswers JSON before any operations
- Dry-run mode: Validation still performed (fail fast)

---

## Persistence Strategy

**No database required**. All state is transient:
- `WizardAnswers`: Collected in memory, optionally read from JSON file
- `WorkspaceConfig`: Tracked in memory during execution
- `OperationLog`: Accumulated in memory, output to stdout/logs at end
- `DocumentationArtifacts`: Generated and written directly to filesystem

**No serialization** of answers to disk unless user explicitly saves as JSON for future non-interactive runs.
