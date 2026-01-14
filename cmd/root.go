package cmd

import (
	"log/slog"
	"os"

	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "abacatepay",
	Short:         "AbacatePay’s developer-first CLI for APIs and local workflows",
	Version:       version.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var Local, Verbose bool

func Exec() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&Local, "local", "l", false, "Use test server")

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

	if err := rootCmd.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
