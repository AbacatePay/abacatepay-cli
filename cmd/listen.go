package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/utils"
	"abacatepay-cli/internal/webhook"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Escutar webhooks e encaminhar para servidor local",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listen()
	},
}

var forwardURL string

func init() {
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "http://localhost:3000", "salve")

	rootCmd.AddCommand(listenCmd)
}

func listen() error {
	cfg := utils.GetConfig(Local)
	store := utils.GetStore(cfg)

	token, err := store.Get()
	if err != nil {
		return err
	}

	if token == "" {
		return fmt.Errorf("não autenticado. Execute 'abacatepay-cli login' primeiro")
	}

	if forwardURL == "" {
		forwardURL = utils.PromptForURL(cfg.DefaultForwardURL)
	}

	logCfg, err := logger.DefaultConfig()
	if err != nil {
		return fmt.Errorf("erro ao configurar logger: %w", err)
	}

	txLogger, err := logger.NewTransactionLogger(logCfg)
	if err != nil {
		return fmt.Errorf("erro ao criar logger de transações: %w", err)
	}

	cli := client.New(cfg)
	listener := webhook.NewListener(cfg, cli, forwardURL, token, txLogger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return listener.Listen(ctx)
}
