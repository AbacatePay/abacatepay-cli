package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"abacatepay-cli/internal/webhook"
)

func StartListener(params *StartListenerParams) error {
	txLogger, err := SetupTransactionLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction logger: %w", err)
	}

	listener := webhook.NewListener(params.Config, params.Client, params.ForwardURL, params.Token, txLogger)

	fmt.Fprintln(os.Stderr)
	if params.Mock {
		slog.Info("Running in MOCK mode", "interval", "5s")
	}
	slog.Info("Listening for webhooks", "forward_to", params.ForwardURL)
	fmt.Fprintln(os.Stderr, "Press Ctrl+C to stop")
	fmt.Fprintln(os.Stderr)

	err = listener.Listen(params.Context, params.Mock)

	fmt.Fprintln(os.Stderr)
	slog.Info("Listener stopped")

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil
		}
		return err
	}

	return nil
}
