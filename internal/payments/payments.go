package payments

import (
	"encoding/json"
	"fmt"

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

func (s *Service) executeRequest(req *resty.Request, method, url string, result any) error {
	var resp *resty.Response
	var err error

	switch method {
	case "POST":
		resp, err = req.Post(url)
	case "GET":
		resp, err = req.Get(url)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}