package wizard

import (
	"github.com/go-playground/validator/v10"
)

// WizardAnswers captures all user inputs from Q&A or config file
type WizardAnswers struct {
	RepoOwner       string   `json:"repo_owner" validate:"required,max=39"`
	RepoName        string   `json:"repo_name" validate:"required,max=100"`
	Visibility      string   `json:"visibility" validate:"required,oneof=public private"`
	License         string   `json:"license" validate:"required,oneof=MIT Apache-2.0 GPL-3.0 BSD-3-Clause MPL-2.0 ISC Unlicense None"`
	Language        string   `json:"language" validate:"required"`
	Framework       string   `json:"framework"`
	TargetOS        []string `json:"target_os" validate:"required,min=1"`
	AgentChoice     string   `json:"agent_choice" validate:"required,oneof=cursor claude copilot gemini"`
	PackageManager  string   `json:"package_manager"`
	TestFramework   string   `json:"test_framework"`
	AppKind         string   `json:"app_kind" validate:"required,oneof=CLI web service"`
	IssueLabels     string   `json:"issue_labels" validate:"oneof=standard extended custom none"`
	WorkspacePath   string   `json:"workspace_path" validate:"required"`
	InitialFeatures []string `json:"initial_features,omitempty"`
}

// Validate validates all fields using struct tags
func (w *WizardAnswers) Validate() error {
	validate := validator.New()
	return validate.Struct(w)
}
