package webhook

import (
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"abacatepay-cli/internal/config"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
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