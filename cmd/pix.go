package cmd

import "github.com/spf13/cobra"

var pixCmd = &cobra.Command{
	Use:   "pix",
	Short: "Send PIX from your store's balance",
}

func init() {
	rootCmd.AddCommand(pixCmd)
}
