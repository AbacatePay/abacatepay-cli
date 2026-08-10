package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/AbacatePay/abacatepay-cli/internal/style"
	"github.com/AbacatePay/abacatepay-cli/internal/webhook"
)

func StartListener(params *StartListenerParams) error {
	txLogger, err := SetupTransactionLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction logger: %w", err)
	}

	listener := webhook.NewListener(params.Config, params.Client, params.ForwardURL, params.Token, txLogger)

	// User-facing status goes to stderr (styled) so piped stdout stays clean
	// webhook-event data; slog.Info here is the file-only diagnostic trail.
	fmt.Fprintln(os.Stderr)
	slog.Info("Listening for webhooks", "forward_to", params.ForwardURL)
	style.FprintInfo(os.Stderr, fmt.Sprintf("Listening for webhooks — forwarding to %s", params.ForwardURL))
	style.FprintInfo(os.Stderr, "Press Ctrl+C to stop")
	fmt.Fprintln(os.Stderr)

	err = listener.Listen(params.Context)

	fmt.Fprintln(os.Stderr)
	slog.Info("Listener stopped")
	style.FprintInfo(os.Stderr, "Listener stopped")

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}
		return err
	}

	return nil
}
