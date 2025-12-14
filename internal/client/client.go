package client

import (
	"github.com/go-resty/resty/v2"

	"abacatepay-cli/internal/config"
)

func New(cfg *config.Config) *resty.Client {
	return resty.New().
		SetTimeout(cfg.HTTPTimeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWaitTime).
		SetHeader("User-Agent", "abacatepay-cli/1.0")
}
