package utils

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/webhook"

	"github.com/charmbracelet/lipgloss"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/go-resty/resty/v2"
)

type StartListenerParams struct {
	Context    context.Context
	Config     *config.Config
	Client     *resty.Client
	Store      auth.TokenStore
	ForwardURL string
	Version    string
}

type Dependencies struct {
	Config *config.Config
	Client *resty.Client
	Store  auth.TokenStore
}

func StartListener(params *StartListenerParams) error {
	go ShowUpdate(params.Version)

	activeProfile, err := params.Store.GetActiveProfile()

	if err != nil || activeProfile == "" {
		return fmt.Errorf("no active profile found")
	}

	token, err := params.Store.GetNamed(activeProfile)

	if err != nil || token == "" {
		return fmt.Errorf("couldn’t load token for profile '%s'", activeProfile)
	}

	logCfg, err := logger.DefaultConfig()

	if err != nil {
		return fmt.Errorf("failed to configure logger: %w", err)
	}

	txLogger, err := logger.NewTransactionLogger(logCfg)

	if err != nil {
		return fmt.Errorf("failed to initialize transaction logger: %w", err)
	}

	listener := webhook.NewListener(params.Config, params.Client, params.ForwardURL, token, txLogger)

	fmt.Fprintln(os.Stderr)
	slog.Info("Listening for webhooks", "forward_to", params.ForwardURL)
	fmt.Fprintln(os.Stderr, "Press Ctrl+C to stop")
	fmt.Fprintln(os.Stderr)

	return listener.Listen(params.Context)
}

func GetConfig(local bool) *config.Config {
	if local {
		return config.Local()
	}

	return config.Default()
}

func GetStore(cfg *config.Config) auth.TokenStore {
	return auth.NewKeyringStore(cfg.ServiceName, cfg.TokenKey)
}

func PromptForURL(defaultURL string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stderr, "\nForward events to [%s]: ", defaultURL)

	input, err := reader.ReadString('\n')

	if err != nil {
		return defaultURL
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return defaultURL
	}

	return input
}

func SetupDependencies(local bool, verbose bool) *Dependencies {
	cfg := GetConfig(local)
	cfg.Verbose = verbose

	cli := client.New(cfg)
	store := GetStore(cfg)

	return &Dependencies{
		Config: cfg,
		Client: cli,
		Store:  store,
	}
}

func CheckUpdate(ctx context.Context, currentVersion string) (*selfupdate.Release, bool, error) {
	slug := "AbacatePay/abacatepay-cli"
	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug(slug))

	if err != nil {
		return nil, false, err
	}

	if !found || latest.LessOrEqual(currentVersion) {
		return nil, false, nil
	}

	return latest, true, nil
}

func ShowUpdate(currentVersion string) {
	latest, found, _ := CheckUpdate(context.Background(), currentVersion)

	if !found {
		return
	}

	var (
		primaryColor = lipgloss.Color("#25D366")
		yellowColor  = lipgloss.Color("#FFFF00")

		boxStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginTop(1).
				MarginBottom(1)

		titleStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

		versionStyle = lipgloss.NewStyle().
				Foreground(yellowColor).
				Bold(true)

		commandStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)
	)

	msg := fmt.Sprintf(
		"🥑 %s %s\n      Current: %s\n\n   To update, run:\n   %s",
		titleStyle.Render("Update available:"),
		versionStyle.Render(latest.Version()),
		currentVersion,
		commandStyle.Render("abacatepay update"),
	)

	fmt.Fprintln(os.Stderr, boxStyle.Render(msg))
}
