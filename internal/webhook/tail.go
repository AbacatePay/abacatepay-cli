package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
)

type TailListener struct {
	cfg     *config.Config
	token   string
	connMu  sync.Mutex
}

func NewTailListener(cfg *config.Config, token string) *TailListener {
	return &TailListener{
		cfg:   cfg,
		token: token,
	}
}

func (t *TailListener) Listen(ctx context.Context) error {
	slog.Info("Starting tail listener...")

	header := http.Header{}
	header.Add("Authorization", "Bearer "+t.token)

	cfg := ws.Config{
		URL:        t.cfg.WebSocketBaseURL,
		Headers:    header,
		MinBackoff: 1 * time.Second,
		MaxBackoff: 15 * time.Second,
		MaxRetries: 5,
	}

	return ws.ConnectWithRetry(ctx, cfg, t.readLoop)
}

func (t *TailListener) readLoop(ctx context.Context, conn *websocket.Conn) error {
	t.setupConn(conn)

	go t.heartbeat(ctx, conn)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		t.connMu.Lock()
		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		t.connMu.Unlock()

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Info("WebSocket connection closed")
				return nil
			}

			if ctx.Err() != nil {
				return nil
			}

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

		t.displayWebhook(raw.Event, raw.Data.ID, message)
	}
}

func (t *TailListener) setupConn(conn *websocket.Conn) {
	conn.SetPongHandler(func(string) error {
		t.connMu.Lock()
		defer t.connMu.Unlock()
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})
}

func (t *TailListener) heartbeat(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.connMu.Lock()
			conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(time.Second),
			)
			t.connMu.Unlock()
			return

		case <-ticker.C:
			t.connMu.Lock()
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second))
			t.connMu.Unlock()

			if err != nil {
				slog.Debug("Ping failed", "error", err)
				return
			}
		}
	}
}

func (t *TailListener) displayWebhook(event, id string, rawBody []byte) {
	style.LogWebhookReceived(event, id)

	if !t.cfg.Verbose {
		return
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, rawBody, "", "  "); err != nil {
		fmt.Println(string(rawBody))
		return
	}
	fmt.Println(buf.String())
}
