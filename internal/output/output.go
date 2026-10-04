package output

import (
	"github.com/AbacatePay/abacatepay-cli/internal/clierr"
	"github.com/AbacatePay/abacatepay-cli/internal/style"
	"github.com/AbacatePay/abacatepay-cli/internal/tui"
)

type Result struct {
	Title  string
	Fields map[string]string
}

func Print(r Result) {
	style.PrintSuccess(r.Title, r.Fields)
}

// RunTask runs fn with an animated spinner (via tui.RunTask), and the
// spinner's final frame becomes the result box - one continuous render, not
// a spinner followed by a separately-printed box.
//
// RunTask itself displays a non-nil error before returning it, so the
// returned error is wrapped with clierr.MarkDisplayed - callers (ultimately
// cmd.Exec's top-level handler) must check clierr.AlreadyDisplayed before
// printing it again.
//
// fn must not print to stdout itself; see tui.RunTask.
func RunTask(message string, fn func() (Result, error)) error {
	err := tui.RunTask(message, func() (tui.TaskResult, error) {
		r, err := fn()
		return tui.TaskResult{Title: r.Title, Fields: r.Fields}, err
	})
	return clierr.MarkDisplayed(err)
}

func Error(msg string) {
	style.PrintError(msg)
}
