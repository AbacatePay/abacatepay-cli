package cmd

import (
	"context"
	"fmt"
	"os"

	"abacatepay-cli/internal/utils"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update CLI to new version if is available",
	RunE: func(cmd *cobra.Command, args []string) error {
		return update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() error {
	ctx := context.Background()

	latest, found, err := utils.CheckUpdate(ctx, rootCmd.Version)
	if err != nil {
		return fmt.Errorf("erro ao verificar atualizações: %w", err)
	}

	if !found {
		fmt.Printf("Você já está na versão mais recente (%s)\n", rootCmd.Version)
		return nil
	}

	fmt.Printf("Nova versão encontrada: %s\n", latest.Version())
	fmt.Println("Baixando e instalando atualização...")

	exe, _ := os.Executable()
	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("erro ao atualizar: %w", err)
	}

	fmt.Println("Atualizado com sucesso!")
	return nil
}
