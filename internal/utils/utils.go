package utils

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/webhook"

	"github.com/go-resty/resty/v2"
)

type StartListenerParams struct {
	Context    context.Context
	Config     *config.Config
	Client     *resty.Client
	Store      auth.TokenStore
	ForwardURL string
}

func StartListener(params *StartListenerParams) error {
	token, err := params.Store.Get()
	if err != nil {
		return fmt.Errorf("erro ao recuperar token: %w", err)
	}

	logCfg, err := logger.DefaultConfig()
	if err != nil {
		return fmt.Errorf("erro ao configurar logger: %w", err)
	}

	txLogger, err := logger.NewTransactionLogger(logCfg)
	if err != nil {
		return fmt.Errorf("erro ao criar logger de transações: %w", err)
	}

	listener := webhook.NewListener(params.Config, params.Client, params.ForwardURL, token, txLogger)

	fmt.Println()
	slog.Info("Iniciando escuta de webhooks...", "forward_url", params.ForwardURL)
	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return listener.Listen(params.Context)
}

func GetConfig(local bool) *config.Config {
	if local {
		return config.Local()
	}
	return config.Default()
}

func GetStore(cfg *config.Config) auth.TokenStore {
	return auth.NewKeyringStore(cfg.ServiceName, cfg.TokenKey)
}

func PromptForURL(defaultURL string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\nURL para encaminhar webhooks [%s]: ", defaultURL)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultURL
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultURL
	}

	return input
}

func Salve() {
}
