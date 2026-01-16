package payments

import (
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/types"
	"encoding/json"
	"fmt"
	"net/http"

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
		return s.handleAPIError(resp)
	}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

func (s *Service) handleAPIError(resp *resty.Response) error {
	var apiErr types.APIError
	_ = json.Unmarshal(resp.Body(), &apiErr)

	errorMessage := apiErr.Message
	if errorMessage == "" {
		errorMessage = fmt.Sprintf("Unexpected error (Status %d)", resp.StatusCode())
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		style.PrintError("Session expired or invalid API Key. Please login again.")
	case http.StatusTooManyRequests:
		style.PrintError("Too many requests. Please slow down.")
	case http.StatusInternalServerError:
		style.PrintError("Server error. Please try again later.")
	default:
		style.PrintError(errorMessage)
	}

	return fmt.Errorf("API Error: %s", errorMessage)
}
