package payments

import (
	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/types"

	v1 "github.com/almeidazs/go-abacate-types/v1"
)

func (s *Service) CreatePixQRCode(body *v1.RESTPostCreateQRCodePixBody, isTrigger bool) (string, error) {
	var result types.PixResponse
	err := s.executeRequest(
		s.Client.R().SetBody(body),
		"POST",
		s.BaseURL+"/v1"+v1.RouteCreatePIXQRCode,
		&result,
	)
	if err != nil {
		return "", err
	}

	if !isTrigger {
		output.Print(output.Result{
			Title: "PIX Payment Created",
			Fields: map[string]string{
				"ID":     result.Data.ID,
				"Status": "PENDING",
			},
			Data: result,
		})
	}

	return result.Data.ID, nil
}

func (s *Service) CheckPixQRCode(id string) error {
	var result types.PixResponse
	err := s.executeRequest(
		s.Client.R().SetQueryParam("id", id),
		"GET",
		s.BaseURL+v1.RouteCheckQRCodePIX,
		&result,
	)
	if err != nil {
		return err
	}

	output.Print(output.Result{
		Title: "PIX Status Check",
		Fields: map[string]string{
			"ID":     id,
			"Status": result.Data.Status,
		},
		Data: result,
	})

	return nil
}

func (s *Service) SimulatePixQRCodePayment(id string, isTrigger bool) error {
	var result types.PixResponse
	err := s.executeRequest(
		s.Client.R().SetQueryParam("id", id),
		"POST",
		s.BaseURL+"/v1"+v1.RouteSimulatePayment,
		&result,
	)
	if err != nil {
		return err
	}

	if !isTrigger {
		output.Print(output.Result{
			Title: "PIX Payment Simulated",
			Fields: map[string]string{
				"ID":     result.Data.ID,
				"Status": result.Data.Status,
			},
			Data: result,
		})
	}

	return nil
}
