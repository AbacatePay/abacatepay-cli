package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AbacatePay/abacatepay-cli/internal/logger"
	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/tui"
	"github.com/AbacatePay/abacatepay-cli/internal/utils"
	"github.com/AbacatePay/abacatepay-cli/internal/webhook"
	"github.com/AbacatePay/abacatepay-cli/internal/ws"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for webhooks and forward them to your local app",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listen(cmd)
	},
}

var (
	forwardURL string
	listenEnv  string
)

func init() {
	listenCmd.Flags().StringVar(&forwardURL, "forward-to", "", "Where incoming events should be sent")
	listenCmd.Flags().StringVar(&listenEnv, "env", "dev", "Which environment events to receive: dev, prod, or all")

	rootCmd.AddCommand(listenCmd)
}

func listen(cmd *cobra.Command) error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	url, err := utils.GetForwardURL(forwardURL, utils.DefaultForwardURL)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// The live dashboard needs a real terminal to render into, and only
	// makes sense for the default text output — piped/scripted/CI use (or
	// -o json/table) keeps the plain sequential-line behavior.
	if tui.IsInteractive() && output.GetFormat() == output.FormatText {
		return listenInteractive(ctx, cancel, deps, url)
	}

	params := &utils.StartListenerParams{
		Context:    ctx,
		Config:     deps.Config,
		Client:     deps.Client,
		ForwardURL: url,
		Store:      deps.Store,
		Token:      deps.Token,
		Version:    cmd.Root().Version,
		Env:        listenEnv,
	}

	return utils.StartListener(params)
}

func listenInteractive(ctx context.Context, cancel context.CancelFunc, deps *utils.Dependencies, forwardURL string) error {
	txLogger, err := utils.SetupTransactionLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize transaction logger: %w", err)
	}

	listener := webhook.NewListener(deps.Config, deps.Client, forwardURL, deps.Token, txLogger)

	// Background slog output (connection retries, heartbeat warnings, ...)
	// would otherwise corrupt the alt-screen dashboard; it still reaches
	// the file logger.
	logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(true)

	model := tui.NewListenModel(forwardURL, cancel)
	program := tea.NewProgram(model, tea.WithAltScreen())

	events := make(chan webhook.Event, 64)
	listener.Emit = func(e webhook.Event) { events <- e }
	listener.OnStatus = func(s ws.Status) { program.Send(tui.ConnStatusMsg{Status: s}) }

	go func() {
		for e := range events {
			program.Send(e)
		}
	}()

	go func() {
		listenErr := listener.Listen(ctx)
		close(events)

		if errors.Is(listenErr, context.Canceled) || errors.Is(listenErr, context.DeadlineExceeded) {
			listenErr = nil
		}
		program.Send(tui.StoppedMsg{Err: listenErr})
	}()

	_, runErr := program.Run()
	cancel()

	return runErr
}
