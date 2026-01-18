package cmd

import (
	"abacatepay-cli/internal/output"
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
		output.Error("You are not authenticated. Use 'abacatepay login' to start.")
		return nil
	}

	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil || activeProfile == "" {
		output.Error("No active profile found.")
		return nil
	}

	output.Print(output.Result{
		Title: "Connected successfully",
		Fields: map[string]string{
			"Profile": activeProfile,
			"Status":  "Online",
		},
		Data: map[string]string{
			"profile": activeProfile,
			"status":  "online",
		},
	})

	return nil
}
