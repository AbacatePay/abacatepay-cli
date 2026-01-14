package cmd

import (
	"fmt"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to another existing profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return switchProfile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func switchProfile(name string) error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}

	deps := utils.SetupDependencies(Local, Verbose)

	token, err := deps.Store.GetNamed(name)
	if err != nil {
		return fmt.Errorf("error searching for profile: %w", err)
	}
	if token == "" {
		return fmt.Errorf("profile '%s' not found", name)
	}

	if err := deps.Store.SetActiveProfile(name); err != nil {
		return fmt.Errorf("error setting active profile: %w", err)
	}

	fmt.Printf("Now using profile: %s\n", name)
	return nil
}
