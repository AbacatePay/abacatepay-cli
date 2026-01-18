package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
)

type TailListener struct {
	BaseListener
}

func NewTailListener(cfg *config.Config, token string) *TailListener {
	return &TailListener{
		BaseListener: BaseListener{
			Cfg:   cfg,
			Token: token,
		},
	}
}

func (t *TailListener) Listen(ctx context.Context) error {
	slog.Info("Starting tail listener...")

	return ws.ConnectWithRetry(ctx, t.WSConfig(), t.readLoop)
}

func (t *TailListener) readLoop(ctx context.Context, conn *websocket.Conn) error {
	t.SetupConn(conn)

	go t.Heartbeat(ctx, conn)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		t.SetReadDeadline(conn)

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

func (t *TailListener) displayWebhook(event, id string, rawBody []byte) {
	style.LogWebhookReceived(event, id)

	if !t.Cfg.Verbose {
		return
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, rawBody, "", "  "); err != nil {
		fmt.Println(string(rawBody))
		return
	}
	fmt.Println(buf.String())
}
