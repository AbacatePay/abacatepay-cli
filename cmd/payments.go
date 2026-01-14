package cmd

import "github.com/spf13/cobra"

var paymentsCmd = &cobra.Command{
	Use: "payments",
}

func init() {
	rootCmd.AddCommand(paymentsCmd)
}
