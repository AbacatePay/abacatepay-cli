package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/style"

	"github.com/spf13/cobra"
)

var eventsSampleCmd = &cobra.Command{
	Use:       "sample <event>",
	Short:     "Generate a sample JSON payload for a specific event",
	Example:   "abacatepay events billing.paid\n  abacatepay events payout.done",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"billing.paid", "payout.done", "payout.failed"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return events(args[0])
	},
}

func init() {
	eventsCmd.AddCommand(eventsSampleCmd)
}

func events(evt string) error {
	var data any

	switch evt {
	case "billing.paid":
		data = mock.MockBillingPaidEvent()
	case "payout.done":
		data = mock.MockPayoutEvent(true)
	case "payout.failed":
		data = mock.MockPayoutEvent(false)
	default:
		return fmt.Errorf("unknown event type: %s. Available: billing.paid, payout.done, payout.failed", evt)
	}

	style.PrintJSON(data)

	return nil
}
