package cmd

import (
	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:     "logout",
	Aliases: []string{"signout"},
	Short:   "Sign out of AbacatePay",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func logout() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	return auth.Logout(deps.Store)
}
