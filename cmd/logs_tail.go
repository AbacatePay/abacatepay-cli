package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/utils"
	"abacatepay-cli/internal/webhook"

	"github.com/spf13/cobra"
)

var logsTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Stream live webhook events in real-time",
	Long:  "Connect to the WebSocket and display incoming webhook events as they arrive",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logsTail()
	},
}

func init() {
	logsCmd.AddCommand(logsTailCmd)
}

func logsTail() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	listener := webhook.NewTailListener(deps.Config, deps.Config.TokenKey)

	fmt.Println("Streaming webhook events...")
	fmt.Println("\nPress Ctrl+C to stop")

	go func() {
		<-ctx.Done()
		fmt.Println("\nListener stopped")
	}()

	err = listener.Listen(ctx)
	if err == nil {
		return nil
	}
	if ctx.Err() != nil {
		return nil
	}
	return fmt.Errorf("couldn't start the tail listener: %w", err)
}
