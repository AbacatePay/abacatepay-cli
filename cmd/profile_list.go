package cmd

import (
	"fmt"

	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProfiles()
	},
}

func init() {
	profileCmd.AddCommand(listProfilesCmd)
}

func listProfiles() error {
	deps := utils.SetupDependencies(Local, Verbose)

	profiles, err := deps.Store.List()
	if err != nil {
		return fmt.Errorf("error listing profiles: %w", err)
	}

	active, err := deps.Store.GetActiveProfile()
	if err != nil {
		active = ""
	}

	if len(profiles) == 0 {
		output.Print(output.Result{
			Title: "No profiles found",
			Fields: map[string]string{
				"Hint": "Use 'abacatepay login' to create one.",
			},
			Data: map[string]any{
				"profiles": []string{},
				"active":   "",
			},
		})
		return nil
	}

	outputProfiles := make([]output.Profile, 0, len(profiles))
	for _, name := range profiles {
		token, _ := deps.Store.GetNamed(name)
		outputProfiles = append(outputProfiles, output.Profile{
			Name:   name,
			Token:  token,
			Active: name == active,
		})
	}

	output.PrintProfiles(outputProfiles, active)
	return nil
}
