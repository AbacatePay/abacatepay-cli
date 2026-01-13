package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show the current profile",
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
		return fmt.Errorf("you’re not signed in")
	}

	token, err := deps.Store.GetNamed(activeProfile)

	if err != nil || token == "" {
		return fmt.Errorf("this profile doesn’t have a valid session", activeProfile)
	}

	user, err := auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return fmt.Errorf("session expired for %s — please sign in again", activeProfile)
	}

	fmt.Printf("Profile: %s\n", activeProfile)
	fmt.Printf("User:    %s\n", user.Name)
	fmt.Printf("Email:   %s\n", user.Email)
	fmt.Printf("Status:  Signed in ✓\n")

	return nil
}
