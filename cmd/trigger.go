package cmd

import (
	"fmt"

	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "trigger a event",
	RunE: func(cmd *cobra.Command, args []string) error {
		return trigger(args)
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}

func trigger(args []string) error {
	if !utils.IsOnline() {
		return fmt.Errorf("you're offline, please stabilish your connection to continue")
	}

	if len(args) == 0 {
		return fmt.Errorf("please add a event to be triggered")
	}

	for _, evt := range args {
		if r := isEvent(evt); !r {
			return fmt.Errorf("please add a valid event to be triggered, we have: billing.paid, withdraw.done and withdraw.failed")
		}
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
