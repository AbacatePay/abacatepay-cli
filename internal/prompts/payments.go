package prompts

import (
	"fmt"
	"strconv"

	"abacatepay-cli/internal/style"

	v1 "github.com/almeidazs/go-abacate-types/v1"
)

func PromptForPIXQRCodeData(body *v1.RESTPostCreateQRCodePixBody) error {
	var amountStr string
	err := style.Input("Valor (em centavos, ex: 1000 para R$10,00)", "1000", &amountStr, func(s string) error {
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
	if err := style.Input("Descrição do Produto", "Minha compra", &desc, nil); err != nil {
		return err
	}
	body.Description = &desc

	if err := style.Input("Nome do Cliente", "João Silva", &body.Customer.Name, nil); err != nil {
		return err
	}

	if err := style.Input("Email do Cliente", "joao@exemplo.com", &body.Customer.Email, nil); err != nil {
		return err
	}

	if err := style.Input("CPF/CNPJ do Cliente", "12345678909", &body.Customer.TaxID, nil); err != nil {
		return err
	}

	return nil
}
