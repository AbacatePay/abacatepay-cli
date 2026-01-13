package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"profile"},
	Short:   "Gerenciar perfis de autenticação salvos",
	Long:    "Permite listar, alternar e remover perfis de acesso do AbacatePay configurados nesta máquina.",
}

var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar todos os perfis configurados",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProfiles()
	},
}

var useProfileCmd = &cobra.Command{
	Use:   "use [nome]",
	Short: "Alternar para outro perfil existente",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return useProfile(args[0])
	},
}

var deleteProfileCmd = &cobra.Command{
	Use:     "delete [nome]",
	Aliases: []string{"remove", "rm"},
	Short:   "Remover um perfil",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return deleteProfile(args[0])
	},
}

func init() {
	profileCmd.AddCommand(listProfilesCmd)
	profileCmd.AddCommand(useProfileCmd)
	profileCmd.AddCommand(deleteProfileCmd)

	rootCmd.AddCommand(profileCmd)
}

func listProfiles() error {
	deps := utils.SetupDependencies(Local, Verbose)

	profiles, err := deps.Store.List()
	if err != nil {
		return fmt.Errorf("erro ao listar perfis: %w", err)
	}

	active, err := deps.Store.GetActiveProfile()
	if err != nil {
		active = ""
	}

	if len(profiles) == 0 {
		fmt.Println("Nenhum perfil encontrado. Use 'abacatepay-cli login' para criar um.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  NOME\tSTATUS")

	for _, p := range profiles {
		status := ""
		prefix := "  "
		if p == active {
			status = "(ativo)"
			prefix = "* "
		}
		fmt.Fprintf(w, "%s%s\t%s\n", prefix, p, status)
	}
	w.Flush()

	return nil
}

func useProfile(name string) error {
	deps := utils.SetupDependencies(Local, Verbose)

	token, err := deps.Store.GetNamed(name)
	if err != nil {
		return fmt.Errorf("erro ao buscar perfil: %w", err)
	}
	if token == "" {
		return fmt.Errorf("perfil '%s' não encontrado", name)
	}

	if err := deps.Store.SetActiveProfile(name); err != nil {
		return fmt.Errorf("erro ao definir perfil ativo: %w", err)
	}

	fmt.Printf("Agora usando o perfil: %s\n", name)
	return nil
}

func deleteProfile(name string) error {
	deps := utils.SetupDependencies(Local, Verbose)

	token, err := deps.Store.GetNamed(name)
	if err != nil {
		return fmt.Errorf("erro ao verificar perfil: %w", err)
	}
	if token == "" {
		return fmt.Errorf("perfil '%s' não encontrado", name)
	}

	active, _ := deps.Store.GetActiveProfile()
	if active == name {
		return fmt.Errorf("não é possível deletar o perfil ativo. Mude para outro perfil primeiro com 'profiles use'")
	}

	if err := deps.Store.DeleteNamed(name); err != nil {
		return fmt.Errorf("erro ao deletar perfil: %w", err)
	}

	fmt.Printf("Perfil '%s' removido com sucesso.\n", name)
	return nil
}

