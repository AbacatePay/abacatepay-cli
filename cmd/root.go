package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "abacatepay-cli",
	Short:   "AbacatePay CLI para executar webhooks localmente",
	Version: "0.0.0",
}
var Local, Verbose bool

func Exec() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Habilitar logs detalhados")
	rootCmd.PersistentFlags().BoolVarP(&Local, "local", "l", false, "Usar servidor de teste")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if Verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	}

	cobra.CheckErr(rootCmd.Execute())
}
