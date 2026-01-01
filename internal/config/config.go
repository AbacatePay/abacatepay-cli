package config

import "time"

type Config struct {
	Verbose           bool
	APIBaseURL        string
	WebSocketBaseURL  string
	PollInterval      time.Duration
	MaxRetries        int
	ServiceName       string
	TokenKey          string
	HTTPTimeout       time.Duration
	RetryCount        int
	RetryWaitTime     time.Duration
	DefaultForwardURL string
}

func Default() *Config {
	return &Config{
		APIBaseURL:        "https://api.abacatepay.com",
		WebSocketBaseURL:  "wss://ws.abacatepay.com/ws",
		PollInterval:      2 * time.Second,
		MaxRetries:        30,
		ServiceName:       "abacatepay-cli",
		TokenKey:          "auth-token",
		HTTPTimeout:       10 * time.Second,
		RetryCount:        3,
		RetryWaitTime:     1 * time.Second,
		DefaultForwardURL: "http://localhost:3000/webhooks",
		Verbose:           false,
	}
}

func Local() *Config {
	cfg := Default()
	cfg.APIBaseURL = "http://191.252.202.128:8080"
	cfg.WebSocketBaseURL = "ws://191.252.202.128:8080/ws"
	return cfg
}
