package cmd

import (
	"fmt"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage your local profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return profiles()
	},
}

var format string

// NOTE: I`ll add subcommands here, similar to CRUD operations to manage profiles`
func init() {
	loginCmd.Flags().StringVar(&format, "format", "json", "Output format")

	rootCmd.AddCommand(profileCmd)
}

func profiles() error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}

	// deps := utils.SetupDependencies(Local, Verbose)

	return nil
}
