package speckit

import (
	"fmt"
	"os/exec"

	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
)

// Integration handles Spec Kit CLI operations
type Integration struct {
	dryRun bool
	opLog  *logger.OperationLog
}

// NewIntegration creates a new Spec Kit integration handler
func NewIntegration(dryRun bool, opLog *logger.OperationLog) *Integration {
	return &Integration{
		dryRun: dryRun,
		opLog:  opLog,
	}
}

// IsInstalled checks if specify CLI is installed
func (s *Integration) IsInstalled() bool {
	_, err := exec.LookPath("specify")
	return err == nil
}

// Init runs specify init in the workspace
func (s *Integration) Init(projectName, agent, workspacePath string) error {
	if !s.IsInstalled() {
		logger.Info("Spec Kit CLI not installed, skipping initialization")
		logger.Info("Install with: npm install -g @specify/cli")
		logger.Info("Then manually run: specify init %s --ai %s", projectName, agent)
		return nil
	}

	cmdStr := fmt.Sprintf("specify init %s --ai %s", projectName, agent)

	if s.dryRun {
		logger.Info("[DRY RUN] Would run: " + cmdStr)
		s.opLog.Add(logger.Operation{
			Type:    logger.OpSpecKitInit,
			Target:  workspacePath,
			Details: cmdStr,
			DryRun:  true,
		})
		return nil
	}

	cmd := exec.Command("specify", "init", projectName, "--ai", agent)
	cmd.Dir = workspacePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to initialize Spec Kit: %v\nOutput: %s", err, string(output))
		logger.Error(errMsg, err)
		s.opLog.Add(logger.Operation{
			Type:    logger.OpSpecKitInit,
			Target:  workspacePath,
			Details: cmdStr,
			Error:   fmt.Errorf("%s", errMsg),
		})
		// Don't fail the entire operation if Spec Kit fails
		logger.Info("Continuing without Spec Kit (can be run manually later)")
		return nil
	}

	logger.Info("Initialized Spec Kit", "project", projectName, "agent", agent)
	s.opLog.Add(logger.Operation{
		Type:    logger.OpSpecKitInit,
		Target:  workspacePath,
		Details: cmdStr,
	})

	return nil
}
