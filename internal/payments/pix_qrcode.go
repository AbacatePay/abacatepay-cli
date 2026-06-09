package payments

import (
	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/types"
)

func (s *Service) SimulateTransparentPayment(id string) error {
	var result types.TransparentPaymentResponse
	err := s.executeRequest(
		s.Client.R().
			SetQueryParam("id", id).
			SetBody(types.SimulatePaymentRequest{Metadata: map[string]any{}}),
		"POST",
		s.BaseURL+"/v2/transparents/simulate-payment",
		&result,
	)
	if err != nil {
		return err
	}

	output.Print(output.Result{
		Title: "Payment Simulated",
		Fields: map[string]string{
			"ID":      result.Data.ID,
			"Status":  result.Data.Status,
			"DevMode": boolLabel(result.Data.DevMode),
		},
		Data: result,
	})

	return nil
}

func boolLabel(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
