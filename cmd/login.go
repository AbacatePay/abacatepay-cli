package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Autenticar com AbacatePay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return login()
	},
}

var name, key string

func init() {
	loginCmd.Flags().StringVar(&key, "key", "", "Abacate Pay's API Key")
	loginCmd.Flags().StringVar(&name, "name", "", "Abacate Pay's Profile Name")

	rootCmd.AddCommand(loginCmd)
}

func login() error {
	deps := utils.SetupDependencies(Local, Verbose)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	params := &auth.LoginParams{
		Config:      deps.Config,
		Client:      deps.Client,
		Store:       deps.Store,
		Context:     ctx,
		APIKey:      key,
		ProfileName: name,
	}

	if err := auth.Login(params); err != nil {
		return err
	}

	return nil
}
