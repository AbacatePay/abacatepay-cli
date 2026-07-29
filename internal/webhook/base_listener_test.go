package webhook

import (
	"net/url"
	"testing"

	"github.com/AbacatePay/abacatepay-cli/internal/config"
)

func TestWSConfig_PutsTokenInSessionQueryParam(t *testing.T) {
	listener := &BaseListener{
		Cfg:   &config.Config{WebSocketBaseURL: "wss://ws.abacatepay.com/ws"},
		Token: "session-token-123",
	}

	cfg := listener.WSConfig()

	parsed, err := url.Parse(cfg.URL)
	if err != nil {
		t.Fatalf("failed to parse dial URL: %v", err)
	}

	if got := parsed.Query().Get("session"); got != "session-token-123" {
		t.Errorf("expected session query param %q, got %q (full URL: %s)", "session-token-123", got, cfg.URL)
	}

	if cfg.Headers.Get("Authorization") != "" {
		t.Errorf("expected no Authorization header, got %q", cfg.Headers.Get("Authorization"))
	}
}
