package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/utils"

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
	listenMock bool
)

func init() {
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "", "Where incoming events should be sent")
	listenCmd.Flags().BoolVar(&listenMock, "mock", false, "Simulate incoming webhooks without connecting to the API")

	rootCmd.AddCommand(listenCmd)
}

func listen(cmd *cobra.Command) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	url, err := utils.GetForwardURL(forwardURL)
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
		Token:      deps.Config.TokenKey,
		Version:    cmd.Root().Version,
		Mock:       listenMock,
	}

	return utils.StartListener(params)
}
