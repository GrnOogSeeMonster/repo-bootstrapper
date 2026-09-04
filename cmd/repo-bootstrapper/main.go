package main

import (
	"fmt"
	"os"

	"github.com/repo-bootstrapper/repo-bootstrapper/internal/git"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/speckit"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/wizard"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	// Version info (set via ldflags during build)
	Version   = "dev"
	BuildTime = "unknown"
	Commit    = "unknown"

	// Global flags
	debugFlag bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "repo-bootstrapper",
		Short: "Bootstrap new repositories with best practices and AI agent integration",
		Long: `Repo Bootstrapper is a CLI tool that automates repository setup through an
interactive Q&A wizard. It creates workspaces, initializes git, sets up Spec Kit,
generates documentation, and publishes to GitHub.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.Init(debugFlag)
		},
	}

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")

	// Add subcommands
	rootCmd.AddCommand(newCommand())
	rootCmd.AddCommand(doctorCommand())
	rootCmd.AddCommand(versionCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	var (
		fromFile string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new bootstrapped repository",
		Long:  `Guides you through an interactive wizard to create a new repository with all necessary setup.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNew(fromFile, dryRun)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load configuration from JSON file (non-interactive mode)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview operations without executing them")

	return cmd
}

func doctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check prerequisites and environment",
		Long:  `Verifies that all required tools are installed and properly configured.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Checking prerequisites...")
			fmt.Println()

			// Check Go
			fmt.Println("✓ Go installed (repo-bootstrapper is running)")

			// TODO: Check other prerequisites
			fmt.Println("✗ Full prerequisite checking not yet implemented")
			fmt.Println()

			return nil
		},
	}
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("repo-bootstrapper %s\n", Version)
			fmt.Printf("Built: %s\n", BuildTime)
			fmt.Printf("Commit: %s\n", Commit)
		},
	}
}

func runNew(fromFile string, dryRun bool) error {
	opLog := &logger.OperationLog{Operations: make([]logger.Operation, 0)}

	logger.Info("Starting repo-bootstrapper")

	if dryRun {
		fmt.Println()
		fmt.Println("=== DRY RUN MODE ===")
		fmt.Println("No changes will be made to your filesystem or GitHub")
		fmt.Println()
	}

	// Step 1: Collect inputs
	var answers *wizard.WizardAnswers
	var err error

	if fromFile != "" {
		logger.Info("Loading configuration from file", "file", fromFile)
		answers, err = wizard.LoadFromFile(fromFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		logger.Info("Starting interactive wizard")
		fmt.Println()
		answers, err = wizard.RunInteractiveWizard()
		if err != nil {
			return fmt.Errorf("wizard cancelled or failed: %w", err)
		}
	}

	fmt.Println()
	logger.Info("Configuration collected", "repo", fmt.Sprintf("%s/%s", answers.RepoOwner, answers.RepoName))

	// Step 2: Create workspace
	creator := workspace.NewCreator(answers.WorkspacePath, dryRun, opLog)
	if err := creator.Create(); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	config := creator.GetConfig()

	// Step 3: Generate documentation
	generator := workspace.NewGenerator(config, answers, dryRun, opLog)
	if err := generator.GenerateAll(); err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	// Step 4: Initialize git
	repo := git.NewRepository(config, dryRun, opLog)
	if err := repo.Init(); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	if err := repo.AddAll(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	commitMsg := fmt.Sprintf("Initial commit from repo-bootstrapper\n\nGenerated with repo-bootstrapper for %s", answers.AgentChoice)
	if err := repo.Commit(commitMsg); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	config.GitInitialized = true

	// Step 5: Initialize Spec Kit
	specKit := speckit.NewIntegration(dryRun, opLog)
	if err := specKit.Init(answers.RepoName, answers.AgentChoice, answers.WorkspacePath); err != nil {
		logger.Error("Spec Kit initialization failed (continuing anyway)", err)
	} else {
		config.SpecKitReady = true
	}

	// Step 6: Create GitHub repository
	gh := git.NewGitHub(dryRun, opLog)
	if err := gh.CreateRepo(answers.RepoOwner, answers.RepoName, answers.Visibility, answers.WorkspacePath); err != nil {
		logger.Error("GitHub repo creation failed (continuing anyway)", err)
	} else {
		config.GitRemoteURL = fmt.Sprintf("https://github.com/%s/%s", answers.RepoOwner, answers.RepoName)
	}

	// Step 7: Create labels if requested
	if answers.IssueLabels != "none" {
		if err := gh.CreateLabels(answers.IssueLabels, answers.WorkspacePath); err != nil {
			logger.Error("Failed to create labels (continuing anyway)", err)
		}
	}

	// Step 8: Print summary
	fmt.Println()
	fmt.Println("=== Summary ===")
	opLog.PrintSummary(dryRun)
	fmt.Println()

	if dryRun {
		fmt.Println("Dry run complete. No changes were made.")
		fmt.Println()
		fmt.Println("To execute for real, run without --dry-run flag:")
		if fromFile != "" {
			fmt.Printf("  repo-bootstrapper new --from-file %s\n", fromFile)
		} else {
			fmt.Println("  repo-bootstrapper new")
		}
	} else {
		fmt.Println("✓ Workspace created:", answers.WorkspacePath)
		fmt.Println("✓ Git initialized")
		if config.SpecKitReady {
			fmt.Println("✓ Spec Kit initialized")
		}
		if config.GitRemoteURL != "" {
			fmt.Println("✓ GitHub repo created:", config.GitRemoteURL)
		}
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. cd %s\n", answers.WorkspacePath)
		fmt.Printf("  2. Open in %s\n", answers.AgentChoice)
		fmt.Println("  3. Read prompts/primer.md for AI workflow instructions")
		fmt.Println("  4. Run /implement to build your application")
	}
	fmt.Println()

	return nil
}
