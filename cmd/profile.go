package cmd

import (
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"profiles"},
	Short:   "Manage saved authentication profiles",
	Long:    "Allows listing, switching, and removing AbacatePay access profiles configured on this machine.",
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
