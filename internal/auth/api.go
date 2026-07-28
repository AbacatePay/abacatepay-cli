// Package auth validates AbacatePay credentials.
package auth

import (
	"fmt"
	"net/http"

	"github.com/AbacatePay/abacatepay-cli/internal/types"

	"github.com/go-resty/resty/v2"
)

func ValidateToken(client *resty.Client, baseURL, token string) (*types.User, error) {
	var result types.StoreResponse

	resp, err := client.R().
		SetAuthToken(token).
		SetResult(&result).
		Get(baseURL + "/v2/stores/get")
	if err != nil {
		return nil, fmt.Errorf("failed to validate API key: %w", err)
	}

	if resp.StatusCode() == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key")
	}

	if resp.IsError() {
		if result.Error != "" {
			return nil, fmt.Errorf("failed to validate API key: %s", result.Error)
		}
		return nil, fmt.Errorf("failed to validate API key: %s", resp.Status())
	}

	name := result.Data.Name
	if name == "" {
		name = result.Data.ID
	}
	if name == "" {
		name = "AbacatePay user"
	}

	return &types.User{Name: name}, nil
}
