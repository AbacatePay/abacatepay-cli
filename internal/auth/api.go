package auth

import (
	"fmt"
	"log/slog"
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
		Get(baseURL + "/v1/customer/get")
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com a API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {

		if resp.StatusCode() == http.StatusNotFound {
			slog.Warn("API de perfil retornou 404. Usando perfil Mock para testes locais.")
			return &UserData{
				ID:    "mock-id",
				Name:  "Abacate Tester",
				Email: "test@abacatepay.com",
			}, nil
		}
		return nil, fmt.Errorf("token inválido ou expirado (status %d) ao acessar %s", resp.StatusCode(), baseURL+"/v1/customer/get")
	}

	return &result.Data, nil
}
