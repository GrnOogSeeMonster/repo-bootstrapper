package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// T007: Test for WizardAnswers validation
func TestWizardAnswers_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid repo name with spaces",
			input: map[string]interface{}{
				"repo_name": "invalid name!",
			},
			wantErr: true,
			errMsg:  "repo_name",
		},
		{
			name: "missing required field",
			input: map[string]interface{}{
				"repo_name": "valid-name",
				// missing repo_owner
			},
			wantErr: true,
			errMsg:  "repo_owner",
		},
		{
			name: "invalid license not in allowlist",
			input: map[string]interface{}{
				"repo_owner": "myorg",
				"repo_name":  "myrepo",
				"license":    "InvalidLicense",
			},
			wantErr: true,
			errMsg:  "license",
		},
		{
			name: "repo owner exceeds 39 chars",
			input: map[string]interface{}{
				"repo_owner": "this-is-way-too-long-for-github-owner-limit-test",
			},
			wantErr: true,
			errMsg:  "max=39",
		},
		{
			name: "valid input",
			input: map[string]interface{}{
				"repo_owner":   "myorg",
				"repo_name":    "myrepo",
				"visibility":   "public",
				"license":      "MIT",
				"language":     "Go",
				"target_os":    []string{"Linux"},
				"agent_choice": "claude",
				"app_kind":     "CLI",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This will fail until WizardAnswers validation is implemented
			t.Skip("Implementation pending - TDD placeholder")

			// TODO: Implement validation test
			// answers := wizard.NewWizardAnswers(tt.input)
			// err := answers.Validate()
			// if tt.wantErr {
			// 	assert.Error(t, err)
			// 	assert.Contains(t, err.Error(), tt.errMsg)
			// } else {
			// 	assert.NoError(t, err)
			// }
		})
	}
}

func TestRepoNameValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid alphanumeric", "myrepo", false},
		{"valid with hyphens", "my-repo", false},
		{"valid with dots", "my.repo", false},
		{"valid with underscores", "my_repo", false},
		{"invalid with spaces", "my repo", true},
		{"invalid special chars", "my@repo", true},
		{"invalid leading hyphen", "-myrepo", true},
		{"invalid trailing hyphen", "myrepo-", true},
		{"exceeds 100 chars", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("Implementation pending - TDD placeholder")
			// TODO: Implement validation
			// err := wizard.ValidateRepoName(tt.input)
			// if tt.wantErr {
			// 	assert.Error(t, err)
			// } else {
			// 	assert.NoError(t, err)
			// }
		})
	}
}

func TestLicenseAllowlist(t *testing.T) {
	validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "MPL-2.0", "ISC", "Unlicense", "None"}

	for _, license := range validLicenses {
		t.Run("valid_"+license, func(t *testing.T) {
			t.Skip("Implementation pending - TDD placeholder")
			// TODO: Implement validation
			// err := wizard.ValidateLicense(license)
			// assert.NoError(t, err)
		})
	}

	t.Run("invalid license", func(t *testing.T) {
		t.Skip("Implementation pending - TDD placeholder")
		// TODO: Implement validation
		// err := wizard.ValidateLicense("InvalidLicense")
		// assert.Error(t, err)
	})
}
