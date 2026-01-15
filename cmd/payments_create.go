package cmd

import (
	"fmt"
	"strconv"

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

	var body *v1.RESTPostCreateQRCodePixBody

	body = mock.CreatePixQRCodeMock()
	if createInteractive {
		body = &v1.RESTPostCreateQRCodePixBody{
			Customer: &v1.APICustomerMetadata{},
		}

		var amountStr string
		err = style.Input("Valor (em centavos, ex: 1000 para R$10,00)\n", "1000", &amountStr, func(s string) error {
			if _, err := strconv.ParseInt(s, 10, 64); err != nil {
				return fmt.Errorf("insira um número válido")
			}
			return nil
		})
		if err != nil {
			return err
		}
		amount, _ := strconv.ParseInt(amountStr, 10, 64)
		body.Amount = int(amount)

		var desc string
		err = style.Input("Descrição do Produto\n", "Minha compra", &desc, nil)
		if err != nil {
			return err
		}
		body.Description = &desc

		err = style.Input("Nome do Cliente\n", "João Silva", &body.Customer.Name, nil)
		if err != nil {
			return err
		}

		err = style.Input("Email do Cliente\n", "joao@exemplo.com", &body.Customer.Email, nil)
		if err != nil {
			return err
		}

		err = style.Input("CPF/CNPJ do Cliente\n", "12345678909", &body.Customer.TaxID, nil)
		if err != nil {
			return err
		}
	}

	switch method {
	case "pix":
		pixService := pix.New(deps.Client, deps.Config.APIBaseURL)
		return pixService.CreateQRCode(body)
	case "card":
		fmt.Println("🚧 Criação de pagamento via Cartão de Crédito em breve!")
		return nil
	default:
		return fmt.Errorf("método de pagamento desconhecido: %s", method)
	}
}

