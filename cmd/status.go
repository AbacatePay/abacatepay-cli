package cmd

import (
	"log/slog"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verificar status da autenticação",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getAuthStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getAuthStatus() error {
	deps := utils.SetupDependencies(Local, Verbose)

	activeProfile, err := deps.Store.GetActiveProfile()
	if err != nil || activeProfile == "" {
		slog.Info("Não autenticado (nenhum perfil ativo)")
		return nil
	}

	token, err := deps.Store.GetNamed(activeProfile)
	if err != nil || token == "" {
		slog.Info("Não autenticado", "profile", activeProfile)
		return nil
	}

	slog.Info("Autenticado", "profile", activeProfile)
	return nil
}