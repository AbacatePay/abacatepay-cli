package webhook

import (
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"time"

	"abacatepay-cli/internal/config"

	"github.com/go-resty/resty/v2"
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
	BaseListener
	client        *resty.Client
	forwardURL    string
	txLogger      *slog.Logger
	signingSecret string
}

func NewListener(cfg *config.Config, client *resty.Client, forwardURL, token string, txLogger *slog.Logger) *Listener {
	return &Listener{
		BaseListener: BaseListener{
			Cfg:   cfg,
			Token: token,
		},
		client:        client,
		forwardURL:    forwardURL,
		txLogger:      txLogger,
		signingSecret: "whsec_mock_" + hex.EncodeToString([]byte(time.Now().Format("150405"))),
	}
}

