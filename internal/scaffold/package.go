package scaffold

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
)

// PackageJSON represents the structure of a package.json file.
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Private         bool              `json:"private,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	DevDependencies map[string]string `json:"devDependencies,omitempty"`
}

// ReadPackageJSON reads and parses a package.json file from the given path.
func ReadPackageJSON(path string) (*PackageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	return &pkg, nil
}

// Write writes the package.json to the specified path.
func (p *PackageJSON) Write(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package.json: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	return nil
}

// MergeDependencies merges dependencies from another package.json into this one.
// The source dependencies will be added to the destination, overwriting any existing ones.
func (p *PackageJSON) MergeDependencies(source *PackageJSON) {
	if source.Dependencies != nil {
		if p.Dependencies == nil {
			p.Dependencies = make(map[string]string)
		}
		maps.Copy(p.Dependencies, source.Dependencies)
	}

	if source.DevDependencies != nil {
		if p.DevDependencies == nil {
			p.DevDependencies = make(map[string]string)
		}
		maps.Copy(p.DevDependencies, source.DevDependencies)
	}
}

// MergeScripts merges scripts from another package.json into this one.
func (p *PackageJSON) MergeScripts(source *PackageJSON) {
	if source.Scripts != nil {
		if p.Scripts == nil {
			p.Scripts = make(map[string]string)
		}
		maps.Copy(p.Scripts, source.Scripts)
	}
}

// LoadAndMerge reads a package.json file and merges its contents into this one.
func (p *PackageJSON) LoadAndMerge(path string) error {
	source, err := ReadPackageJSON(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to merge
		}
		return err
	}

	p.MergeDependencies(source)
	p.MergeScripts(source)
	return nil
}

// LoadScriptsAndMerge reads a scripts.json file and merges its scripts into this package.json.
func (p *PackageJSON) LoadScriptsAndMerge(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to merge
		}
		return fmt.Errorf("failed to read scripts file: %w", err)
	}

	var scriptsWrapper struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &scriptsWrapper); err != nil {
		return fmt.Errorf("failed to parse scripts file: %w", err)
	}

	if scriptsWrapper.Scripts != nil {
		if p.Scripts == nil {
			p.Scripts = make(map[string]string)
		}
		maps.Copy(p.Scripts, scriptsWrapper.Scripts)
	}

	return nil
}
