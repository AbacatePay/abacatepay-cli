package webhook

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/AbacatePay/abacatepay-cli/internal/config"
	"github.com/AbacatePay/abacatepay-cli/internal/ws"

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

	go func() {
		_ = t.Heartbeat(ctx, conn)
	}()

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

		event, id, err := parseEnvelope(message)
		if err != nil {
			t.emit(Event{Kind: EventInvalid, Time: time.Now()})
			continue
		}

		t.displayWebhook(event, id, message)
	}
}

func (t *TailListener) displayWebhook(event, id string, rawBody []byte) {
	t.emit(Event{
		Kind:    EventReceived,
		Time:    time.Now(),
		Name:    event,
		ID:      id,
		RawJSON: prettyJSON(t.Cfg.Verbose, rawBody),
	})
}
