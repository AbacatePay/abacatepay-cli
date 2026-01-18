package cmd

import (
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View webhook transaction logs",
	Long:  "Stream live webhook events (tail) or view historical logs (list)",
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
