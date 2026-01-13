package cmd

import (
	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out of AbacatePay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func logout() error {
	deps := utils.SetupDependencies(Local, Verbose)

	return auth.Logout(deps.Store)
}
