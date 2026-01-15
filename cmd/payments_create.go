package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/payments/pix"
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

	pixService := pix.New(deps.Client, deps.Config.APIBaseURL)

	options := map[string]string{
		"PIX QR Code":       "pix",
		"Cartão de Crédito": "card",
	}

	method, err := style.Select("🥑 Escolha o método de pagamento\n", options)
	if err != nil {
		return err
	}

	if !createInteractive {
		body := mock.CreatePixQRCodeMock()
		return pixService.CreateQRCode(body)
	}

	if method == "card" {
		fmt.Println("🚧 Criação de pagamento via Cartão de Crédito em breve!")
		return nil
	}

	body := &v1.RESTPostCreateQRCodePixBody{
		Customer: &v1.APICustomerMetadata{},
	}

	if err := prompts.PromptForPIXQRCodeData(body); err != nil {
		return fmt.Errorf("error to prompt pix qrcode data: %w", err)
	}

	return pixService.CreateQRCode(body)
}
