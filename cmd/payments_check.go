package cmd

import (
	"abacatepay-cli/internal/payments"
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
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	pixService := payments.New(deps.Client, deps.Config.APIBaseURL)

	return pixService.CheckPixQRCode(paymentID)
}
