package payments

import (
	"encoding/json"
	"fmt"
	"net/http"

	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/types"

	"github.com/go-resty/resty/v2"
)

type Service struct {
	Client  *resty.Client
	BaseURL string
	Verbose bool
}

func New(client *resty.Client, baseURL string, verbose bool) *Service {
	return &Service{
		Client:  client,
		BaseURL: baseURL,
		Verbose: verbose,
	}
}

func (s *Service) executeRequest(req *resty.Request, method, url string, result any) error {
	if s.Verbose {
		fmt.Printf("Request: %s %s\n", method, url)
		if body := req.Body; body != nil {
			if b, ok := body.([]byte); ok {
				var pretty any
				if err := json.Unmarshal(b, &pretty); err == nil {
					style.PrintJSON(pretty)
				} else {
					fmt.Println(string(b))
				}
			} else {
				style.PrintJSON(body)
			}
		}
		fmt.Println()
	}

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

	if s.Verbose {
		fmt.Printf("Response: %s\n", resp.Status())
	}

	if resp.IsError() {
		return s.handleAPIError(resp)
	}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if s.Verbose {
		style.PrintJSON(result)
		fmt.Println()
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
		output.Error("Session expired or invalid API Key. Please login again.")
	case http.StatusTooManyRequests:
		output.Error("Too many requests. Please slow down.")
	case http.StatusInternalServerError:
		output.Error("Server error. Please try again later.")
	default:
		output.Error(errorMessage)
	}

	return fmt.Errorf("API Error: %s", errorMessage)
}
