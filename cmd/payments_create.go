package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/payments"
	"abacatepay-cli/internal/prompts"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/types"
	"abacatepay-cli/internal/utils"

	v1 "github.com/almeidazs/go-abacate-types/v1"
	"github.com/spf13/cobra"
)

var createInteractive bool

var createPaymentCmd = &cobra.Command{
	Use:   "create [pix|checkout]",
	Short: "Create a new payment charge",
	Long: `Create a new payment charge.
By default, creates a mock payment with random data.
Use -i to enter interactive mode and specify details.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var method string
		if len(args) > 0 {
			method = args[0]
		}

		return createPayment(method)
	},
}

func init() {
	paymentsCmd.AddCommand(createPaymentCmd)

	createPaymentCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Enable interactive mode")
}

func createPayment(method string) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	if method == "" {
		options := map[string]string{
			"PIX QR Code": "pix",
			"Checkout":    "checkout",
		}

		method, err = style.Select("🥑 Choose payment method\n", options)
		if err != nil {
			return err
		}
	}

	service := payments.New(deps.Client, deps.Config.APIBaseURL, Verbose)

	switch method {
	case "pix", "pix_qrcode":
		if !createInteractive {
			body := mock.CreatePixQRCodeMock()
			_, err := service.CreatePixQRCode(body, false)
			return err
		}

		body := &v1.RESTPostCreateQRCodePixBody{
			Customer: &v1.APICustomerMetadata{},
		}
		if err := prompts.PromptForPIXQRCodeData(body); err != nil {
			return fmt.Errorf("failed to prompt pix qrcode data: %w", err)
		}
		_, err := service.CreatePixQRCode(body, false)
		return err

	case "checkout":
		if !createInteractive {
			body := mock.CreateCheckoutMock()
			return service.CreateCheckout(body)
		}

		body := &types.CreateCheckoutRequest{
			Customer: &types.Customer{},
		}
		if err := prompts.PromptForCheckout(body); err != nil {
			return fmt.Errorf("failed to prompt checkout data: %w", err)
		}
		return service.CreateCheckout(body)

	default:
		return fmt.Errorf("invalid payment method: %s. Use 'pix' or 'checkout'", method)
	}
}
