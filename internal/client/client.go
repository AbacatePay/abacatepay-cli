package client

import (
	"github.com/go-resty/resty/v2"

	"github.com/AbacatePay/abacatepay-cli/internal/config"
)

func New(cfg *config.Config) *resty.Client {
	return resty.New().
		SetTimeout(cfg.HTTPTimeout).
		SetHeader("User-Agent", "github.com/AbacatePay/abacatepay-cli/1.0")
}
