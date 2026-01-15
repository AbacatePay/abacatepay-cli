package pix

import (
	"encoding/json"
	"fmt"

	"abacatepay-cli/internal/style"

	v1 "github.com/almeidazs/go-abacate-types/v1"
	"github.com/go-resty/resty/v2"
)

type Service struct {
	Client  *resty.Client
	BaseURL string
}

func New(client *resty.Client, baseURL string) *Service {
	return &Service{
		Client:  client,
		BaseURL: baseURL,
	}
}

type pixResponse struct {
	Data struct {
		ID     string `json:"id"`
		BRCode string `json:"brCode"`
		Status string `json:"status"`
	} `json:"data"`
}

func (s *Service) sendRequest(req *resty.Request, method, url string) (*pixResponse, error) {
	var resp *resty.Response
	var err error

	switch method {
	case "POST":
		resp, err = req.Post(url)
	case "GET":
		resp, err = req.Get(url)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	var result pixResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.ID == "" {
		return nil, fmt.Errorf("no ID found in response")
	}

	return &result, nil
}

func (s *Service) CreateQRCode(body *v1.RESTPostCreateQRCodePixBody) error {
	result, err := s.sendRequest(
		s.Client.R().SetBody(body),
		"POST",
		s.BaseURL+v1.RouteCreatePIXQRCode,
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

func (s *Service) CheckQRCode(id string) error {
	result, err := s.sendRequest(
		s.Client.R().SetQueryParam("id", id),
		"GET",
		s.BaseURL+v1.RouteCheckQRCodePIX,
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

func (s *Service) SimulateQRCodePayment(id string) error {
	result, err := s.sendRequest(
		s.Client.R().SetQueryParam("id", id),
		"POST",
		s.BaseURL+v1.RouteSimulatePayment,
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