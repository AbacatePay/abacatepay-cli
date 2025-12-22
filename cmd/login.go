package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Autenticar com AbacatePay e iniciar listener",
	RunE: func(cmd *cobra.Command, args []string) error {
		return login()
	},
}

var name, key string

func init() {
	loginCmd.Flags().StringVar(&key, "key", "", "Abacate Pay's API Key")
	loginCmd.Flags().StringVar(&name, "name", "", "Name for the profile (Min 3, Max 50 chars.)")

	rootCmd.AddCommand(loginCmd)
}

func login() error {
	cfg := utils.GetConfig(Local)
	cli := client.New(cfg)
	store := utils.GetStore(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := auth.Login(ctx, cfg, cli, store); err != nil {
		return err
	}

	return nil
}
