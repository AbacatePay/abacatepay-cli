package cmd

import (
	"fmt"

	"abacatepay-cli/internal/auth"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Exibir o perfil atual e status de autenticação",
	RunE: func(cmd *cobra.Command, args []string) error {
		return whoami()
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func whoami() error {
	deps := utils.SetupDependencies(Local, Verbose)
	activeProfile, err := deps.Store.GetActiveProfile()

	if err != nil || activeProfile == "" {
		return fmt.Errorf("nenhum perfil ativo encontrado. Por favor, faça login primeiro")
	}

	token, err := deps.Store.GetNamed(activeProfile)

	if err != nil || token == "" {
		return fmt.Errorf("token não encontrado para o perfil ativo: %s", activeProfile)
	}

	user, err := auth.ValidateToken(deps.Client, deps.Config.APIBaseURL, token)
	if err != nil {
		return fmt.Errorf("sessão expirada para o perfil %s: %w", activeProfile, err)
	}

	fmt.Printf("● Perfil Ativo: %s\n", activeProfile)
	fmt.Printf("● Usuário:      %s\n", user.Name)
	fmt.Printf("● Email:        %s\n", user.Email)
	fmt.Printf("● Status:       Autenticado ✅\n")

	return nil
}
