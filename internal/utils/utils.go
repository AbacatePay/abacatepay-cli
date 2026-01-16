package utils

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/webhook"

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

func SetupClient(local, verbose bool) (*Dependencies, error) {
	if !IsOnline() {
		return nil, fmt.Errorf("you’re offline — check your connection and try again")
	}

	deps := SetupDependencies(local, verbose)
	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil {
		return nil, fmt.Errorf("failed to get active profile: %w", err)
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		return nil, fmt.Errorf("token not found for active profile: %s", activeProfile)
	}

	_, err = auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return nil, fmt.Errorf("session expired for profile %s: %w", activeProfile, err)
	}

	deps.Client.SetAuthToken(token)
	return deps, nil
}

func CheckUpdate(ctx context.Context, currentVersion string) (*selfupdate.Release, bool, error) {
	slug := "AbacatePay/abacatepay-cli"
	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug(slug))
	if err != nil {
		return nil, false, err
	}

	if !found || currentVersion == "" || latest.LessOrEqual(currentVersion) {
		return nil, false, nil
	}

	return latest, true, nil
}

func ShowUpdate(currentVersion string) {
	latest, found, _ := CheckUpdate(context.Background(), currentVersion)

	if !found {
		return
	}

	msg := fmt.Sprintf(
		"🥑 %s %s\n      Current: %s\n\n   To update, run:\n   %s",
		style.TitleStyle.Render("Update available:"),
		style.VersionStyle.Render(latest.Version()),
		currentVersion,
		style.CommandStyle.Render("abacatepay update"),
	)

	fmt.Fprintln(os.Stderr, style.BoxStyle.Render(msg))
}

func GenerateValidCPF(r *rand.Rand) string {
	digits := make([]int, 11)
	for i := range 9 {
		digits[i] = r.Intn(10)
	}

	sum := 0
	for i := range 9 {
		sum += digits[i] * (10 - i)
	}
	digits[9] = calculateDigit(sum)

	sum = 0
	for i := range 10 {
		sum += digits[i] * (11 - i)
	}
	digits[10] = calculateDigit(sum)

	cpf := ""
	for _, d := range digits {
		cpf += fmt.Sprintf("%d", d)
	}
	return cpf
}

func calculateDigit(sum int) int {
	remainder := (sum * 10) % 11
	if remainder >= 10 {
		return 0
	}
	return remainder
}
