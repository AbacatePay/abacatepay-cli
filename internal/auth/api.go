package auth

import (
	"abacatepay-cli/internal/types"

	"github.com/go-resty/resty/v2"
)

// TODO: I`ll made this func when i get the right endpoint later
func ValidateToken(client *resty.Client, baseURL, token string) (*types.User, error) {
	return &types.User{Name: "Mock User", Email: "mock@example.com"}, nil
}

