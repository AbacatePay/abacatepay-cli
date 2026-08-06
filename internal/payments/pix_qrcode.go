package payments

import (
	"github.com/AbacatePay/abacatepay-cli/internal/clierr"
	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/types"
)

func (s *Service) SimulateTransparentPayment(id string) error {
	action := func() (output.Result, error) {
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
			return output.Result{}, err
		}

		return output.Result{
			Title: "Payment Simulated",
			Fields: map[string]string{
				"ID":      result.Data.ID,
				"Status":  result.Data.Status,
				"DevMode": boolLabel(result.Data.DevMode),
			},
			Data: result,
		}, nil
	}

	// Verbose mode prints its own request/response dump from inside
	// executeRequest; a spinner animating concurrently would race with (and
	// garble) that output, so skip the TUI and go straight through.
	if s.Verbose {
		r, err := action()
		if err != nil {
			if !clierr.AlreadyDisplayed(err) {
				output.Error(err.Error())
			}
			return err
		}
		output.Print(r)
		return nil
	}

	return output.RunTask("Simulating payment...", action)
}

func boolLabel(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
