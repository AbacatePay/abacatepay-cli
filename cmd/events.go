package cmd

import (
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use: "events",
}

func init() {
	rootCmd.AddCommand(eventsCmd)
}
