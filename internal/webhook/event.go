package webhook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/AbacatePay/abacatepay-cli/internal/style"
)

// EventKind identifies what a listener Event represents.
type EventKind int

const (
	// EventReceived is emitted for every webhook message read off the socket.
	EventReceived EventKind = iota
	// EventForwarded is emitted after an attempt to forward a webhook to the
	// configured local URL.
	EventForwarded
	// EventInvalid is emitted when a socket message couldn't be parsed as JSON.
	EventInvalid
)

// Event describes a single occurrence in a listener's lifecycle: a webhook
// received from the relay, or the result of forwarding one to the user's
// local server. Listeners report these through BaseListener.Emit, which
// defaults to printing plain lines (today's behavior) when unset, and can be
// overridden by a caller (e.g. the interactive `listen` dashboard) to route
// events elsewhere instead.
type Event struct {
	Kind       EventKind
	Time       time.Time
	Name       string // webhook event type, e.g. "payment.completed"
	ID         string
	StatusCode int
	StatusText string
	RawJSON    string // syntax-highlighted body (via style.RenderJSON), only populated in verbose mode
}

// emit routes e through BaseListener.Emit if one is set, or falls back to
// defaultEmit to preserve the plain-line output the CLI has always printed.
func (b *BaseListener) emit(e Event) {
	if b.Emit != nil {
		b.Emit(e)
		return
	}
	defaultEmit(e)
}

func defaultEmit(e Event) {
	switch e.Kind {
	case EventReceived:
		style.LogWebhookReceived(e.Name, e.ID)
		if e.RawJSON != "" {
			fmt.Println(e.RawJSON)
		}
	case EventForwarded:
		style.LogWebhookForwarded(e.StatusCode, e.StatusText, e.Name)
	case EventInvalid:
		style.PrintError("Received invalid JSON from WebSocket")
	}
}

// prettyJSON returns a syntax-highlighted rendering of rawBody (matching
// every other JSON the CLI prints, via style.RenderJSON), or "" when verbose
// is false (the caller should skip showing the body at all in that case).
func prettyJSON(verbose bool, rawBody []byte) string {
	if !verbose {
		return ""
	}

	return style.RenderJSON(json.RawMessage(rawBody))
}
