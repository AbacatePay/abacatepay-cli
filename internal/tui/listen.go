package tui

import (
	"fmt"
	"strings"

	"github.com/AbacatePay/abacatepay-cli/internal/style"
	"github.com/AbacatePay/abacatepay-cli/internal/webhook"
	"github.com/AbacatePay/abacatepay-cli/internal/ws"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConnStatusMsg reports a change in connection status. Send it into a
// running Program (via Program.Send) to update the dashboard header.
type ConnStatusMsg struct {
	Status ws.Status
}

// StoppedMsg signals that the listener's Listen call returned - the socket
// loop exited, for any reason (including the user quitting). Send it into
// the Program when that happens so the header reflects it if still open.
type StoppedMsg struct{ Err error }

const maxEventLines = 500

type listenKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Bottom key.Binding
	Quit   key.Binding
}

func (k listenKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Bottom, k.Quit}
}

func (k listenKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var listenKeys = listenKeyMap{
	Up:     key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Bottom: key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "latest")),
	Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"), key.WithHelp("q", "quit")),
}

// ListenModel is the bubbletea dashboard shown by the interactive `listen`
// command: a live, scrollable feed of received/forwarded webhooks.
type ListenModel struct {
	forwardURL string
	cancel     func()

	spin spinner.Model
	vp   viewport.Model
	help help.Model

	status  ws.Status
	stopped bool
	stopErr error

	lines []string
	ready bool
}

// NewListenModel builds the dashboard model. cancel is invoked exactly once,
// when the user quits, so the caller can tear down the context driving the
// underlying websocket loop.
func NewListenModel(forwardURL string, cancel func()) *ListenModel {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(style.Palette.Green)

	return &ListenModel{
		forwardURL: forwardURL,
		cancel:     cancel,
		spin:       s,
		help:       help.New(),
		status:     ws.StatusConnecting,
	}
}

func (m *ListenModel) Init() tea.Cmd {
	return m.spin.Tick
}

func (m *ListenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, listenKeys.Quit) {
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		}
		if key.Matches(msg, listenKeys.Bottom) {
			m.vp.GotoBottom()
			return m, nil
		}

	case webhook.Event:
		m.appendEvent(msg)
		return m, nil

	case ConnStatusMsg:
		m.status = msg.Status
		return m, nil

	case StoppedMsg:
		m.stopped = true
		m.stopErr = msg.Err
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}

	if !m.ready {
		return m, nil
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *ListenModel) resize(width, height int) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
	vpHeight := height - headerHeight - footerHeight
	if vpHeight < 0 {
		vpHeight = 0
	}

	if !m.ready {
		m.vp = viewport.New(width, vpHeight)
		m.vp.SetContent(strings.Join(m.lines, "\n"))
		m.ready = true
	} else {
		m.vp.Width = width
		m.vp.Height = vpHeight
	}

	m.help.Width = width
}

func (m *ListenModel) appendEvent(e webhook.Event) {
	var line string
	switch e.Kind {
	case webhook.EventReceived:
		line = style.RenderWebhookReceived(e.Name, e.ID)
		if e.RawJSON != "" {
			line += "\n" + e.RawJSON
		}
	case webhook.EventForwarded:
		line = style.RenderWebhookForwarded(e.StatusCode, e.StatusText, e.Name)
	case webhook.EventInvalid:
		line = lipgloss.NewStyle().Foreground(style.Palette.SoftRed).Render("received invalid JSON from WebSocket")
	default:
		return
	}

	m.lines = append(m.lines, line)
	if len(m.lines) > maxEventLines {
		m.lines = m.lines[len(m.lines)-maxEventLines:]
	}

	if m.ready {
		m.vp.SetContent(strings.Join(m.lines, "\n"))
	}
}

func (m *ListenModel) headerView() string {
	var status string
	switch {
	case m.stopped:
		mark := lipgloss.NewStyle().Foreground(style.Palette.SoftRed).Bold(true).Render("✗ stopped")
		status = mark
		if m.stopErr != nil {
			status += lipgloss.NewStyle().Foreground(style.Palette.Gray).Render(" (" + m.stopErr.Error() + ")")
		}
	case m.status == ws.StatusConnected:
		status = lipgloss.NewStyle().Foreground(style.Palette.Green).Bold(true).Render("● connected")
	case m.status == ws.StatusRetrying:
		status = m.spin.View() + " reconnecting..."
	default:
		status = m.spin.View() + " connecting..."
	}

	title := style.TitleStyle.Render("🥑 abacatepay listen")
	forward := style.LabelStyle.Render("forwarding to ") + style.ValueStyle.Render(m.forwardURL)

	return lipgloss.JoinVertical(lipgloss.Left, title, status+"  "+forward, "")
}

func (m *ListenModel) footerView() string {
	return "\n" + m.help.View(listenKeys)
}

func (m *ListenModel) View() string {
	if !m.ready {
		return fmt.Sprintf("%s\nInitializing…", m.headerView())
	}

	return lipgloss.JoinVertical(lipgloss.Left, m.headerView(), m.vp.View(), m.footerView())
}
