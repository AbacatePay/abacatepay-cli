package config

import "time"

type Config struct {
	Verbose           bool
	APIBaseURL        string
	WebSocketBaseURL  string
	ServiceName       string
	TokenKey          string
	HTTPTimeout       time.Duration
	DefaultForwardURL string
}

func Default() *Config {
	return &Config{
		APIBaseURL:        "https://api.abacatepay.com",
		WebSocketBaseURL:  "wss://ws.abacatepay.com/ws",
		ServiceName:       "abacatepay-cli",
		TokenKey:          "auth-token",
		HTTPTimeout:       15 * time.Second,
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
