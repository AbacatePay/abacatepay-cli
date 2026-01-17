package cmd

import (
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logs()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func logs() error {
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}
	return nil
}
