package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/payments"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var simulatePaymentCmd = &cobra.Command{
	Use:   "simulate",
	Args:  cobra.ExactArgs(1),
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return simulate(args[0])
	},
}

func init() {
	paymentsCmd.AddCommand(simulatePaymentCmd)
}

func simulate(paymentID string) error {
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

	pixService := payments.New(deps.Client, deps.Config.APIBaseURL)

	return pixService.SimulatePixQRCodePayment(paymentID)
}
