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

	token, err := params.Store.Get()
	if err != nil {
		return fmt.Errorf("erro ao recuperar token: %w", err)
	}

	logCfg, err := logger.DefaultConfig()
	if err != nil {
		return fmt.Errorf("erro ao configurar logger: %w", err)
	}

	txLogger, err := logger.NewTransactionLogger(logCfg)
	if err != nil {
		return fmt.Errorf("erro ao criar logger de transações: %w", err)
	}

	listener := webhook.NewListener(params.Config, params.Client, params.ForwardURL, token, txLogger)

	fmt.Fprintln(os.Stderr)
	slog.Info("Iniciando escuta de webhooks...", "forward_url", params.ForwardURL)
	fmt.Fprintln(os.Stderr, "Pressione Ctrl+C para parar")
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

	fmt.Fprintf(os.Stderr, "\nURL para encaminhar webhooks [%s]: ", defaultURL)
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

	// Estilos do Lipgloss
	var (
		primaryColor = lipgloss.Color("#25D366") // Verde Abacate
		yellowColor  = lipgloss.Color("#FFFF00") // Amarelo Destaque

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

	// Construção da mensagem
	msg := fmt.Sprintf(
		"🥑 %s %s\n      Atual: %s\n\n   Para atualizar execute:\n   %s",
		titleStyle.Render("Nova versão disponível:"),
		versionStyle.Render(latest.Version()),
		currentVersion,
		commandStyle.Render("abacatepay-cli update"),
	)

	fmt.Fprintln(os.Stderr, boxStyle.Render(msg))
}
