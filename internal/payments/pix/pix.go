package pix

import (
	"fmt"

	v1 "github.com/almeidazs/go-abacate-types/v1"
	"github.com/go-resty/resty/v2"
)

type CreatePixQRCodeParams struct {
	Client  *resty.Client
	Body    *v1.RESTPostCreateQRCodePixBody
	BaseURL string
	Token   string
}

func CreateQRCode(params *CreatePixQRCodeParams) error {
	resp, err := params.Client.R().
		SetBody(params.Body).
		SetHeader("Authorization", "Bearer "+params.Token).
		Post(params.BaseURL + v1.RouteCreatePIXQRCode)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), resp.String())
	}

	fmt.Println("✅ Mock PIX Payment Created!")
	fmt.Println(resp.String())

	return nil
}
