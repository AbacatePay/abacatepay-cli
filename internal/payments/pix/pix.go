package pix

import (
	"encoding/json"
	"fmt"

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

func (s *Service) CreateQRCode(body *v1.RESTPostCreateQRCodePixBody) error {
	resp, err := s.Client.R().
		SetBody(body).
		Post(s.BaseURL + v1.RouteCreatePIXQRCode)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	var result pixResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.ID == "" {
		return fmt.Errorf("no ID found in response")
	}

	fmt.Println("\n✅ Mock PIX Payment Created!")
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("PIX ID: %s\n", result.Data.ID)

	return nil
}

func (s *Service) CheckQRCode(id string) error {
	resp, err := s.Client.R().
		SetQueryParam("id", id).
		Get(s.BaseURL + v1.RouteCheckQRCodePIX)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	var result pixResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.ID == "" {
		return fmt.Errorf("no ID found in response")
	}

	fmt.Println("---------------------------------------------------------")
	fmt.Printf("PIX STATUS: %s\n", result.Data.Status)

	return nil
}

func (s *Service) SimulateQRCodePayment(id string) error {
	resp, err := s.Client.R().
		SetQueryParam("id", id).
		Post(s.BaseURL + v1.RouteSimulatePayment)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	var result pixResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Data.ID == "" {
		return fmt.Errorf("no ID found in response")
	}

	fmt.Println("---------------------------------------------------------")
	fmt.Printf("PIX ID: %s\n", result.Data.ID)
	fmt.Printf("PIX STATUS: %s\n", result.Data.Status)

	return nil
}

