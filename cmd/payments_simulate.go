package cmd

import (
	"github.com/AbacatePay/abacatepay-cli/internal/payments"

	"github.com/spf13/cobra"
)

var simulatePaymentCmd = &cobra.Command{
	Use:     "simulate <charge-id>",
	Aliases: []string{"simulate-payment"},
	Args:    cobra.ExactArgs(1),
	Short:   "Simulate payment for a transparent checkout in dev mode",
	Long: "Simulates payment for a transparent checkout created with a dev-mode API key. " +
		"This calls the API v2 endpoint POST /v2/transparents/simulate-payment?id=<charge-id>.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return simulate(args[0])
	},
}

func init() {
	paymentsCmd.AddCommand(simulatePaymentCmd)
}

func simulate(paymentID string) error {
	return payments.ExecutePaymentAction(Local, Verbose, func(s *payments.Service) error {
		return s.SimulateTransparentPayment(paymentID)
	})
}
