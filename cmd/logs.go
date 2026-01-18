package cmd

import (
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Display recent webhook transaction logs",
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
