package cmd

import (
	"abacatepay-cli/internal/payments"

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
	return payments.ExecutePaymentAction(Local, Verbose, func(s *payments.Service) error {
		return s.SimulatePixQRCodePayment(paymentID, false)
	})
}
