package cmd

import (
	"fmt"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/payments"
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
	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	return handleEvent(deps, evt)
}

func handleEvent(deps *utils.Dependencies, evt string) error {
	service := payments.New(deps.Client, deps.Config.APIBaseURL)

	switch evt {
	case "billing.paid":
		body := mock.CreatePixQRCodeMock()

		pixID, err := service.CreatePixQRCode(body, true)
		if err != nil {
			return err
		}

		if err := service.SimulatePixQRCodePayment(pixID, true); err != nil {
			return err
		}

		output.Print(output.Result{
			Title: "Billing.paid Triggered",
			Fields: map[string]string{
				"Charge ID": pixID,
				"Status":    "Simulated",
				"Note":      "Check your 'listen' terminal for the webhook event",
			},
			Data: map[string]any{
				"event":    evt,
				"chargeId": pixID,
			},
		})

		return nil

	case "payout.done", "payout.failed":
		isDone := evt == "payout.done"
		mockEvent := mock.MockPayoutEvent(isDone)

		output.Print(output.Result{
			Title: fmt.Sprintf("Mock %s Triggered", evt),
			Fields: map[string]string{
				"Event ID":       mockEvent.ID,
				"Transaction ID": mockEvent.Data.Transaction.ID,
				"Status":         mockEvent.Data.Transaction.Status,
			},
			Data: mockEvent,
		})

		fmt.Printf("\nTip: Use 'abacatepay events resend %s' to send this mock to your local server.\n", mockEvent.ID)

		// Optional: We could automatically append this mock to the local log file
		// so it appears in 'logs list' immediately.
		return nil

	default:
		return fmt.Errorf("invalid event type: %s. Supported: billing.paid, payout.done, payout.failed", evt)
	}
}
