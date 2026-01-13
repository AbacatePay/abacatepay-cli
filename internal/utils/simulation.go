package utils

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type SimulationClient struct {
	cli   *resty.Client
	token string
}

func NewSimulationClient(resty *resty.Client, token string) *SimulationClient {
	return &SimulationClient{
		cli:   resty,
		token: token,
	}
}

func (s *SimulationClient) SimulateBillingPaid(ctx context.Context, billingID string) error {
	resp, err := s.cli.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.token).
		SetQueryParam("id", billingID).
		SetBody(map[string]any{"metadata": map[string]string{}}).
		Post("/v1/pixQrCode/simulate-payment")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("API Error: %d - %s", resp.StatusCode(), resp.String())
	}

	return nil
}
