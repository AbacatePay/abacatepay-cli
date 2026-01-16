package cmd

import (
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"doctor"},
	Short:   "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getAuthStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getAuthStatus() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		style.PrintError("You are not authenticated. Use 'abacatepay login' to start.")
		return nil
	}

	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil || activeProfile == "" {
		style.PrintError("No active profile found.")
		return nil
	}

	style.PrintSuccess("Connected successfully", map[string]string{
		"Profile": activeProfile,
		"Status":  "Online",
	})

	return nil
}
