package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/AbacatePay/abacatepay-cli/internal/config"
	"github.com/AbacatePay/abacatepay-cli/internal/store"
)

func newTestClient() *resty.Client {
	return resty.New().SetTimeout(5 * time.Second)
}

func writeJSON(w http.ResponseWriter, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func TestLoginWithDeviceFlow_HappyPath(t *testing.T) {
	const requestedID = "cli_test123"
	pollCount := 0

	mux := http.NewServeMux()
	mux.HandleFunc("/app/oauth/cli/request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["device"]; !ok {
			t.Error("expected a device field in the request body")
		}

		writeJSON(w, map[string]any{
			"success": true,
			"data":    map[string]string{"publicId": requestedID},
			"error":   nil,
		})
	})
	mux.HandleFunc("/app/oauth/cli/"+requestedID+"/status", func(w http.ResponseWriter, r *http.Request) {
		pollCount++
		if pollCount < 2 {
			writeJSON(w, map[string]any{
				"success": true,
				"data":    map[string]string{"status": "PENDING"},
				"error":   nil,
			})
			return
		}

		writeJSON(w, map[string]any{
			"success": true,
			"data":    map[string]string{"status": "APPROVED", "token": "session-token-123"},
			"error":   nil,
		})
	})
	mux.HandleFunc("/v2/stores/get", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"success": true,
			"data":    map[string]string{"id": "store_1", "name": "Test Store"},
			"error":   nil,
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	memStore := store.NewMemoryStore()
	params := &LoginParams{
		Config: &config.Config{
			APIBaseURL: server.URL,
			AppBaseURL: "https://app.test",
		},
		Client:      newTestClient(),
		Store:       memStore,
		Context:     context.Background(),
		OpenBrowser: func(string) error { return nil },
	}

	if err := Login(params); err != nil {
		t.Fatalf("expected login to succeed, got error: %v", err)
	}

	token, err := memStore.GetNamed("default")
	if err != nil {
		t.Fatalf("failed to read saved token: %v", err)
	}
	if token != "session-token-123" {
		t.Errorf("expected saved token %q, got %q", "session-token-123", token)
	}

	activeProfile, err := memStore.GetActiveProfile()
	if err != nil || activeProfile != "default" {
		t.Errorf("expected active profile %q, got %q (err: %v)", "default", activeProfile, err)
	}
}

func TestPollForToken_Rejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"success": true,
			"data":    map[string]string{"status": "REJECTED"},
			"error":   nil,
		})
	}))
	defer server.Close()

	cfg := &config.Config{APIBaseURL: server.URL}
	if _, err := pollForTokenLoop(context.Background(), cfg, newTestClient(), "cli_rejected"); err == nil {
		t.Fatal("expected an error for a rejected login request")
	}
}

func TestPollForToken_Expired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{
			"success": true,
			"data":    map[string]string{"status": "EXPIRED"},
			"error":   nil,
		})
	}))
	defer server.Close()

	cfg := &config.Config{APIBaseURL: server.URL}
	if _, err := pollForTokenLoop(context.Background(), cfg, newTestClient(), "cli_expired"); err == nil {
		t.Fatal("expected an error for an expired login request")
	}
}

func TestPollForToken_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &config.Config{APIBaseURL: server.URL}
	if _, err := pollForTokenLoop(context.Background(), cfg, newTestClient(), "cli_missing"); err == nil {
		t.Fatal("expected an error for a not-found login request")
	}
}
