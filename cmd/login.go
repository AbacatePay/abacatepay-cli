package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"signin"},
	Short:   "Sign in to AbacatePay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return login()
	},
}

var name, key string

func init() {
	loginCmd.Flags().StringVar(&name, "name", "", "Profile name")
	loginCmd.Flags().StringVar(&key, "key", "", "Your AbacatePay API key")

	rootCmd.AddCommand(loginCmd)
}

func login() error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}

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
		OpenBrowser: utils.OpenBrowser,
	}

	if err := auth.Login(params); err != nil {
		return err
	}

	return nil
}
