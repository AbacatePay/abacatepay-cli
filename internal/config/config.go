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
		DefaultForwardURL: "http://localhost:3000/webhooks/abacatepay",
		Verbose:           false,
	}
}

func Local() *Config {
	// Kept for backwards compatibility with the old --local flag. API v2 uses the
	// same public endpoint for production and dev mode; the API key determines
	// the environment.
	return Default()
}
