package tui

import (
	"fmt"

	"github.com/AbacatePay/abacatepay-cli/internal/style"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TaskResult is what a Task renders once it finishes successfully.
type TaskResult struct {
	Title  string
	Fields map[string]string
}

// RunTask runs fn with an animated spinner next to message, then renders its
// outcome - success or error text, using the exact same rendering
// RenderSuccess/RenderError draw everywhere else - as part of the same
// continuous render. There is no separate print step after the spinner: the
// spinner frame is simply replaced by the final text in place, so a
// command's entire visible output is one uninterrupted TUI render instead of
// a bubbletea animation followed by a disconnected plain fmt.Println.
//
// fn must not print to stdout itself - anything it writes while the spinner
// is animating races the renderer and corrupts the display. Callers that
// need to show verbose/debug output should skip RunTask entirely for that
// invocation and print directly instead.
//
// When stdout isn't a terminal, RunTask prints message once and renders the
// final result as a plain line, matching how every other styled output
// falls back for piped/scripted use.
func RunTask(message string, fn func() (TaskResult, error)) error {
	if !IsInteractive() {
		style.PrintInfo(message)

		res, err := fn()
		if err != nil {
			style.PrintError(err.Error())
			return err
		}

		style.PrintSuccess(res.Title, res.Fields)
		return nil
	}

	m := newTaskModel(message, fn)

	final, runErr := tea.NewProgram(m).Run()
	if runErr != nil {
		return runErr
	}

	return final.(taskModel).err
}

type taskDoneMsg struct {
	result TaskResult
	err    error
}

type taskModel struct {
	spin    spinner.Model
	message string
	fn      func() (TaskResult, error)
	done    bool
	result  TaskResult
	err     error
}

func newTaskModel(message string, fn func() (TaskResult, error)) taskModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(style.Palette.Green)

	return taskModel{spin: s, message: message, fn: fn}
}

func (m taskModel) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, runTaskFn(m.fn))
}

func runTaskFn(fn func() (TaskResult, error)) tea.Cmd {
	return func() tea.Msg {
		res, err := fn()
		return taskDoneMsg{result: res, err: err}
	}
}

func (m taskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case taskDoneMsg:
		m.done = true
		m.result = msg.result
		m.err = msg.err
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m taskModel) View() string {
	if m.done {
		if m.err != nil {
			return style.RenderError(m.err.Error())
		}
		return style.RenderSuccess(m.result.Title, m.result.Fields)
	}

	return fmt.Sprintf("%s %s\n", m.spin.View(), m.message)
}
