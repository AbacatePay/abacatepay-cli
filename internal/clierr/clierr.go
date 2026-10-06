// Package clierr marks errors that have already been shown to the user, so
// the top-level command runner doesn't print them a second time. It has no
// internal dependencies so every layer (ws, payments, output, tui, cmd) can
// import it without creating an import cycle.
package clierr

import "errors"

type displayed struct{ err error }

func (e *displayed) Error() string { return e.err.Error() }
func (e *displayed) Unwrap() error { return e.err }

// MarkDisplayed wraps err to record that it has already been shown to the
// user (a styled box, a TUI render, ...). Returns nil unchanged.
func MarkDisplayed(err error) error {
	if err == nil {
		return nil
	}
	return &displayed{err: err}
}

// AlreadyDisplayed reports whether err, or anything it wraps, was marked via
// MarkDisplayed.
func AlreadyDisplayed(err error) bool {
	var d *displayed
	return errors.As(err, &d)
}
