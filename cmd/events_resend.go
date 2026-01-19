package cmd

import (
	"fmt"
	"net/http"
	"time"

	"abacatepay-cli/internal/crypto"
	"abacatepay-cli/internal/logger"
	"abacatepay-cli/internal/style"
	"abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

var resendForwardURL string

var eventsResendCmd = &cobra.Command{
	Use:   "resend <event-id>",
	Short: "Resend a past event to your local webhook endpoint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return resendEvent(args[0])
	},
}

func init() {
	eventsResendCmd.Flags().StringVar(&resendForwardURL, "forward-to", "", "URL to forward the event to")
	eventsCmd.AddCommand(eventsResendCmd)
}

func resendEvent(id string) error {
	entry, err := logger.FindLogEntryByID(id)
	if err != nil {
		return err
	}

	if entry.RawMessage == "" {
		return fmt.Errorf("event %s found but has no raw payload stored", id)
	}

	defaultURL := utils.DefaultForwardURL
	if entry.URL != "" {
		defaultURL = entry.URL
	}

	url, err := utils.GetForwardURL(resendForwardURL != "", resendForwardURL, defaultURL)
	if err != nil {
		return err
	}

	deps := utils.SetupDependencies(Local, Verbose)

	secret := "whsec_abacate_local_dev_secret"
	timestamp := time.Now().Unix()
	signature := crypto.SignWebhookPayload(secret, timestamp, []byte(entry.RawMessage))

	style.LogSigningSecret(secret)
	fmt.Printf("Resending event %s to %s...\n", id, url)

	startTime := time.Now()
	resp, err := deps.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Abacate-Signature", fmt.Sprintf("t=%d,v1=%s", timestamp, signature)).
		SetBody(entry.RawMessage).
		Post(url)

	duration := time.Since(startTime)

	if err != nil {
		style.PrintError(fmt.Sprintf("Failed to forward: %v", err))
		return nil
	}

	style.LogWebhookForwarded(resp.StatusCode(), http.StatusText(resp.StatusCode()), entry.Event)

	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		style.PrintSuccess("Event resent successfully", map[string]string{
			"ID":       id,
			"Status":   fmt.Sprintf("%d %s", resp.StatusCode(), http.StatusText(resp.StatusCode())),
			"Duration": fmt.Sprintf("%dms", duration.Milliseconds()),
		})

		return nil
	}

	fmt.Printf("\nServer responded with error status: %d\n", resp.StatusCode())
	fmt.Println(string(resp.Body()))

	return nil
}
