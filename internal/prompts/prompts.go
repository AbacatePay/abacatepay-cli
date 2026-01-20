package prompts

import (
	"fmt"
	"strconv"

	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/types"

	v1 "github.com/almeidazs/go-abacate-types/v1"
)

func PromptForPIXQRCodeData(body *v1.RESTPostCreateQRCodePixBody) error {
	var amountStr string
	err := style.Input("Amount (in cents, e.g. 1000 for R$10.00)", "1000", &amountStr, func(s string) error {
		if _, err := strconv.ParseInt(s, 10, 64); err != nil {
			return fmt.Errorf("please enter a valid number")
		}
		return nil
	})
	if err != nil {
		return err
	}
	amount, _ := strconv.ParseInt(amountStr, 10, 64)
	body.Amount = int(amount)

	var desc string
	if err := style.Input("Product Description", "My purchase", &desc, nil); err != nil {
		return err
	}
	body.Description = &desc

	if err := style.Input("Customer Name", "John Doe", &body.Customer.Name, nil); err != nil {
		return err
	}

	if err := style.Input("Customer Email", "john@example.com", &body.Customer.Email, nil); err != nil {
		return err
	}

	if err := style.Input("Customer TaxID (CPF/CNPJ)", "12345678909", &body.Customer.TaxID, nil); err != nil {
		return err
	}

	return nil
}

func PromptForCheckout(body *types.CreateCheckoutRequest) error {
	if err := style.Input("Customer Name", "John Doe", &body.Customer.Name, nil); err != nil {
		return err
	}

	if err := style.Input("Customer Email", "john@example.com", &body.Customer.Email, nil); err != nil {
		return err
	}

	if err := style.Input("Customer TaxID (CPF/CNPJ)", "12345678909", &body.Customer.TaxID, nil); err != nil {
		return err
	}
	var productID string
	if err := style.Input("Product ID", "prod_abc123xyz", &productID, nil); err != nil {
		return err
	}

	var qtdStr string
	if err := style.Input("Product Quantity", "1", &qtdStr, nil); err != nil {
		return err
	}
	qtd, err := strconv.Atoi(qtdStr)
	if err != nil {
		return err
	}

	body.Items = []types.Item{
		{
			ID:       productID,
			Quantity: qtd,
		},
	}

	return nil
}
