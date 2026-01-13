package cmd

import (
	"log/slog"
	"os"

	"abacatepay-cli/internal/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "abacatepay-cli",
	Short:   "AbacatePay CLI para executar webhooks localmente",
	Version: "0.0.1",
}
var Local, Verbose bool

func Exec() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Habilitar logs detalhados")
	rootCmd.PersistentFlags().BoolVarP(&Local, "local", "l", false, "Usar servidor de teste")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		level := slog.LevelInfo
		if Verbose {
			level = slog.LevelDebug
		}

		cfg, err := logger.DefaultConfig()
		if err != nil {
			h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			slog.SetDefault(slog.New(h))
			return
		}

		cfg.Level = level

		if _, err := logger.Setup(cfg); err != nil {
			h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			slog.SetDefault(slog.New(h))
		}
	}

	cobra.CheckErr(rootCmd.Execute())
}
