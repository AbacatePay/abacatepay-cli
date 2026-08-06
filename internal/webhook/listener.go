package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbacatePay/abacatepay-cli/internal/crypto"
	"github.com/AbacatePay/abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

func (l *Listener) Listen(ctx context.Context) error {
	slog.Info("Starting webhook listener...")

	return ws.ConnectWithRetry(ctx, l.WSConfig(), l.readLoop)
}

func (l *Listener) readLoop(ctx context.Context, conn *websocket.Conn) error {
	const requestLimitPerSecond int = 10

	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(requestLimitPerSecond)

	l.SetupConn(conn)

	g.Go(func() error {
		return l.Heartbeat(gCtx, conn)
	})

	for {
		select {
		case <-gCtx.Done():
			return g.Wait()

		default:
		}

		l.SetReadDeadline(conn)

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
			l.emit(Event{Kind: EventInvalid, Time: time.Now()})
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

func (l *Listener) displayWebhook(meta webhookMetadata, rawBody []byte) {
	l.txLogger.Info("webhook_received",
		"event", meta.Event,
		"id", meta.ID,
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(rawBody),
		"raw_message", string(rawBody),
	)

	l.emit(Event{
		Kind:    EventReceived,
		Time:    time.Now(),
		Name:    meta.Event,
		ID:      meta.ID,
		RawJSON: prettyJSON(l.Cfg.Verbose, rawBody),
	})
}

func (l *Listener) forward(ctx context.Context, message []byte, event string) error {
	startTime := time.Now()

	signature := crypto.SignWebhookPayload(message)

	resp, err := l.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Webhook-Signature", signature).
		SetBody(message).
		Post(l.forwardURL)

	duration := time.Since(startTime)

	if err != nil {
		l.txLogger.Error("webhook_forward_failed",
			"event", event,
			"url", l.forwardURL,
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("failed to forward webhook: %w", err)
	}

	statusCode := resp.StatusCode()
	l.emit(Event{
		Kind:       EventForwarded,
		Time:       time.Now(),
		Name:       event,
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
	})

	if statusCode < 200 || statusCode >= 300 {
		l.txLogger.Error("webhook_forward_error",
			"event", event,
			"url", l.forwardURL,
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
			"response_body", string(resp.Body()),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return nil
	}

	l.txLogger.Info("webhook_forwarded",
		"event", event,
		"url", l.forwardURL,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(message),
	)

	return nil
}
