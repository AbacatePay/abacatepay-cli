package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/client"
	"abacatepay-cli/internal/config"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/webhook"
)

var (
	local   bool
	verbose bool
)

func main() {
	logCfg, err := logger.DefaultConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao configurar logger: %v\n", err)
		os.Exit(1)
	}

	if _, err := logger.Setup(logCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar logger: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd().Execute(); err != nil {
		slog.Error("Erro ao executar comando", "error", err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "abacatepay-cli",
		Short:   "AbacatePay CLI para executar webhooks localmente",
		Version: "1.0.0",
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Habilitar logs detalhados")
	cmd.PersistentFlags().BoolVarP(&local, "local", "l", false, "Usar servidor de teste")

	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	}

	cmd.AddCommand(
		loginCmd(),
		logoutCmd(),
		statusCmd(),
		listenCmd(),
	)

	return cmd
}

func getConfig() *config.Config {
	if local {
		return config.Local()
	}
	return config.Default()
}

func getStore(cfg *config.Config) auth.TokenStore {
	return auth.NewKeyringStore(cfg.ServiceName, cfg.TokenKey)
}

func promptForURL(defaultURL string) string {
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

func startListener(ctx context.Context, cfg *config.Config, cli *resty.Client, store auth.TokenStore, forwardURL string) error {
	token, err := store.Get()
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

	listener := webhook.NewListener(cfg, cli, forwardURL, token, txLogger)

	fmt.Println()
	slog.Info("Iniciando escuta de webhooks...", "forward_url", forwardURL)
	fmt.Println("Pressione Ctrl+C para parar")
	fmt.Println()

	return listener.Listen(ctx)
}

func loginCmd() *cobra.Command {
	var forwardURL string
	var skipListen bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Autenticar com AbacatePay e iniciar listener",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getConfig()
			cli := client.New(cfg)
			store := getStore(cfg)

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			if err := auth.Login(ctx, cfg, cli, store); err != nil {
				return err
			}

			if skipListen {
				return nil
			}

			if forwardURL == "" {
				forwardURL = promptForURL(cfg.DefaultForwardURL)
			}

			return startListener(ctx, cfg, cli, store, forwardURL)
		},
	}

	cmd.Flags().StringVarP(&forwardURL, "forward", "f", "", "URL para encaminhar webhooks")
	cmd.Flags().BoolVar(&skipListen, "no-listen", false, "Não iniciar listener após login")

	return cmd
}

func logoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Sair do AbacatePay",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getConfig()
			store := getStore(cfg)
			return auth.Logout(store)
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Verificar status da autenticação",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getConfig()
			store := getStore(cfg)

			token, err := store.Get()
			if err != nil {
				return err
			}

			if token != "" {
				slog.Info("Autenticado")
			} else {
				slog.Info("Não autenticado")
				fmt.Println("\nExecute 'abacatepay-cli login' para autenticar")
			}
			return nil
		},
	}
}

func listenCmd() *cobra.Command {
	var forwardURL string

	cmd := &cobra.Command{
		Use:   "listen",
		Short: "Escutar webhooks e encaminhar para servidor local",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := getConfig()
			store := getStore(cfg)

			token, err := store.Get()
			if err != nil {
				return err
			}

			if token == "" {
				return fmt.Errorf("não autenticado. Execute 'abacatepay-cli login' primeiro")
			}

			if forwardURL == "" {
				forwardURL = promptForURL(cfg.DefaultForwardURL)
			}

			logCfg, err := logger.DefaultConfig()
			if err != nil {
				return fmt.Errorf("erro ao configurar logger: %w", err)
			}

			txLogger, err := logger.NewTransactionLogger(logCfg)
			if err != nil {
				return fmt.Errorf("erro ao criar logger de transações: %w", err)
			}

			cli := client.New(cfg)
			listener := webhook.NewListener(cfg, cli, forwardURL, token, txLogger)

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			fmt.Println("Pressione Ctrl+C para parar")
			fmt.Println()

			return listener.Listen(ctx)
		},
	}

	cmd.Flags().StringVarP(&forwardURL, "forward", "f", "", "URL para encaminhar webhooks")

	return cmd
}
