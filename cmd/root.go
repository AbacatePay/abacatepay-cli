package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "abacatepay-cli",
	Short:   "AbacatePay CLI para executar webhooks localmente",
	Version: "1.0.0",
}

func init() {
}

func Exec() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Habilitar logs detalhados")
	rootCmd.PersistentFlags().BoolVarP(&local, "local", "l", false, "Usar servidor de teste")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	}

	cobra.CheckErr(rootCmd.Execute())
}
