package scaffold

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config holds the scaffold configuration
type Config struct {
	ProjectName string
	Framework   string
	Linter      string
	BetterAuth  bool
}

// ScaffoldProject creates a new project by composing templates
func ScaffoldProject(cfg Config, targetDir string) error {
	// Create temporary directory for cloning
	tempDir, err := os.MkdirTemp("", "abacate-templates-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone templates repository
	if err := cloneTemplates(tempDir); err != nil {
		return fmt.Errorf("failed to clone templates: %w", err)
	}

	templatesDir := filepath.Join(tempDir, "templates")

	// Create project directory
	projectPath := filepath.Join(targetDir, cfg.ProjectName)
	if err := os.MkdirAll(projectPath, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// 1. Copy base template
	basePath := filepath.Join(templatesDir, "base", cfg.Framework)
	if err := copyDir(basePath, projectPath); err != nil {
		return fmt.Errorf("failed to copy base template: %w", err)
	}

	// 2. Apply linter
	if err := applyLinter(templatesDir, cfg.Framework, cfg.Linter, projectPath); err != nil {
		return fmt.Errorf("failed to apply linter: %w", err)
	}

	// 3. Apply BetterAuth if enabled
	if cfg.BetterAuth {
		if err := applyBetterAuth(templatesDir, cfg.Framework, projectPath); err != nil {
			return fmt.Errorf("failed to apply better-auth: %w", err)
		}
	}

	// 4. Merge package.json files
	if err := mergePackageJSONs(projectPath, templatesDir, cfg.Framework, cfg.Linter, cfg.BetterAuth); err != nil {
		return fmt.Errorf("failed to merge package.json: %w", err)
	}

	// 5. Update project name in package.json
	if err := updateProjectName(projectPath, cfg.ProjectName); err != nil {
		return fmt.Errorf("failed to update project name: %w", err)
	}

	return nil
}

// cloneTemplates clones the templates repository
func cloneTemplates(destDir string) error {
	repoURL := "https://github.com/AbacatePay/templates.git"

	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// applyLinter applies the linter configuration to the project
func applyLinter(templatesDir, framework, linter, projectPath string) error {
	linterDir := filepath.Join(templatesDir, "linters", linter)

	// Copy config file
	var configFile string
	if linter == "eslint" {
		if framework == "next" {
			configFile = "eslint.config.next.js"
		} else {
			configFile = "eslint.config.elysia.js"
		}
		if err := copyFile(
			filepath.Join(linterDir, configFile),
			filepath.Join(projectPath, "eslint.config.js"),
		); err != nil {
			return err
		}
	} else if linter == "biome" {
		if err := copyFile(
			filepath.Join(linterDir, "biome.json"),
			filepath.Join(projectPath, "biome.json"),
		); err != nil {
			return err
		}
	}

	return nil
}

// applyBetterAuth applies better-auth feature to the project
func applyBetterAuth(templatesDir, framework, projectPath string) error {
	authDir := filepath.Join(templatesDir, "features", "betterauth", framework)

	// Copy auth files
	if err := copyDir(authDir, projectPath); err != nil {
		return err
	}

	return nil
}

// mergePackageJSONs merges all package.json fragments into the main package.json
func mergePackageJSONs(projectPath, templatesDir, framework, linter string, betterAuth bool) error {
	mainPackagePath := filepath.Join(projectPath, "package.json")

	// Read main package.json
	mainPkg, err := readPackageJSON(mainPackagePath)
	if err != nil {
		return err
	}

	// Merge linter package.json
	linterPackagePath := filepath.Join(templatesDir, "linters", linter, "package.json")
	if err := mergePackageJSON(mainPkg, linterPackagePath); err != nil {
		return err
	}

	// Merge linter scripts
	linterScriptsPath := filepath.Join(templatesDir, "linters", linter, "scripts.json")
	if err := mergeScripts(mainPkg, linterScriptsPath); err != nil {
		return err
	}

	// Merge better-auth if enabled
	if betterAuth {
		authPackagePath := filepath.Join(templatesDir, "features", "betterauth", framework, "package.json")
		if err := mergePackageJSON(mainPkg, authPackagePath); err != nil {
			return err
		}
	}

	// Write merged package.json
	if err := writePackageJSON(mainPackagePath, mainPkg); err != nil {
		return err
	}

	return nil
}

// readPackageJSON reads and parses a package.json file
func readPackageJSON(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}

// mergePackageJSON merges a package.json file into the main package
func mergePackageJSON(main map[string]interface{}, mergePath string) error {
	data, err := os.ReadFile(mergePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var merge map[string]interface{}
	if err := json.Unmarshal(data, &merge); err != nil {
		return err
	}

	// Merge dependencies
	if deps, ok := merge["dependencies"].(map[string]interface{}); ok {
		if main["dependencies"] == nil {
			main["dependencies"] = make(map[string]interface{})
		}
		mainDeps := main["dependencies"].(map[string]interface{})
		for k, v := range deps {
			mainDeps[k] = v
		}
	}

	// Merge devDependencies
	if devDeps, ok := merge["devDependencies"].(map[string]interface{}); ok {
		if main["devDependencies"] == nil {
			main["devDependencies"] = make(map[string]interface{})
		}
		mainDevDeps := main["devDependencies"].(map[string]interface{})
		for k, v := range devDeps {
			mainDevDeps[k] = v
		}
	}

	return nil
}

// mergeScripts merges scripts from a scripts.json file
func mergeScripts(main map[string]interface{}, scriptsPath string) error {
	data, err := os.ReadFile(scriptsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var scripts map[string]any
	if err := json.Unmarshal(data, &scripts); err != nil {
		return err
	}

	if scriptsMap, ok := scripts["scripts"].(map[string]any); ok {
		if main["scripts"] == nil {
			main["scripts"] = make(map[string]any)
		}
		mainScripts := main["scripts"].(map[string]any)
		maps.Copy(mainScripts, scriptsMap)
	}

	return nil
}

func updateProjectName(projectPath, name string) error {
	pkg, err := readPackageJSON(filepath.Join(projectPath, "package.json"))
	if err != nil {
		return err
	}

	pkg["name"] = name

	return writePackageJSON(filepath.Join(projectPath, "package.json"), pkg)
}

func writePackageJSON(path string, pkg map[string]any) error {
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func GetNextSteps(cfg Config) []string {
	steps := []string{
		fmt.Sprintf("cd %s", cfg.ProjectName),
		"bun install",
		"cp .env.example .env",
	}

	if cfg.Framework == "next" {
		steps = append(steps, "bun run db:push")
		steps = append(steps, "bun run dev")
	} else {
		steps = append(steps, "bun run db:push")
		steps = append(steps, "bun run dev")
	}

	return steps
}

func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("project name cannot contain '%s'", char)
		}
	}

	return nil
}
