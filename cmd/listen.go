package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AbacatePay/abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for webhooks and forward them to your local app",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listen(cmd)
	},
}

var (
	forwardURL string
	listenEnv  string
)

func init() {
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "", "Where incoming events should be sent")
	listenCmd.Flags().StringVar(&listenEnv, "env", "dev", "Which environment events to receive: dev, prod, or all")

	rootCmd.AddCommand(listenCmd)
}

func listen(cmd *cobra.Command) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	url, err := utils.GetForwardURL(forwardURL, utils.DefaultForwardURL)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	params := &utils.StartListenerParams{
		Context:    ctx,
		Config:     deps.Config,
		Client:     deps.Client,
		ForwardURL: url,
		Store:      deps.Store,
		Token:      deps.Token,
		Version:    cmd.Root().Version,
		Env:        listenEnv,
	}

	return utils.StartListener(params)
}
