package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"

	"abacatepay-cli/internal/config"
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

func (l *Listener) Listen(ctx context.Context) error {
	slog.Info("Iniciando listener de webhooks...")

	conn, err := l.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	slog.Info("Listener iniciado", "forward_url", l.forwardURL)

	return l.readLoop(ctx, conn)
}

func (l *Listener) connect(ctx context.Context) (*websocket.Conn, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+l.token)

	dialer := websocket.Dialer{
		HandshakeTimeout: l.cfg.HTTPTimeout,
	}

	conn, _, err := dialer.DialContext(ctx, l.cfg.WebSocketBaseURL, header)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar WebSocket: %w", err)
	}

	return conn, nil
}

func (l *Listener) readLoop(ctx context.Context, conn *websocket.Conn) error {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	conn.SetPongHandler(func(string) error {
		l.connMu.Lock()
		defer l.connMu.Unlock()
		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		return nil
	})

	g.Go(func() error {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-gCtx.Done():
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
					slog.Debug("Falha ao enviar ping", "error", err)
					return err
				}
			}
		}
	})

	for {
		select {
		case <-gCtx.Done():
			_ = g.Wait()
			return nil
		default:
		}

		l.connMu.Lock()
		conn.SetReadDeadline(time.Now().Add(90 * time.Second))
		l.connMu.Unlock()

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Info("Conexão WebSocket fechada normalmente")
				_ = g.Wait()
				return nil
			}

			if gCtx.Err() != nil {
				_ = g.Wait()
				return nil
			}

			_ = g.Wait()
			return fmt.Errorf("erro ao ler mensagem: %w", err)
		}

		l.logWebhook(message)

		g.Go(func() error {
			if err := l.forward(gCtx, message); err != nil {
				slog.Error("Falha ao encaminhar webhook", "error", err)
			}
			return nil
		})
	}
}

func (l *Listener) logWebhook(message []byte) {
	var webhook Message
	if err := json.Unmarshal(message, &webhook); err == nil {
		slog.Info("Webhook recebido", "event", webhook.Event, "size_bytes", len(message))

		l.txLogger.Info("webhook_received",
			"event", webhook.Event,
			"timestamp", time.Now().Format(time.RFC3339),
			"size_bytes", len(message),
			"data", string(webhook.Data),
		)
	} else {
		slog.Info("Webhook recebido", "size_bytes", len(message))

		l.txLogger.Info("webhook_received",
			"timestamp", time.Now().Format(time.RFC3339),
			"size_bytes", len(message),
			"raw_message", string(message),
		)
	}

	var prettyJSON interface{}
	if err := json.Unmarshal(message, &prettyJSON); err == nil {
		formatted, err := json.MarshalIndent(prettyJSON, "", "  ")
		if err != nil {
			fmt.Println(string(message))
		} else {
			fmt.Println(string(formatted))
		}
	} else {
		fmt.Println(string(message))
	}
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
		return fmt.Errorf("falha ao encaminhar: %w", err)
	}

	statusCode := resp.StatusCode()

	if statusCode < 200 || statusCode >= 300 {
		l.txLogger.Error("webhook_forward_error",
			"url", l.forwardURL,
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
			"response_body", string(resp.Body()),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("falha ao encaminhar webhook: status %d", statusCode)
	}

	l.txLogger.Info("webhook_forwarded",
		"url", l.forwardURL,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
		"size_bytes", len(message),
	)

	slog.Debug("Webhook encaminhado",
		"status", statusCode,
		"url", l.forwardURL,
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}
