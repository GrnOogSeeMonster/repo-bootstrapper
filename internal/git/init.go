package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/workspace"
)

// Repository handles git operations
type Repository struct {
	config *workspace.WorkspaceConfig
	repo   *git.Repository
	dryRun bool
	opLog  *logger.OperationLog
}

// NewRepository creates a new git repository handler
func NewRepository(config *workspace.WorkspaceConfig, dryRun bool, opLog *logger.OperationLog) *Repository {
	return &Repository{
		config: config,
		dryRun: dryRun,
		opLog:  opLog,
	}
}

// Init initializes a git repository
func (r *Repository) Init() error {
	if r.dryRun {
		logger.Info("[DRY RUN] Would initialize git repository", "path", r.config.RootPath)
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitInit,
			Target:  r.config.RootPath,
			Details: "Initialize git repository",
			DryRun:  true,
		})
		return nil
	}

	repo, err := git.PlainInit(r.config.RootPath, false)
	if err != nil {
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitInit,
			Target:  r.config.RootPath,
			Details: "Initialize git repository",
			Error:   err,
		})
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	r.repo = repo
	r.config.GitInitialized = true

	logger.Info("Initialized git repository")
	r.opLog.Add(logger.Operation{
		Type:    logger.OpGitInit,
		Target:  r.config.RootPath,
		Details: "Initialize git repository",
	})

	return nil
}

// AddAll adds all files to git staging
func (r *Repository) AddAll() error {
	if r.dryRun {
		logger.Info("[DRY RUN] Would run: git add .")
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitAdd,
			Target:  r.config.RootPath,
			Details: "Add all files to git",
			DryRun:  true,
		})
		return nil
	}

	if r.repo == nil {
		return fmt.Errorf("repository not initialized")
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add all files
	if err := w.AddGlob("."); err != nil {
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitAdd,
			Target:  r.config.RootPath,
			Details: "Add all files to git",
			Error:   err,
		})
		return fmt.Errorf("failed to add files: %w", err)
	}

	logger.Info("Added files to git")
	r.opLog.Add(logger.Operation{
		Type:    logger.OpGitAdd,
		Target:  r.config.RootPath,
		Details: "Add all files to git",
	})

	return nil
}

// Commit creates an initial commit
func (r *Repository) Commit(message string) error {
	if r.dryRun {
		logger.Info("[DRY RUN] Would commit", "message", message)
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitCommit,
			Target:  r.config.RootPath,
			Details: fmt.Sprintf("Commit: %s", message),
			DryRun:  true,
		})
		return nil
	}

	if r.repo == nil {
		return fmt.Errorf("repository not initialized")
	}

	w, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get git config for author info
	authorName := getGitConfig("user.name", "Repo Bootstrapper")
	authorEmail := getGitConfig("user.email", "noreply@example.com")

	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
		},
	})

	if err != nil {
		r.opLog.Add(logger.Operation{
			Type:    logger.OpGitCommit,
			Target:  r.config.RootPath,
			Details: fmt.Sprintf("Commit: %s", message),
			Error:   err,
		})
		return fmt.Errorf("failed to commit: %w", err)
	}

	logger.Info("Created git commit", "message", message)
	r.opLog.Add(logger.Operation{
		Type:    logger.OpGitCommit,
		Target:  r.config.RootPath,
		Details: fmt.Sprintf("Commit: %s", message),
	})

	return nil
}

func getGitConfig(key, defaultValue string) string {
	// Try to read from global git config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return defaultValue
	}

	gitConfigPath := filepath.Join(homeDir, ".gitconfig")
	if _, err := os.Stat(gitConfigPath); os.IsNotExist(err) {
		return defaultValue
	}

	// For simplicity, return default
	// A full implementation would parse .gitconfig
	return defaultValue
}
