package cmd

import (
	"fmt"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger <event>",
	Args:  cobra.ExactArgs(1),
	Short: "Trigger test events",
	RunE: func(cmd *cobra.Command, args []string) error {
		return trigger(args[0])
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}

func trigger(evt string) error {
	if !utils.IsOnline() {
		return fmt.Errorf("you’re offline — check your connection and try again")
	}

	if r := isEvent(evt); !r {
		return fmt.Errorf("unknown event '%s'. Available events: billing.paid, withdraw.done, withdraw.failed", evt)
	}

	return nil
}

func isEvent(evt string) bool {
	switch evt {
	case "billing.paid":
		return true
	case "withdraw.done":
		return true
	case "withdraw.failed":
		return true
	default:
		return false
	}
}
