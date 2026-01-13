package ws

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Handler func(ctx context.Context, conn *websocket.Conn) error

type Config struct {
	URL        string
	Headers    http.Header
	MaxRetries int
	MinBackoff time.Duration
	MaxBackoff time.Duration
}

func ConnectWithRetry(ctx context.Context, cfg Config, handler Handler) error {
	backoff := cfg.MinBackoff

	for {
		select {

		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		slog.Debug("Connecting...", "url", cfg.URL)

		conn, _, err := websocket.DefaultDialer.DialContext(ctx, cfg.URL, cfg.Headers)
		if err != nil {
			slog.Warn(
				"Connection failed, retrying…",
				"error", err,
				"backoff", backoff,
			)

			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-time.After(backoff):
				backoff *= 2
				if backoff > cfg.MaxBackoff {
					backoff = cfg.MaxBackoff
				}

				continue
			}
		}

		slog.Info("WebSocket connected")

		backoff = cfg.MinBackoff

		if err := handler(ctx, conn); err != nil {
			slog.Warn("Connection lost", "error", err)
		}

		conn.Close()
	}
}
