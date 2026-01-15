package payments

import (
	"abacatepay-cli/internal/style"
	"fmt"
)

type CreateCheckoutRequest struct {
	Items         []Item    `json:"items"`
	Method        string    `json:"method,omitempty"` // PIX | CARD (default: PIX)
	ReturnURL     string    `json:"returnUrl,omitempty"`
	CompletionURL string    `json:"completionUrl,omitempty"`
	CustomerID    string    `json:"customerId,omitempty"`
	Customer      *Customer `json:"customer,omitempty"`
	Coupons       []string  `json:"coupons,omitempty"`
	ExternalID    string    `json:"externalId,omitempty"`
}

type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Customer struct {
	Name      string `json:"name,omitempty"`
	Cellphone string `json:"cellphone,omitempty"`
	Email     string `json:"email,omitempty"`
	TaxID     string `json:"taxId,omitempty"`
}

type checkoutResponse struct {
	Data struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Status string `json:"status"`
		Amount int    `json:"amount"`
	} `json:"data"`
}

func (s *Service) CreateCheckout(body *CreateCheckoutRequest) error {
	var result checkoutResponse
	err := s.executeRequest(
		s.Client.R().SetBody(body),
		"POST",
		s.BaseURL+"/v2/checkouts/create",
		&result,
	)
	if err != nil {
		return err
	}

	style.PrintSuccess("Checkout Created", map[string]string{
		"ID":     result.Data.ID,
		"URL":    result.Data.URL,
		"Amount": fmt.Sprintf("R$ %.2f", float64(result.Data.Amount)/100),
		"Status": result.Data.Status,
	})

	return nil
}