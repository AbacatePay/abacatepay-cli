package auth

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type UserData struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	Data UserData `json:"data"`
}

func ValidateToken(client *resty.Client, baseURL, token string) (*UserData, error) {
	var result UserResponse
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetResult(&result).
		Get(baseURL + "/customer/get")

	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com a API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("token inválido ou expirado (status %d)", resp.StatusCode())
	}

	return &result.Data, nil
}
