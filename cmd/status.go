package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verificar status da autenticação",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getAuthStatus()
	},
}

func getAuthStatus() error {
	cfg := getConfig()
	store := getStore(cfg)

	token, err := store.Get()
	if err != nil {
		return err
	}

	if token != "" {
		slog.Info("Autenticado")
		return nil
	}

	slog.Info("Não autenticado")
	fmt.Println("\nExecute 'abacatepay-cli login' para autenticar")

	return nil
}
