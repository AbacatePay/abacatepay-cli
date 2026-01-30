package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
)

const templatesRepoURL = "https://github.com/AbacatePay/templates.git"

// ProjectBuilder handles the construction of a new project by composing template layers.
type ProjectBuilder struct {
	config       Config
	templatesDir string
	projectPath  string
}

// NewProjectBuilder creates a new ProjectBuilder instance.
func NewProjectBuilder(cfg Config, targetDir string) (*ProjectBuilder, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &ProjectBuilder{
		config:      cfg,
		projectPath: filepath.Join(targetDir, cfg.ProjectName),
	}, nil
}

// Build executes the complete project scaffolding process.
// It clones templates, applies layers, and finalizes the project.
func (pb *ProjectBuilder) Build() error {
	// Create temporary directory for templates
	tempDir, err := os.MkdirTemp("", "abacate-templates-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone templates repository
	if err := pb.cloneTemplates(tempDir); err != nil {
		return err
	}

	pb.templatesDir = filepath.Join(tempDir, "templates")

	// Create project directory
	if err := EnsureDir(pb.projectPath, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Apply template layers
	if err := pb.applyBaseTemplate(); err != nil {
		return err
	}

	if err := pb.applyLinter(); err != nil {
		return err
	}

	if err := pb.applyFeatures(); err != nil {
		return err
	}

	if err := pb.finalizePackageJSON(); err != nil {
		return err
	}

	return nil
}

// cloneTemplates clones the templates repository into the temp directory.
func (pb *ProjectBuilder) cloneTemplates(tempDir string) error {
	if err := GitClone(tempDir); err != nil {
		return fmt.Errorf("failed to clone templates repository: %w", err)
	}
	return nil
}

// applyBaseTemplate copies the base template for the selected framework.
func (pb *ProjectBuilder) applyBaseTemplate() error {
	basePath := filepath.Join(pb.templatesDir, "base", pb.config.Framework)

	if !DirExists(basePath) {
		return fmt.Errorf("base template not found for framework: %s", pb.config.Framework)
	}

	if err := CopyDir(basePath, pb.projectPath); err != nil {
		return fmt.Errorf("failed to copy base template: %w", err)
	}

	return nil
}

// applyLinter applies the selected linter configuration to the project.
func (pb *ProjectBuilder) applyLinter() error {
	linterDir := filepath.Join(pb.templatesDir, "linters", pb.config.Linter)

	switch pb.config.Linter {
	case "eslint":
		if err := pb.applyESLint(linterDir); err != nil {
			return err
		}
	case "biome":
		if err := pb.applyBiome(linterDir); err != nil {
			return err
		}
	}

	return nil
}

// applyESLint copies the ESLint configuration file.
func (pb *ProjectBuilder) applyESLint(linterDir string) error {
	configFile := "eslint.config.elysia.js"
	if pb.config.Framework == "next" {
		configFile = "eslint.config.next.js"
	}

	src := filepath.Join(linterDir, configFile)
	dst := filepath.Join(pb.projectPath, "eslint.config.js")

	if err := CopyFile(src, dst); err != nil {
		return fmt.Errorf("failed to copy ESLint config: %w", err)
	}

	return nil
}

// applyBiome copies the Biome configuration file.
func (pb *ProjectBuilder) applyBiome(linterDir string) error {
	src := filepath.Join(linterDir, "biome.json")
	dst := filepath.Join(pb.projectPath, "biome.json")

	if err := CopyFile(src, dst); err != nil {
		return fmt.Errorf("failed to copy Biome config: %w", err)
	}

	return nil
}

// applyFeatures applies selected features (like BetterAuth) to the project.
func (pb *ProjectBuilder) applyFeatures() error {
	if pb.config.BetterAuth {
		if err := pb.applyBetterAuth(); err != nil {
			return err
		}
	}

	return nil
}

// applyBetterAuth applies the BetterAuth feature to the project.
func (pb *ProjectBuilder) applyBetterAuth() error {
	authDir := filepath.Join(pb.templatesDir, "features", "betterauth", pb.config.Framework)

	if !DirExists(authDir) {
		return fmt.Errorf("betterauth template not found for framework: %s", pb.config.Framework)
	}

	if err := CopyDir(authDir, pb.projectPath); err != nil {
		return fmt.Errorf("failed to copy better-auth files: %w", err)
	}

	return nil
}

// finalizePackageJSON merges all package.json fragments and updates the project name.
func (pb *ProjectBuilder) finalizePackageJSON() error {
	mainPackagePath := filepath.Join(pb.projectPath, "package.json")

	// Read main package.json
	pkg, err := ReadPackageJSON(mainPackagePath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	// Merge linter dependencies and scripts
	linterDir := filepath.Join(pb.templatesDir, "linters", pb.config.Linter)
	if err := pkg.LoadAndMerge(filepath.Join(linterDir, "package.json")); err != nil {
		return fmt.Errorf("failed to merge linter dependencies: %w", err)
	}
	if err := pkg.LoadScriptsAndMerge(filepath.Join(linterDir, "scripts.json")); err != nil {
		return fmt.Errorf("failed to merge linter scripts: %w", err)
	}

	// Merge better-auth dependencies if enabled
	if pb.config.BetterAuth {
		authDir := filepath.Join(pb.templatesDir, "features", "betterauth", pb.config.Framework)
		if err := pkg.LoadAndMerge(filepath.Join(authDir, "package.json")); err != nil {
			return fmt.Errorf("failed to merge better-auth dependencies: %w", err)
		}
	}

	// Update project name
	pkg.Name = pb.config.ProjectName

	// Write final package.json
	if err := pkg.Write(mainPackagePath); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	return nil
}

// ScaffoldProject is a convenience function that creates and builds a project in one call.
// It's the main entry point for project scaffolding.
func ScaffoldProject(cfg Config, targetDir string) error {
	builder, err := NewProjectBuilder(cfg, targetDir)
	if err != nil {
		return err
	}

	return builder.Build()
}
