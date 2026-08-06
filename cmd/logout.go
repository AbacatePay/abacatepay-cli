package cmd

import (
	"github.com/AbacatePay/abacatepay-cli/internal/auth"
	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/utils"

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
	deps := utils.SetupDependencies(Local, Verbose)

	return output.RunTask("Signing out...", func() (output.Result, error) {
		profile, err := auth.Logout(deps.Store)
		if err != nil {
			return output.Result{}, err
		}

		return output.Result{
			Title: "Signed out successfully",
			Fields: map[string]string{
				"Profile": profile,
			},
		}, nil
	})
}
