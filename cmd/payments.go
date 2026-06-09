package cmd

import "github.com/spf13/cobra"

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "Payment utilities for AbacatePay dev mode",
}

func init() {
	rootCmd.AddCommand(paymentsCmd)
}
