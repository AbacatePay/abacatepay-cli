package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/payments"
	"abacatepay-cli/internal/prompts"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"

	v1 "github.com/almeidazs/go-abacate-types/v1"
	"github.com/spf13/cobra"
)

var createInteractive bool

var createPaymentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new payment charge",
	Long: `Create a new payment charge.
By default, creates a mock payment with random data.
Use -i to enter interactive mode and specify details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return createPayment()
	},
}

func init() {
	paymentsCmd.AddCommand(createPaymentCmd)

	createPaymentCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Enable interactive mode")
}

func createPayment() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	options := map[string]string{
		"PIX QR Code": "pix_qrcode",
		"Checkout ":   "checkout",
	}

	method, err := style.Select("🥑 Escolha o método de pagamento\n", options)
	if err != nil {
		return err
	}

	switch method {
	case "pix_qrcode":
		body := &v1.RESTPostCreateQRCodePixBody{
			Customer: &v1.APICustomerMetadata{},
		}
		pixService := payments.New(deps.Client, deps.Config.APIBaseURL)

		if createInteractive {
			if err := prompts.PromptForPIXQRCodeData(body); err != nil {
				return fmt.Errorf("error to prompt pix qrcode data: %w", err)
			}

			return pixService.CreatePixQRCode(body)
		}

		body = mock.CreatePixQRCodeMock()
		return pixService.CreatePixQRCode(body)

	case "checkout":
		fmt.Println("🚧 Criação de pagamento via Cartão de Crédito em breve!")
		return nil

	default:
		return nil
	}
}
