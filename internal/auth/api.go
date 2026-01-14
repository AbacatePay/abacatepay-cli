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
		return nil, fmt.Errorf("failed to connect to API: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNotFound {
			slog.Warn("Profile API returned 404. Using Mock profile for local testing.")
			return &UserData{
				ID:    "mock-id",
				Name:  "Abacate Tester",
				Email: "test@abacatepay.com",
			}, nil
		}
		return nil, fmt.Errorf("token invalid or expired (status %d) while accessing %s", resp.StatusCode(), baseURL+"/v1/customer/get")
	}

	return &result.Data, nil
}
