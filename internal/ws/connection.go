package ws

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbacatePay/abacatepay-cli/internal/clierr"
	"github.com/AbacatePay/abacatepay-cli/internal/style"

	"github.com/gorilla/websocket"
)

type Handler func(ctx context.Context, conn *websocket.Conn) error

// Status describes a transition in the connection lifecycle, reported
// through Config.OnStatus when set.
type Status int

const (
	StatusConnecting Status = iota
	StatusConnected
	StatusRetrying
)

type Config struct {
	URL        string
	Headers    http.Header
	MaxRetries int
	MinBackoff time.Duration
	MaxBackoff time.Duration

	// OnStatus, if set, is called on every connection lifecycle transition.
	// It must return quickly and must not block.
	OnStatus func(Status)
}

func (cfg Config) reportStatus(s Status) {
	if cfg.OnStatus != nil {
		cfg.OnStatus(s)
	}
}

func ConnectWithRetry(ctx context.Context, cfg Config, handler Handler) error {
	backoff := cfg.MinBackoff
	retries := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		slog.Debug("Connecting...", "url", cfg.URL)
		cfg.reportStatus(StatusConnecting)

		conn, _, err := websocket.DefaultDialer.DialContext(ctx, cfg.URL, cfg.Headers)
		if err != nil {
			retries++

			if cfg.MaxRetries > 0 && retries >= cfg.MaxRetries {
				errMsg := fmt.Sprintf("Failed to connect to %s after %d retries: %v", cfg.URL, retries, err)
				style.PrintError(errMsg)
				return clierr.MarkDisplayed(fmt.Errorf("%s", errMsg))
			}

			slog.Warn(
				"Connection failed, retrying…",
				"error", err,
				"backoff", backoff,
				"retry", fmt.Sprintf("%d/%d", retries, cfg.MaxRetries),
			)
			cfg.reportStatus(StatusRetrying)

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
		cfg.reportStatus(StatusConnected)
		backoff = cfg.MinBackoff
		retries = 0

		if err := handler(ctx, conn); err != nil {
			slog.Warn("Connection lost", "error", err)
		}

		conn.Close()
	}
}
