package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
)

// Creator handles workspace creation
type Creator struct {
	config *WorkspaceConfig
	dryRun bool
	opLog  *logger.OperationLog
}

// NewCreator creates a new workspace creator
func NewCreator(rootPath string, dryRun bool, opLog *logger.OperationLog) *Creator {
	return &Creator{
		config: NewWorkspaceConfig(rootPath),
		dryRun: dryRun,
		opLog:  opLog,
	}
}

// Create creates the workspace directory and structure
func (c *Creator) Create() error {
	// Check if directory exists
	if _, err := os.Stat(c.config.RootPath); err == nil {
		// Directory exists - prompt user
		if c.dryRun {
			logger.Info("[DRY RUN] Directory exists, would prompt for overwrite", "path", c.config.RootPath)
		} else {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Directory %s already exists. Overwrite", c.config.RootPath),
				IsConfirm: true,
			}
			_, err := prompt.Run()
			if err != nil {
				return fmt.Errorf("cancelled: directory already exists")
			}
		}
	}

	// Create directory
	if c.dryRun {
		logger.Info("[DRY RUN] Would create directory", "path", c.config.RootPath)
		c.opLog.Add(logger.Operation{
			Type:    logger.OpCreateDir,
			Target:  c.config.RootPath,
			Details: "Create workspace directory",
			DryRun:  true,
		})
	} else {
		if err := os.MkdirAll(c.config.RootPath, 0755); err != nil {
			c.opLog.Add(logger.Operation{
				Type:    logger.OpCreateDir,
				Target:  c.config.RootPath,
				Details: "Create workspace directory",
				Error:   err,
			})
			return fmt.Errorf("failed to create directory: %w", err)
		}
		logger.Info("Created workspace directory", "path", c.config.RootPath)
		c.opLog.Add(logger.Operation{
			Type:    logger.OpCreateDir,
			Target:  c.config.RootPath,
			Details: "Create workspace directory",
		})
		c.config.AddCreatedFile(c.config.RootPath)
	}

	// Create subdirectories
	subdirs := []string{"docs", "prompts"}
	for _, subdir := range subdirs {
		path := filepath.Join(c.config.RootPath, subdir)
		if c.dryRun {
			logger.Info("[DRY RUN] Would create directory", "path", path)
			c.opLog.Add(logger.Operation{
				Type:    logger.OpCreateDir,
				Target:  path,
				Details: fmt.Sprintf("Create %s directory", subdir),
				DryRun:  true,
			})
		} else {
			if err := os.MkdirAll(path, 0755); err != nil {
				return fmt.Errorf("failed to create %s directory: %w", subdir, err)
			}
			c.opLog.Add(logger.Operation{
				Type:    logger.OpCreateDir,
				Target:  path,
				Details: fmt.Sprintf("Create %s directory", subdir),
			})
			c.config.AddCreatedFile(path)
		}
	}

	return nil
}

// GetConfig returns the workspace configuration
func (c *Creator) GetConfig() *WorkspaceConfig {
	return c.config
}
