package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current profile and authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return whoami()
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func whoami() error {
	deps := utils.SetupDependencies(Local, Verbose)
	activeProfile, err := deps.Store.GetActiveProfile()

	if err != nil || activeProfile == "" {
		return fmt.Errorf("no active profile found. Please login first")
	}


token, err := deps.Store.GetNamed(activeProfile)

	if err != nil || token == "" {
		return fmt.Errorf("token not found for active profile: %s", activeProfile)
	}

	user, err := auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return fmt.Errorf("session expired for profile %s: %w", activeProfile, err)
	}

	fmt.Printf("● Active Profile: %s\n", activeProfile)
	fmt.Printf("● User:           %s\n", user.Name)
	fmt.Printf("● Email:          %s\n", user.Email)
	fmt.Printf("● Status:         Authenticated ✅\n")

	return nil
}