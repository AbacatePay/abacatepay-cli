package cmd

import (
	"errors"

	"abacatepay-cli/internal/mock"
	"abacatepay-cli/internal/output"
	"abacatepay-cli/internal/payments"
	"abacatepay-cli/internal/utils"

	v1 "github.com/almeidazs/go-abacate-types/v1"
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

	if err := handleEvent(deps, evt); err != nil {
		return err
	}

	return nil
}

func handleEvent(deps *utils.Dependencies, evt string) error {
	service := payments.New(deps.Client, deps.Config.APIBaseURL)

	switch evt {
	case "billing.paid":
		body := &v1.RESTPostCreateQRCodePixBody{
			Customer: &v1.APICustomerMetadata{},
		}

		body = mock.CreatePixQRCodeMock()

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
				"ID": pixID,
			},
			Data: map[string]any{
				"event": evt,
				"id":    pixID,
			},
		})

		return nil
	case "payout.done":
		return nil
	case "payout.failed":
		return nil
	default:
		return errors.New("invalid event")
	}
}
