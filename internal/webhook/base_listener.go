package webhook

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
)

type BaseListener struct {
	Cfg    *config.Config
	Token  string
	ConnMu sync.Mutex
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
			conn.WriteControl(
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
	header := http.Header{}
	header.Add("Authorization", "Bearer "+b.Token)

	return ws.Config{
		URL:        b.Cfg.WebSocketBaseURL,
		Headers:    header,
		MinBackoff: 1 * time.Second,
		MaxBackoff: 15 * time.Second,
		MaxRetries: 5,
	}
}

func (b *BaseListener) SetReadDeadline(conn *websocket.Conn) {
	b.ConnMu.Lock()
	conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	b.ConnMu.Unlock()
}
