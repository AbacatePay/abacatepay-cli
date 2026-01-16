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

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

type Message struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type Listener struct {
	cfg        *config.Config
	client     *resty.Client
	forwardURL string
	token      string
	txLogger   *slog.Logger
	connMu     sync.Mutex
}

func NewListener(cfg *config.Config, client *resty.Client, forwardURL, token string, txLogger *slog.Logger) *Listener {
	return &Listener{
		cfg:        cfg,
		client:     client,
		forwardURL: forwardURL,
		token:      token,
		txLogger:   txLogger,
	}
}

func (l *Listener) Listen(ctx context.Context, mock bool) error {
	if mock {
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
			mockData := map[string]any{
				"event": "billing.paid",
				"data": map[string]any{
					"id":         fmt.Sprintf("pix_char_%d", time.Now().Unix()),
					"externalId": "order_123",
					"amount":     1000,
					"status":     "PAID",
				},
			}

			message, _ := json.Marshal(mockData)
			l.logWebhook(message)

			if err := l.forward(ctx, message); err != nil {
				slog.Error("Mock webhook forward failed", "error", err)
			}
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

		l.logWebhook(message)

		g.Go(func() error {
			if err := l.forward(gCtx, message); err != nil {
				slog.Error("Webhook forward failed", "error", err)
			}
			return nil
		})
	}
}

func (l *Listener) SetupConn(conn *websocket.Conn) {
	conn.SetPongHandler(func(string) error {
		l.connMu.Lock()
		defer l.connMu.Unlock()
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})
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

func (l *Listener) logWebhook(message []byte) {
	var webhook struct {
		Event string `json:"event"`
		Data  struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(message, &webhook); err != nil {
		fmt.Println(string(message))
		return
	}

	style.LogWebhookReceived(webhook.Event, webhook.Data.ID)

	l.txLogger.Info("webhook_received",
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(message),
		"raw_message", string(message),
	)

	if l.cfg.Verbose {
		var buf bytes.Buffer
		if err := json.Indent(&buf, message, "", "  "); err == nil {
			fmt.Println(buf.String())

			return
		}
	}

	fmt.Println(string(message))
}

func (l *Listener) forward(ctx context.Context, message []byte) error {
	startTime := time.Now()

	resp, err := l.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(message).
		Post(l.forwardURL)

	duration := time.Since(startTime)

	if err != nil {
		l.txLogger.Error("webhook_forward_failed",
			"url", l.forwardURL,
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("failed to forward webhook: %w", err)
	}

	statusCode := resp.StatusCode()

	var webhook struct {
		Event string `json:"event"`
	}
	_ = json.Unmarshal(message, &webhook)

	style.LogWebhookForwarded(statusCode, http.StatusText(statusCode), webhook.Event)

	if statusCode < 200 || statusCode >= 300 {
		l.txLogger.Error("webhook_forward_error",
			"url", l.forwardURL,
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
			"response_body", string(resp.Body()),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("failed to forward webhook: server returned %d", statusCode)
	}

	l.txLogger.Info("webhook_forwarded",
		"url", l.forwardURL,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(message),
	)

	return nil
}

