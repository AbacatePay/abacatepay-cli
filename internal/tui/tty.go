package tui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsInteractive reports whether stdout is attached to a terminal. Dynamic
// (bubbletea-driven) output should only be used when this is true; otherwise
// callers should fall back to plain sequential output so piped/scripted use
// keeps working.
func IsInteractive() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}
