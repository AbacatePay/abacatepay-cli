package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"abacatepay-cli/internal/crypto"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/ws"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

func (l *Listener) Listen(ctx context.Context, mock bool) error {
	if mock {
		style.LogSigningSecret(l.signingSecret)
		return l.mockListen(ctx)
	}

	slog.Info("Starting webhook listener...")

	header := http.Header{}
	header.Add("Authorization", "Bearer "+l.token)

	cfg := ws.Config{
		URL:        l.cfg.WebSocketBaseURL,
		Headers:    header,
		MinBackoff: 1 * time.Second,
		MaxBackoff: 15 * time.Second,
		MaxRetries: 5,
	}

	return ws.ConnectWithRetry(ctx, cfg, l.readLoop)
}

func (l *Listener) mockListen(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			id := fmt.Sprintf("pix_char_%d", time.Now().Unix())
			event := "billing.paid"

			mockData := map[string]any{
				"event": event,
				"data": map[string]any{
					"id":         id,
					"externalId": "order_123",
					"amount":     1000,
					"status":     "PAID",
				},
			}

			message, _ := json.Marshal(mockData)
			l.displayWebhook(webhookMetadata{Event: event, ID: id}, message)

			go func() {
				_ = l.forward(ctx, message, event)
			}()
		}
	}
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

func (l *Listener) displayWebhook(meta webhookMetadata, rawBody []byte) {
	style.LogWebhookReceived(meta.Event, meta.ID)

	l.txLogger.Info("webhook_received",
		"event", meta.Event,
		"id", meta.ID,
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(rawBody),
		"raw_message", string(rawBody),
	)

	if !l.cfg.Verbose {
		return
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, rawBody, "", "  "); err != nil {
		fmt.Println(string(rawBody))
		return
	}
	fmt.Println(buf.String())
}

func (l *Listener) forward(ctx context.Context, message []byte, event string) error {
	startTime := time.Now()
	timestamp := time.Now().Unix()

	signature := crypto.SignWebhookPayload(l.signingSecret, timestamp, message)

	resp, err := l.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Abacate-Signature", fmt.Sprintf("t=%d,v1=%s", timestamp, signature)).
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
	style.LogWebhookForwarded(statusCode, http.StatusText(statusCode), event)

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
