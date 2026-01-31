// Package scaffold provides project scaffolding functionality for AbacatePay CLI.
// It handles cloning templates, composing project layers, and setting up new projects.
package scaffold

import (
	"fmt"
	"strings"
)

// Config holds the scaffold configuration for creating a new project.
type Config struct {
	ProjectName string // Name of the project (directory name)
	Framework   string // Framework to use: "next" or "elysia"
	Linter      string // Linter to use: "eslint" or "biome"
	BetterAuth  bool   // Whether to include BetterAuth configuration
}

// Validate checks if the configuration is valid.
func (c Config) Validate() error {
	if err := ValidateProjectName(c.ProjectName); err != nil {
		return err
	}

	validFrameworks := map[string]bool{"next": true, "elysia": true}
	if !validFrameworks[c.Framework] {
		return fmt.Errorf("invalid framework: %s (must be 'next' or 'elysia')", c.Framework)
	}

	validLinters := map[string]bool{"eslint": true, "biome": true}
	if !validLinters[c.Linter] {
		return fmt.Errorf("invalid linter: %s (must be 'eslint' or 'biome')", c.Linter)
	}

	return nil
}

// ValidateProjectName checks if the project name is valid.
// It ensures the name is not empty and doesn't contain invalid characters.
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}

	// Check for invalid characters that are not allowed in directory names
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(name, char) {
			return fmt.Errorf("project name cannot contain '%s'", char)
		}
	}

	return nil
}

// GetNextSteps returns the list of commands the user should run after project creation.
func (c Config) GetNextSteps() []string {
	steps := []string{
		fmt.Sprintf("cd %s", c.ProjectName),
		"bun install",
		"cp .env.example .env",
	}

	// Both frameworks use the same commands for now
	steps = append(steps, "bun run db:push")
	steps = append(steps, "bun run dev")

	return steps
}
