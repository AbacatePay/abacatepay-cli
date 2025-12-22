package cmd

import (
	"abacatepay-cli/internal/auth"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sair do AbacatePay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logout()
	},
}

func logout() error {
	cfg := getConfig()
	store := getStore(cfg)
	return auth.Logout(store)
}
