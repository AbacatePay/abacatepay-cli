package config

import (
	"os"
	"time"
)

type Config struct {
	Verbose           bool
	APIBaseURL        string
	AppBaseURL        string
	WebSocketBaseURL  string
	ServiceName       string
	TokenKey          string
	HTTPTimeout       time.Duration
	DefaultForwardURL string
}

// envOr returns the value of the given environment variable, or fallback if unset.
// These overrides exist for local testing against a non-production stack; they are
// not documented CLI flags.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Default() *Config {
	return &Config{
		APIBaseURL:        envOr("ABACATEPAY_API_URL", "https://api.abacatepay.com"),
		AppBaseURL:        envOr("ABACATEPAY_APP_URL", "https://app.abacatepay.com"),
		WebSocketBaseURL:  envOr("ABACATEPAY_WS_URL", "wss://ws.abacatepay.com/ws"),
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
