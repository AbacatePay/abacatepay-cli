package scaffold

import (
	"fmt"
	"os"
	"os/exec"
)

// GitRepositoryURL is the URL of the templates repository.
const GitRepositoryURL = "https://github.com/albuquerquesz/abacatepay-templates.git"

// GitClone clones the templates repository into the specified directory.
// It performs a shallow clone (--depth 1) to save bandwidth and time.
func GitClone(destDir string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", GitRepositoryURL, destDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	return nil
}
