package cmd

import (
	"fmt"

	"abacatepay-cli/internal/scaffold"
	"abacatepay-cli/internal/style"

	"github.com/spf13/cobra"
)

type projectConfig struct {
	Name       string
	Framework  string
	Linter     string
	BetterAuth bool
}

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new AbacatePay project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		}
		return initializeProject(name)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializeProject(name string) error {
	cfg := projectConfig{}

	if name != "" {
		cfg.Name = name
	}

	if cfg.Name == "" {
		if err := style.Input("Project name", "my-app", &cfg.Name, scaffold.ValidateProjectName); err != nil {
			return err
		}
	}

	if err := scaffold.ValidateProjectName(cfg.Name); err != nil {
		return err
	}

	frameworkOptions := map[string]string{
		"Next.js": "next",
		"Elysia":  "elysia",
	}
	framework, err := style.Select("Which framework do you want to use?", frameworkOptions)
	if err != nil {
		return err
	}
	cfg.Framework = framework

	linterOptions := map[string]string{
		"ESLint": "eslint",
		"Biome":  "biome",
	}
	linter, err := style.Select("Which linter do you want to use?", linterOptions)
	if err != nil {
		return err
	}
	cfg.Linter = linter

	if err := style.Confirm("Do you want BetterAuth already configured?", &cfg.BetterAuth); err != nil {
		return err
	}

	betterAuthLabel := "No"
	if cfg.BetterAuth {
		betterAuthLabel = "Yes"
	}

	frameworkLabel := "Next.js"
	if cfg.Framework == "elysia" {
		frameworkLabel = "Elysia"
	}

	linterLabel := "ESLint"
	if cfg.Linter == "biome" {
		linterLabel = "Biome"
	}

	style.PrintSuccess("Project configured!", map[string]string{
		"Name":       cfg.Name,
		"Framework":  frameworkLabel,
		"Linter":     linterLabel,
		"BetterAuth": betterAuthLabel,
	})

	fmt.Printf("\n🥑 Creating project %s...\n\n", cfg.Name)

	scaffoldCfg := scaffold.Config{
		ProjectName: cfg.Name,
		Framework:   cfg.Framework,
		Linter:      cfg.Linter,
		BetterAuth:  cfg.BetterAuth,
	}

	if err := scaffold.ScaffoldProject(scaffoldCfg, "."); err != nil {
		return fmt.Errorf("failed to scaffold project: %w", err)
	}

	fmt.Println("✅ Project created successfully!")
	fmt.Println("Next steps:")

	steps := scaffoldCfg.GetNextSteps()
	for _, step := range steps {
		fmt.Printf("  %s\n", step)
	}

	return nil
}
