package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/style"
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

func setupListenDeps() (*utils.Dependencies, string, error) {
	if listenMock {
		deps := utils.SetupDependencies(Local, Verbose)
		return deps, "", nil
	}

	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return nil, "", err
	}

	profile, _ := deps.Store.GetActiveProfile()
	token, _ := deps.Store.GetNamed(profile)
	return deps, token, nil
}

func listen(cmd *cobra.Command) error {
	deps, token, err := setupListenDeps()
	if err != nil {
		return err
	}

	defaultURL := "http://localhost:3000/webhooks/abacatepay"
	if !cmd.Flags().Changed("forward-to") {
		err := style.Input("Forward events to", defaultURL, &forwardURL, nil)
		if err != nil {
			return err
		}

		if forwardURL == "" {
			forwardURL = defaultURL
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer cancel()

	params := &utils.StartListenerParams{
		Context:    ctx,
		Config:     deps.Config,
		Client:     deps.Client,
		ForwardURL: forwardURL,
		Store:      deps.Store,
		Token:      token,
		Version:    cmd.Root().Version,
		Mock:       listenMock,
	}

	return utils.StartListener(params)
}
