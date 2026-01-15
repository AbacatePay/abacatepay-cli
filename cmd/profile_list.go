package cmd

import (
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"
	"fmt"

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
		fmt.Println("No profiles found. Use 'abacatepay login' to create one.")
		return nil
	}

	profileMap := make(map[string]string)
	for _, p := range profiles {
		token, _ := deps.Store.GetNamed(p)
		profileMap[p] = token
	}

	style.ProfileSimpleList(profileMap, active)

	return nil
}
