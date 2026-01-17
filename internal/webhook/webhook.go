package webhook

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/style"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type webhookMetadata struct {
	Event string
	ID    string
}

type Listener struct {
	cfg           *config.Config
	client        *resty.Client
	forwardURL    string
	token         string
	txLogger      *slog.Logger
	connMu        sync.Mutex
	signingSecret string
}

func NewListener(cfg *config.Config, client *resty.Client, forwardURL, token string, txLogger *slog.Logger) *Listener {
	return &Listener{
		cfg:           cfg,
		client:        client,
		forwardURL:    forwardURL,
		token:         token,
		txLogger:      txLogger,
		signingSecret: "whsec_mock_" + hex.EncodeToString([]byte(time.Now().Format("150405"))),
	}
}

func (l *Listener) SetupConn(conn *websocket.Conn) {
	conn.SetPongHandler(func(string) error {
		l.connMu.Lock()
		defer l.connMu.Unlock()
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})
}

func (l *Listener) readLoop(ctx context.Context, conn *websocket.Conn) error {
	const requestLimitPerSecond int = 10

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(requestLimitPerSecond)

	l.SetupConn(conn)

	g.Go(func() error {
		return l.heartbeat(gCtx, conn)
	})

	for {
		select {
		case <-gCtx.Done():
			return g.Wait()

		default:
		}

		l.connMu.Lock()
		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		l.connMu.Unlock()

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Info("WebSocket connection closed")
				_ = g.Wait()
				return nil
			}

			if gCtx.Err() != nil {
				_ = g.Wait()
				return nil
			}

			_ = g.Wait()
			return fmt.Errorf("failed to read websocket message: %w", err)
		}

		var raw struct {
			Event string `json:"event"`
			Data  struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		if err := json.Unmarshal(message, &raw); err != nil {
			style.PrintError("Received invalid JSON from WebSocket")
			continue
		}

		meta := webhookMetadata{Event: raw.Event, ID: raw.Data.ID}
		l.displayWebhook(meta, message)

		g.Go(func() error {
			_ = l.forward(gCtx, message, meta.Event)
			return nil
		})
	}
}

func (l *Listener) heartbeat(ctx context.Context, conn *websocket.Conn) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.connMu.Lock()
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second),
			)
			l.connMu.Unlock()
			return nil

		case <-ticker.C:
			l.connMu.Lock()
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second))
			l.connMu.Unlock()

			if err != nil {
				slog.Debug("Ping failed", "error", err)
				return err
			}
		}
	}
}
