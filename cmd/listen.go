package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/utils"

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
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "http://localhost:3000/webhooks/abacatepay", "salve")

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

	cli := client.New(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	params := utils.StartListenerParams{
		Context:    ctx,
		Config:     cfg,
		Client:     cli,
		ForwardURL: forwardURL,
		Store:      store,
	}
	if err := utils.StartListener(&params); err != nil {
		return fmt.Errorf("error to start listener: %w", err)
	}

	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return nil
}
