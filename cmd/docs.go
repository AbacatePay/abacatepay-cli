package cmd

import (
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Open the CLI documentation in your browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		return utils.OpenBrowser("https://docs.abacatepay.com/pages/cli")
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
