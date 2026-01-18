package cmd

import (
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
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	pixService := payments.New(deps.Client, deps.Config.APIBaseURL)

	return pixService.SimulatePixQRCodePayment(paymentID, false)
}
