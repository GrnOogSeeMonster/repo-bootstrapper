package wizard

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	repoNameRegex  = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	repoOwnerRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
)

// ValidateRepoName validates a repository name
func ValidateRepoName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("repository name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("repository name cannot exceed 100 characters")
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("repository name cannot start or end with a hyphen")
	}
	if !repoNameRegex.MatchString(name) {
		return fmt.Errorf("repository name can only contain alphanumeric characters, dots, underscores, and hyphens")
	}
	return nil
}

// ValidateRepoOwner validates a repository owner/org name
func ValidateRepoOwner(owner string) error {
	if len(owner) == 0 {
		return fmt.Errorf("repository owner cannot be empty")
	}
	if len(owner) > 39 {
		return fmt.Errorf("repository owner cannot exceed 39 characters (GitHub limit)")
	}
	if !repoOwnerRegex.MatchString(owner) {
		return fmt.Errorf("repository owner can only contain alphanumeric characters and hyphens")
	}
	return nil
}

// ValidateLicense validates a license choice
func ValidateLicense(license string) error {
	validLicenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause", "MPL-2.0", "ISC", "Unlicense", "None"}
	for _, valid := range validLicenses {
		if strings.EqualFold(license, valid) {
			return nil
		}
	}
	return fmt.Errorf("invalid license: must be one of %v", validLicenses)
}

// ValidateAgentChoice validates an AI agent choice
func ValidateAgentChoice(agent string) error {
	validAgents := []string{"cursor", "claude", "copilot", "gemini"}
	for _, valid := range validAgents {
		if strings.EqualFold(agent, valid) {
			return nil
		}
	}
	return fmt.Errorf("invalid agent: must be one of %v", validAgents)
}

// ValidateAppKind validates an application kind
func ValidateAppKind(kind string) error {
	validKinds := []string{"CLI", "web", "service"}
	for _, valid := range validKinds {
		if strings.EqualFold(kind, valid) {
			return nil
		}
	}
	return fmt.Errorf("invalid app kind: must be one of %v", validKinds)
}

// ValidatePath validates and cleans a file path
func ValidatePath(path string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path to prevent directory traversal
	cleaned := filepath.Clean(path)

	// Expand home directory if present
	if strings.HasPrefix(cleaned, "~") {
		home, err := getHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine home directory: %w", err)
		}
		cleaned = filepath.Join(home, cleaned[1:])
	}

	return cleaned, nil
}

// getHomeDir returns the user's home directory
func getHomeDir() (string, error) {
	home, err := filepath.Abs("~")
	if err == nil && home != "~" {
		return home, nil
	}

	// Fallback to environment variables
	home = getEnv("HOME")
	if home == "" {
		home = getEnv("USERPROFILE") // Windows
	}
	if home == "" {
		return "", fmt.Errorf("could not determine home directory")
	}

	return home, nil
}

func getEnv(key string) string {
	// This will be replaced with actual os.Getenv in integration
	return ""
}
