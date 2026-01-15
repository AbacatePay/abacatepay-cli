package ws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestConnectWithRetry_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	cfg := Config{
		URL:        wsURL,
		MinBackoff: 10 * time.Millisecond,
		MaxBackoff: 50 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	connectedCount := 0
	var mu sync.Mutex

	handler := func(ctx context.Context, conn *websocket.Conn) error {
		mu.Lock()
		connectedCount++
		mu.Unlock()
		return nil
	}

	err := ConnectWithRetry(ctx, cfg, handler)

	if err != context.DeadlineExceeded {
		t.Errorf("Expected error context.DeadlineExceeded, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if connectedCount == 0 {
		t.Error("The handler should have been called at least once")
	}
}

func TestConnectWithRetry_ConnectionRefused(t *testing.T) {
	cfg := Config{
		URL:        "ws://localhost:54321",
		MinBackoff: 1 * time.Millisecond,
		MaxBackoff: 5 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := ConnectWithRetry(ctx, cfg, func(ctx context.Context, conn *websocket.Conn) error {
		t.Error("Should not have connected")
		return nil
	})

	if err != context.DeadlineExceeded {
		t.Errorf("Expected error context.DeadlineExceeded when canceling retries, got: %v", err)
	}
}