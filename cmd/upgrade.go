package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/AbacatePay/abacatepay-cli/internal/version"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Atualiza a CLI para a nova versão caso esteja disponível",
	RunE: func(cmd *cobra.Command, args []string) error {
		return upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() error {
	ctx := context.Background()

	latest, found, err := version.CheckUpdate(ctx, version.Version)
	if err != nil {
		return fmt.Errorf("não foi possível checar por atualizações: %w", err)
	}

	if !found {
		fmt.Printf("Você já está na última versão (%s)\n", version.Version)
		return nil
	}

	fmt.Printf("Atualização disponível: %s\n", latest.Version())
	fmt.Println("Baixando e instalando...")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("não foi possível obter o caminho do executável: %w", err)
	}

	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("falha na atualização: %w", err)
	}

	fmt.Println("Atualização concluída ✨")

	return nil
}
