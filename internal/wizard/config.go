package wizard

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFromFile loads wizard answers from a JSON configuration file
func LoadFromFile(path string) (*WizardAnswers, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var answers WizardAnswers
	if err := json.Unmarshal(data, &answers); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Validate the loaded configuration
	if err := answers.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &answers, nil
}

// SaveToFile saves wizard answers to a JSON file (for --save-config option)
func SaveToFile(answers *WizardAnswers, path string) error {
	data, err := json.MarshalIndent(answers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
