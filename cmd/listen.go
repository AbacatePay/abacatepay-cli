package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"
	"abacatepay-cli/internal/version"

	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for webhooks and forward them to your local app",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listen(cmd)
	},
}

var forwardURL string

func init() {
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "", "Where incoming events should be sent")

	rootCmd.AddCommand(listenCmd)
}

func listen(cmd *cobra.Command) error {
	deps := utils.SetupDependencies(Local, Verbose)

	activeProfile, err := deps.Store.GetActiveProfile()

	if err != nil || activeProfile == "" {
		return fmt.Errorf("you’re not logged in — run `abacatepay login` to continue")
	}

	token, err := deps.Store.GetNamed(activeProfile)

	if err != nil || token == "" {
		return fmt.Errorf("this profile doesn’t have a valid token — try logging in again")
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
		Version:    version.Version,
	}
	if err := utils.StartListener(params); err != nil {
		return fmt.Errorf("couldn’t start the webhook listener: %w", err)
	}

	fmt.Printf("Forwarding events to %s\n\n", forwardURL)
	fmt.Println("Press Ctrl+C to stop")

	go func() {
		<-ctx.Done()
		fmt.Println("\nListener stopped")
	}()

	return nil
}
