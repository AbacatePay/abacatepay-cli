package cmd

import (
	"abacatepay-cli/internal/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var deleteProfileCmd = &cobra.Command{
	Use:     "delete [name]",
	Short:   "Remove a profile",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"remove", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return deleteProfile(args[0])
	},
}

func init() {
	profileCmd.AddCommand(deleteProfileCmd)
}

func deleteProfile(name string) error {
	deps := utils.SetupDependencies(Local, Verbose)

	token, err := deps.Store.GetNamed(name)

	if err != nil {
		return fmt.Errorf("error verifying profile: %w", err)
	}
	if token == "" {
		return fmt.Errorf("profile '%s' not found", name)
	}

	active, _ := deps.Store.GetActiveProfile()

	if active == name {
		return fmt.Errorf("cannot delete the active profile. Switch to another profile first using 'profile use'")
	}

	if err := deps.Store.DeleteNamed(name); err != nil {
		return fmt.Errorf("error deleting profile: %w", err)
	}

	fmt.Printf("Profile '%s' successfully removed.\n", name)

	return nil
}
