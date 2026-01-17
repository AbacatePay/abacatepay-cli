package payments

import (
	"fmt"

	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/types"
)

func (s *Service) CreateCheckout(body *types.CreateCheckoutRequest) error {
	var result types.CheckoutResponse
	err := s.executeRequest(
		s.Client.R().SetBody(body),
		"POST",
		s.BaseURL+"/v2/checkouts/create",
		&result,
	)
	if err != nil {
		return err
	}

	output.Print(output.Result{
		Title: "Checkout Created",
		Fields: map[string]string{
			"ID":     result.Data.ID,
			"URL":    result.Data.URL,
			"Amount": fmt.Sprintf("R$ %.2f", float64(result.Data.Amount)/100),
			"Status": result.Data.Status,
		},
		Data: result,
	})

	return nil
}