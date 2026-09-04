package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
)

// GitHub handles GitHub CLI operations
type GitHub struct {
	dryRun bool
	opLog  *logger.OperationLog
}

// NewGitHub creates a new GitHub handler
func NewGitHub(dryRun bool, opLog *logger.OperationLog) *GitHub {
	return &GitHub{
		dryRun: dryRun,
		opLog:  opLog,
	}
}

// IsAuthenticated checks if gh CLI is authenticated
func (g *GitHub) IsAuthenticated() bool {
	cmd := exec.Command("gh", "auth", "status")
	err := cmd.Run()
	return err == nil
}

// IsInstalled checks if gh CLI is installed
func (g *GitHub) IsInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// CreateRepo creates a GitHub repository
func (g *GitHub) CreateRepo(owner, name, visibility, workspacePath string) error {
	if !g.IsInstalled() {
		logger.Info("GitHub CLI not installed, skipping repo creation")
		logger.Info("To create repo manually, run: gh repo create %s/%s --%s --source %s --remote origin --push", owner, name, visibility, workspacePath)
		return nil
	}

	if !g.IsAuthenticated() {
		logger.Info("GitHub CLI not authenticated, skipping repo creation")
		logger.Info("Run 'gh auth login' to authenticate, then manually create repo")
		return nil
	}

	visFlag := "--" + visibility
	cmdStr := fmt.Sprintf("gh repo create %s/%s %s --source %s --remote origin --push", owner, name, visFlag, workspacePath)

	if g.dryRun {
		logger.Info("[DRY RUN] Would run: " + cmdStr)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpGHRepoCreate,
			Target:  fmt.Sprintf("%s/%s", owner, name),
			Details: cmdStr,
			DryRun:  true,
		})
		return nil
	}

	cmd := exec.Command("gh", "repo", "create", fmt.Sprintf("%s/%s", owner, name), visFlag, "--source", workspacePath, "--remote", "origin", "--push")
	cmd.Dir = workspacePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create GitHub repo: %v\nOutput: %s", err, string(output))
		logger.Error(errMsg, err)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpGHRepoCreate,
			Target:  fmt.Sprintf("%s/%s", owner, name),
			Details: cmdStr,
			Error:   fmt.Errorf("%s", errMsg),
		})
		return fmt.Errorf("%s", errMsg)
	}

	logger.Info("Created GitHub repository", "repo", fmt.Sprintf("%s/%s", owner, name))
	g.opLog.Add(logger.Operation{
		Type:    logger.OpGHRepoCreate,
		Target:  fmt.Sprintf("%s/%s", owner, name),
		Details: cmdStr,
	})

	return nil
}

// CreateLabels creates issue labels
func (g *GitHub) CreateLabels(labelType, workspacePath string) error {
	if !g.IsInstalled() || !g.IsAuthenticated() {
		logger.Info("Skipping label creation (gh CLI not available)")
		return nil
	}

	if labelType == "none" {
		return nil
	}

	labels := getStandardLabels()
	if labelType == "extended" {
		labels = append(labels, getExtendedLabels()...)
	}

	for _, label := range labels {
		cmdStr := fmt.Sprintf("gh label create %s --color %s --description \"%s\"", label.Name, label.Color, label.Description)

		if g.dryRun {
			logger.Info("[DRY RUN] Would run: " + cmdStr)
			g.opLog.Add(logger.Operation{
				Type:    logger.OpGHLabelCreate,
				Target:  label.Name,
				Details: cmdStr,
				DryRun:  true,
			})
			continue
		}

		cmd := exec.Command("gh", "label", "create", label.Name, "--color", label.Color, "--description", label.Description)
		cmd.Dir = workspacePath

		// Ignore errors (label might already exist)
		output, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(output), "already exists") {
			logger.Debug("Label creation warning", "label", label.Name, "output", string(output))
		}

		g.opLog.Add(logger.Operation{
			Type:    logger.OpGHLabelCreate,
			Target:  label.Name,
			Details: cmdStr,
		})
	}

	logger.Info("Created issue labels", "type", labelType)
	return nil
}

// Label represents a GitHub issue label
type Label struct {
	Name        string
	Description string
	Color       string
}

func getStandardLabels() []Label {
	return []Label{
		{Name: "bug", Description: "Something isn't working", Color: "d73a4a"},
		{Name: "enhancement", Description: "New feature or request", Color: "a2eeef"},
		{Name: "documentation", Description: "Improvements or additions to documentation", Color: "0075ca"},
		{Name: "question", Description: "Further information is requested", Color: "d876e3"},
	}
}

func getExtendedLabels() []Label {
	return []Label{
		{Name: "performance", Description: "Performance improvement", Color: "ffaa00"},
		{Name: "security", Description: "Security issue or improvement", Color: "ff0000"},
		{Name: "technical-debt", Description: "Code refactoring needed", Color: "cccccc"},
	}
}
