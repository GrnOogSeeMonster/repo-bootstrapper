package wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// RunInteractiveWizard runs the interactive Q&A wizard
func RunInteractiveWizard() (*WizardAnswers, error) {
	answers := &WizardAnswers{}

	// Repo Owner
	ownerPrompt := promptui.Prompt{
		Label: "Repository Owner/Organization",
		Validate: func(input string) error {
			return ValidateRepoOwner(input)
		},
	}
	owner, err := ownerPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.RepoOwner = owner

	// Repo Name
	namePrompt := promptui.Prompt{
		Label: "Repository Name",
		Validate: func(input string) error {
			return ValidateRepoName(input)
		},
	}
	name, err := namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.RepoName = name

	// Visibility
	visibilitySelect := promptui.Select{
		Label: "Repository Visibility",
		Items: []string{"public", "private"},
	}
	_, visibility, err := visibilitySelect.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.Visibility = visibility

	// License
	licenseSelect := promptui.Select{
		Label: "License",
		Items: []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "MPL-2.0", "ISC", "Unlicense", "None"},
	}
	_, license, err := licenseSelect.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.License = license

	// Language
	languagePrompt := promptui.Prompt{
		Label:   "Primary Language",
		Default: "Go",
	}
	language, err := languagePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.Language = language

	// Framework (optional)
	frameworkPrompt := promptui.Prompt{
		Label:   "Framework (optional, press Enter to skip)",
		Default: "",
	}
	framework, _ := frameworkPrompt.Run()
	answers.Framework = framework

	// Target OS
	osPrompt := promptui.Select{
		Label: "Target Operating Systems (select primary, others assumed)",
		Items: []string{"Linux", "Windows", "macOS", "All"},
	}
	_, osChoice, err := osPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	if osChoice == "All" {
		answers.TargetOS = []string{"Windows", "macOS", "Linux"}
	} else {
		answers.TargetOS = []string{osChoice}
	}

	// Agent Choice
	agentSelect := promptui.Select{
		Label: "AI Coding Agent",
		Items: []string{"claude", "cursor", "copilot", "gemini"},
	}
	_, agent, err := agentSelect.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.AgentChoice = agent

	// Package Manager (optional)
	pkgMgrPrompt := promptui.Prompt{
		Label:   "Package Manager (optional)",
		Default: "",
	}
	pkgMgr, _ := pkgMgrPrompt.Run()
	answers.PackageManager = pkgMgr

	// Test Framework (optional)
	testFwPrompt := promptui.Prompt{
		Label:   "Test Framework (optional)",
		Default: "",
	}
	testFw, _ := testFwPrompt.Run()
	answers.TestFramework = testFw

	// App Kind
	appKindSelect := promptui.Select{
		Label: "Application Kind",
		Items: []string{"CLI", "web", "service"},
	}
	_, appKind, err := appKindSelect.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.AppKind = appKind

	// Issue Labels
	labelsSelect := promptui.Select{
		Label: "GitHub Issue Labels",
		Items: []string{"standard", "extended", "none"},
	}
	_, labels, err := labelsSelect.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	answers.IssueLabels = labels

	// Workspace Path
	home, _ := os.UserHomeDir()
	defaultPath := filepath.Join(home, "Documents", answers.RepoName)

	pathPrompt := promptui.Prompt{
		Label:   "Workspace Path",
		Default: defaultPath,
		Validate: func(input string) error {
			_, err := ValidatePath(input)
			return err
		},
	}
	path, err := pathPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("cancelled")
	}
	cleanedPath, _ := ValidatePath(path)
	answers.WorkspacePath = cleanedPath

	// Initial Features (optional)
	featuresPrompt := promptui.Prompt{
		Label:   "Initial Features (comma-separated, optional)",
		Default: "",
	}
	features, _ := featuresPrompt.Run()
	if features != "" {
		answers.InitialFeatures = strings.Split(features, ",")
		for i := range answers.InitialFeatures {
			answers.InitialFeatures[i] = strings.TrimSpace(answers.InitialFeatures[i])
		}
	}

	return answers, nil
}
