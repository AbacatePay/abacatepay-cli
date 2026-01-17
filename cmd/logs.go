package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/utils"
	"abacatepay-cli/internal/webhook"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:       "logs",
	Short:     "Display recent webhook transaction logs",
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"tail", "list"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return logs(args)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func logs(args []string) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer cancel()

	txLogger, err := utils.SetupTra
	webhook.NewListener(deps.Config, deps.Client, forwardURL, deps.Config.TokenKey, txLogger)

	return nil
}
