package webhook

import (
	"encoding/json"
	"log/slog"

	"github.com/AbacatePay/abacatepay-cli/internal/config"

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
	client     *resty.Client
	forwardURL string
	txLogger   *slog.Logger
	env        string
}

func NewListener(cfg *config.Config, client *resty.Client, forwardURL, token, env string, txLogger *slog.Logger) *Listener {
	return &Listener{
		BaseListener: BaseListener{
			Cfg:   cfg,
			Token: token,
		},
		client:     client,
		forwardURL: forwardURL,
		txLogger:   txLogger,
		env:        env,
	}
}
