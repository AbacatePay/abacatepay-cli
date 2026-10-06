package webhook

import (
	"context"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/AbacatePay/abacatepay-cli/internal/config"
	"github.com/AbacatePay/abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
)

type BaseListener struct {
	Cfg    *config.Config
	Token  string
	ConnMu sync.Mutex

	// Emit, if set, receives every Event instead of the default plain-line
	// output. Left nil, listeners behave exactly as before.
	Emit func(Event)

	// OnStatus, if set, receives connection lifecycle transitions (dialing,
	// connected, retrying). Left nil, no-op.
	OnStatus func(ws.Status)
}

func (b *BaseListener) SetupConn(conn *websocket.Conn) {
	conn.SetPongHandler(func(string) error {
		b.ConnMu.Lock()
		defer b.ConnMu.Unlock()
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})
}

func (b *BaseListener) Heartbeat(ctx context.Context, conn *websocket.Conn) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			b.ConnMu.Lock()
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second),
			)
			b.ConnMu.Unlock()
			return nil

		case <-ticker.C:
			b.ConnMu.Lock()
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second))
			b.ConnMu.Unlock()

			if err != nil {
				slog.Debug("Ping failed", "error", err)
				return err
			}
		}
	}
}

func (b *BaseListener) WSConfig() ws.Config {
	// services/ws-relay authenticates the connection via a ?session= query
	// param on the upgrade request, not an Authorization header.
	dialURL := b.Cfg.WebSocketBaseURL

	if u, err := url.Parse(b.Cfg.WebSocketBaseURL); err == nil {
		q := u.Query()
		q.Set("session", b.Token)
		u.RawQuery = q.Encode()
		dialURL = u.String()
	} else {
		slog.Debug("Failed to parse websocket base URL", "error", err)
	}

	return ws.Config{
		URL:        dialURL,
		MinBackoff: 1 * time.Second,
		MaxBackoff: 15 * time.Second,
		MaxRetries: 5,
		OnStatus:   b.OnStatus,
	}
}

func (b *BaseListener) SetReadDeadline(conn *websocket.Conn) {
	b.ConnMu.Lock()
	_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	b.ConnMu.Unlock()
}
