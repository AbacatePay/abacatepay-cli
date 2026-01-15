package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/payments/pix"
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
		return createPayment(cmd)
	},
}

func init() {
	paymentsCmd.AddCommand(createPaymentCmd)

	createPaymentCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Enable interactive mode")
}

func createPayment(cmd *cobra.Command) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	method := "pix"
	options := map[string]string{
		"PIX QR Code":       "pix",
		"Cartão de Crédito": "card",
	}

	selected, err := style.Select("🥑 Escolha o método de pagamento\n", options)
	if err != nil {
		return err
	}

	method = selected
	if createInteractive {
		salve := []string{"Amount", "ExpiresIn"}
		body := &v1.RESTPostCreateQRCodePixBody{}
		for i, field := range salve {
			body[field] = "salve"
		}
	}

	switch method {
	case "pix":
		b := mock.CreatePixQRCodeMock()
		pixService := pix.New(deps.Client, deps.Config.APIBaseURL)
		return pixService.CreateQRCode(b)
	case "card":
		fmt.Println("🚧 Criação de pagamento via Cartão de Crédito em breve!")
		return nil
	default:
		return fmt.Errorf("método de pagamento desconhecido: %s", method)
	}
}

