package cmd

import (
	"log/slog"
	"os"

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

		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		slog.SetDefault(slog.New(handler))
	}

	cobra.CheckErr(rootCmd.Execute())
}
