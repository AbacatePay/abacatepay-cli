package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/payments/pix"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var checkPaymentCmd = &cobra.Command{
	Use:   "check",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return check(args[0])
	},
}

func init() {
	paymentsCmd.AddCommand(checkPaymentCmd)
}

func check(paymentID string) error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}

	deps := utils.SetupDependencies(Local, Verbose)
	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil {
		return fmt.Errorf("failed to get active profile: %w", err)
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		return fmt.Errorf("token not found for active profile: %s", activeProfile)
	}

	_, err = auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return fmt.Errorf("session expired for profile %s: %w", activeProfile, err)
	}

	deps.Client.SetAuthToken(token)

	pixService := pix.New(deps.Client, deps.Config.APIBaseURL)

	return pixService.CheckQRCode(paymentID)
}
