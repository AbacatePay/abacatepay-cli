package payments

import (
	"abacatepay-cli/internal/style"

	v1 "github.com/almeidazs/go-abacate-types/v1"
)

type pixResponse struct {
	Data struct {
		ID     string `json:"id"`
		BRCode string `json:"brCode"`
		Status string `json:"status"`
	} `json:"data"`
}

func (s *Service) CreatePixQRCode(body *v1.RESTPostCreateQRCodePixBody) error {
	var result pixResponse
	err := s.executeRequest(
		s.Client.R().SetBody(body),
		"POST",
		s.BaseURL+v1.RouteCreatePIXQRCode,
		&result,
	)
	if err != nil {
		return err
	}

	style.PrintSuccess("PIX Payment Created", map[string]string{
		"ID":     result.Data.ID,
		"Status": "PENDING",
	})

	return nil
}

func (s *Service) CheckPixQRCode(id string) error {
	var result pixResponse
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

func (s *Service) SimulatePixQRCodePayment(id string) error {
	var result pixResponse
	err := s.executeRequest(
		s.Client.R().SetQueryParam("id", id),
		"POST",
		s.BaseURL+v1.RouteSimulatePayment,
		&result,
	)
	if err != nil {
		return err
	}

	style.PrintSuccess("PIX Payment Simulated", map[string]string{
		"ID":     result.Data.ID,
		"Status": result.Data.Status,
	})

	return nil
}