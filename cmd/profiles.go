package cmd

import (
	"fmt"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profiles",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return profiles()
	},
}

var format string

// NOTE: I`ll add subcommands here, similar to CRUD operations to manage profiles`
func init() {
	loginCmd.Flags().StringVar(&format, "format", "json", "Visualization format")

	rootCmd.AddCommand(profileCmd)
}

func profiles() error {
	if !utils.IsOnline() {
		return fmt.Errorf("you're offline, please stabilish your connection to continue")
	}

	// deps := utils.SetupDependencies(Local, Verbose)

	return nil
}
