package payments

import (
	"abacatepay-cli/internal/style"
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
		style.PrintSuccess("PIX Payment Created", map[string]string{
			"ID":     result.Data.ID,
			"Status": "PENDING",
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

	style.PrintSuccess("PIX Status Check", map[string]string{
		"ID":     id,
		"Status": result.Data.Status,
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
		style.PrintSuccess("PIX Payment Simulated", map[string]string{
			"ID":     result.Data.ID,
			"Status": result.Data.Status,
		})
	}

	return nil
}
