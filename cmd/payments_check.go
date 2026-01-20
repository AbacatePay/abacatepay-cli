package cmd

import (
	"abacatepay-cli/internal/payments"

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
	return payments.ExecutePaymentAction(Local, Verbose, func(s *payments.Service) error {
		return s.CheckPixQRCode(paymentID)
	})
}
