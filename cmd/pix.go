package cmd

import "github.com/spf13/cobra"

var pixCmd = &cobra.Command{
	Use:   "pix",
	Short: "PIX related utilities and transactions",
}

func init() {
	rootCmd.AddCommand(pixCmd)
}
