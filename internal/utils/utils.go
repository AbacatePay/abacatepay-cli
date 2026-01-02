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

	fmt.Println()
	slog.Info("Iniciando escuta de webhooks...", "forward_url", params.ForwardURL)
	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

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

	fmt.Printf("\nURL para encaminhar webhooks [%s]: ", defaultURL)
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

func ShowUpdate(currentVersion string) {
	fmt.Println("[DEBUG] Verificando atualizações para versão:", currentVersion)
	ctx := context.Background()

	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug("AbacatePay/abacatepay-cli"))
	if err != nil {
		fmt.Println("[DEBUG] Erro ao buscar update:", err)
		return
	}
	if !found {
		fmt.Println("[DEBUG] Nenhuma versão encontrada no GitHub")
		return
	}

	fmt.Printf("[DEBUG] Última versão encontrada: %s\n", latest.Version())

	if latest.LessOrEqual(currentVersion) {
		fmt.Println("[DEBUG] Versão atual já é a mais recente")
		return
	}

	green := "\033[32m"
	yellow := "\033[33m"
	reset := "\033[0m"
	bold := "\033[1m"

	fmt.Println()
	fmt.Println(green + "╭──────────────────────────────────────────────────────────────╮" + reset)
	fmt.Println(green + "│" + reset + "                                                              " + green + "│" + reset)
	fmt.Printf(green+"│"+reset+"   🥑 %sNova versão disponível:%s %-18s         "+green+"│"+reset+"\n", bold, reset, yellow+latest.Version()+reset)
	fmt.Printf(green+"│"+reset+"      Atual: %-41s "+green+"│"+reset+"\n", currentVersion)
	fmt.Println(green + "│" + reset + "                                                              " + green + "│" + reset)
	fmt.Printf(green + "│" + reset + "   Para atualizar execute:                                    " + green + "│" + reset + "\n")
	fmt.Printf(green+"│"+reset+"   %sabacatepay-cli update%s                                  "+green+"│"+reset+"\n", bold, reset)
	fmt.Println(green + "│" + reset + "                                                              " + green + "│" + reset)
	fmt.Println(green + "╰──────────────────────────────────────────────────────────────╯" + reset)
	fmt.Println()
}
