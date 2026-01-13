package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "http://localhost:3000/webhooks/abacatepay", "URL local para ouvir eventos")

	rootCmd.AddCommand(listenCmd)
}

func listen() error {
	deps := utils.SetupDependencies(Local, Verbose)

	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil || activeProfile == "" {
		return fmt.Errorf("não autenticado. Execute 'abacatepay-cli login' primeiro")
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		return fmt.Errorf("token não encontrado para o perfil %s. Por favor, faça login novamente", activeProfile)
	}

	if forwardURL == "" || forwardURL == "http://localhost:3000/webhooks/abacatepay" {
		forwardURL = utils.PromptForURL(deps.Config.DefaultForwardURL)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	params := &utils.StartListenerParams{
		Context:    ctx,
		Config:     deps.Config,
		Client:     deps.Client,
		ForwardURL: forwardURL,
		Store:      deps.Store,
		Version:    rootCmd.Version,
	}
	if err := utils.StartListener(params); err != nil {
		return fmt.Errorf("error to start listener: %w", err)
	}

	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return nil
}