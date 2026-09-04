package workspace

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/repo-bootstrapper/repo-bootstrapper/internal/logger"
	"github.com/repo-bootstrapper/repo-bootstrapper/internal/wizard"
)

//go:embed templates/*
var templatesFS embed.FS

// TemplateVars holds variables for template rendering
type TemplateVars struct {
	ProjectName    string
	RepoOwner      string
	RepoName       string
	License        string
	Language       string
	Framework      string
	AgentChoice    string
	AppKind        string
	CurrentDate    string
	BranchStrategy string
	TargetOS       string
	PackageManager string
	TestFramework  string
}

// Generator handles document generation from templates
type Generator struct {
	config *WorkspaceConfig
	vars   *TemplateVars
	dryRun bool
	opLog  *logger.OperationLog
}

// NewGenerator creates a new document generator
func NewGenerator(config *WorkspaceConfig, answers *wizard.WizardAnswers, dryRun bool, opLog *logger.OperationLog) *Generator {
	return &Generator{
		config: config,
		vars:   answersToVars(answers),
		dryRun: dryRun,
		opLog:  opLog,
	}
}

func answersToVars(answers *wizard.WizardAnswers) *TemplateVars {
	targetOS := "cross-platform"
	if len(answers.TargetOS) > 0 {
		targetOS = answers.TargetOS[0]
		if len(answers.TargetOS) > 1 {
			targetOS = "Windows, macOS, Linux"
		}
	}

	return &TemplateVars{
		ProjectName:    answers.RepoName,
		RepoOwner:      answers.RepoOwner,
		RepoName:       answers.RepoName,
		License:        answers.License,
		Language:       answers.Language,
		Framework:      answers.Framework,
		AgentChoice:    answers.AgentChoice,
		AppKind:        answers.AppKind,
		CurrentDate:    time.Now().Format("2006-01-02"),
		BranchStrategy: "Git Flow: main, develop, feature/*, release/*, hotfix/*",
		TargetOS:       targetOS,
		PackageManager: answers.PackageManager,
		TestFramework:  answers.TestFramework,
	}
}

// GenerateAll generates all documentation files
func (g *Generator) GenerateAll() error {
	files := []struct {
		template string
		output   string
	}{
		{"README.tmpl", "README.md"},
		{"constitution.tmpl", "docs/constitution.md"},
		{"spec.tmpl", "docs/spec.md"},
		{"plan.tmpl", "docs/plan.md"},
		{"tasks.tmpl", "docs/tasks.md"},
		{"primer.tmpl", "prompts/primer.md"},
	}

	for _, f := range files {
		if err := g.generateFile(f.template, f.output); err != nil {
			return fmt.Errorf("failed to generate %s: %w", f.output, err)
		}
	}

	// Generate .gitignore
	if err := g.generateGitignore(); err != nil {
		return fmt.Errorf("failed to generate .gitignore: %w", err)
	}

	// Generate LICENSE if not "None"
	if g.vars.License != "None" {
		if err := g.generateLicense(); err != nil {
			return fmt.Errorf("failed to generate LICENSE: %w", err)
		}
	}

	return nil
}

func (g *Generator) generateFile(templateName, outputPath string) error {
	// Read template
	templatePath := "templates/" + templateName
	content, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	// Parse and execute template
	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.vars); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to file
	fullPath := filepath.Join(g.config.RootPath, outputPath)

	if g.dryRun {
		logger.Info("[DRY RUN] Would create file", "path", fullPath, "size", buf.Len())
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: fmt.Sprintf("Generate %s (%d bytes)", outputPath, buf.Len()),
			DryRun:  true,
		})
	} else {
		if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
			g.opLog.Add(logger.Operation{
				Type:    logger.OpCreateFile,
				Target:  fullPath,
				Details: fmt.Sprintf("Generate %s", outputPath),
				Error:   err,
			})
			return fmt.Errorf("failed to write file: %w", err)
		}
		logger.Info("Generated file", "path", outputPath)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: fmt.Sprintf("Generate %s (%d bytes)", outputPath, buf.Len()),
		})
		g.config.AddGeneratedDoc(fullPath)
		g.config.AddCreatedFile(fullPath)
	}

	return nil
}

func (g *Generator) generateGitignore() error {
	content := getGitignoreContent(g.vars.Language)
	fullPath := filepath.Join(g.config.RootPath, ".gitignore")

	if g.dryRun {
		logger.Info("[DRY RUN] Would create .gitignore", "path", fullPath)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: "Generate .gitignore",
			DryRun:  true,
		})
	} else {
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: "Generate .gitignore",
		})
		g.config.AddCreatedFile(fullPath)
	}

	return nil
}

func (g *Generator) generateLicense() error {
	content, err := getLicenseContent(g.vars.License, g.vars.RepoOwner)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(g.config.RootPath, "LICENSE")

	if g.dryRun {
		logger.Info("[DRY RUN] Would create LICENSE", "path", fullPath, "license", g.vars.License)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: fmt.Sprintf("Generate LICENSE (%s)", g.vars.License),
			DryRun:  true,
		})
	} else {
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
		logger.Info("Generated LICENSE", "license", g.vars.License)
		g.opLog.Add(logger.Operation{
			Type:    logger.OpCreateFile,
			Target:  fullPath,
			Details: fmt.Sprintf("Generate LICENSE (%s)", g.vars.License),
		})
		g.config.AddCreatedFile(fullPath)
	}

	return nil
}

func getGitignoreContent(language string) string {
	// Basic gitignore patterns
	base := `# OS
.DS_Store
Thumbs.db

# IDE
.vscode/
.idea/
*.swp

# Temporary
*.tmp
*.log
`

	// Add language-specific patterns
	switch language {
	case "Go":
		base += `
# Go
*.exe
*.test
*.out
vendor/
`
	case "Python":
		base += `
# Python
__pycache__/
*.py[cod]
*.so
.Python
venv/
`
	case "JavaScript", "TypeScript":
		base += `
# Node
node_modules/
dist/
*.log
`
	}

	return base
}

func getLicenseContent(license, owner string) (string, error) {
	year := time.Now().Year()

	switch license {
	case "MIT":
		return fmt.Sprintf(`MIT License

Copyright (c) %d %s

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`, year, owner), nil

	case "Apache-2.0":
		return fmt.Sprintf(`Apache License
Version 2.0, January 2004
http://www.apache.org/licenses/

Copyright %d %s

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`, year, owner), nil

	default:
		return fmt.Sprintf("Copyright (c) %d %s. All rights reserved.\n", year, owner), nil
	}
}
