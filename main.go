package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"abacatepay-cli/cmd"
	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/webhook"
)

var (
	local   bool
	verbose bool
)

func main() {
	cmd.Exec()

	logCfg, err := logger.DefaultConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao configurar logger: %v\n", err)
		os.Exit(1)
	}

	if _, err := logger.Setup(logCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd().Execute(); err != nil {
		slog.Error("Erro ao executar comando", "error", err)
		os.Exit(1)
	}
}

func getConfig() *config.Config {
	if local {
		return config.Local()
	}
	return config.Default()
}

func getStore(cfg *config.Config) auth.TokenStore {
	return auth.NewKeyringStore(cfg.ServiceName, cfg.TokenKey)
}

func promptForURL(defaultURL string) string {
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

func startListener(ctx context.Context, cfg *config.Config, cli *resty.Client, store auth.TokenStore, forwardURL string) error {
	token, err := store.Get()
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

	listener := webhook.NewListener(cfg, cli, forwardURL, token, txLogger)

	fmt.Println()
	slog.Info("Iniciando escuta de webhooks...", "forward_url", forwardURL)
	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return listener.Listen(ctx)
}

func loginCmd() *cobra.Command {
	cmd.Flags().StringVarP(&forwardURL, "forward", "f", "", "URL para encaminhar webhooks")
	cmd.Flags().BoolVar(&skipListen, "no-listen", false, "Não iniciar listener após login")

	return cmd
}
