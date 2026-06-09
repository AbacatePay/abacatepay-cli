// Package cmd...
package cmd

import (
	"log/slog"
	"os"

	"github.com/AbacatePay/abacatepay-cli/internal/logger"
	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "abacatepay",
	Short:         "AbacatePay’s developer-first CLI for APIs and local workflows",
	Version:       version.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var (
	Local, Verbose bool
	OutputFormat   string
)

func Exec() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&Local, "local", "l", false, "Deprecated: API environment is determined by the API key")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format: text, json, table")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		format, err := output.ParseFormat(OutputFormat)
		if err != nil {
			return err
		}
		output.SetFormat(format)

		level := slog.LevelInfo
		if Verbose {
			level = slog.LevelDebug
		}

		cfg, err := logger.DefaultConfig()
		if err != nil {
			h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			slog.SetDefault(slog.New(h))
			return nil
		}

		cfg.Level = level
		if _, err := logger.Setup(cfg); err != nil {
			h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
			slog.SetDefault(slog.New(h))
		}
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}
