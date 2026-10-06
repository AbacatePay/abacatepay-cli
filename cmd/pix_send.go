package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AbacatePay/abacatepay-cli/internal/output"
	"github.com/AbacatePay/abacatepay-cli/internal/style"
	"github.com/AbacatePay/abacatepay-cli/internal/tui"
	"github.com/AbacatePay/abacatepay-cli/internal/utils"

	"github.com/spf13/cobra"
)

const pixSendMinAmount = 100

var pixKeyTypes = []string{"CPF", "CNPJ", "EMAIL", "PHONE", "RANDOM"}

var (
	pixSendAmount      int
	pixSendKeyType     string
	pixSendKey         string
	pixSendExternalID  string
	pixSendOTP         string
	pixSendDescription string
	pixSendYes         bool
)

var pixSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a PIX from your store's balance to any key",
	Long: "Sends a PIX from the active profile's store balance. When the store is in production this moves real money.\n\n" +
		"It requires the store owner's two-factor code (--otp, from the authenticator app set up in the dashboard) " +
		"and a unique --external-id: sending the same --external-id twice is refused.",
	Example: "  abacatepay pix send --amount 1000 --pix-key-type EMAIL --pix-key someone@example.com --external-id order-123 --otp 123456",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pixSend()
	},
}

func init() {
	pixSendCmd.Flags().IntVar(&pixSendAmount, "amount", 0, "Amount in cents (minimum 100)")
	pixSendCmd.Flags().StringVar(&pixSendKeyType, "pix-key-type", "", "PIX key type: "+strings.Join(pixKeyTypes, ", "))
	pixSendCmd.Flags().StringVar(&pixSendKey, "pix-key", "", "Destination PIX key")
	pixSendCmd.Flags().StringVar(&pixSendExternalID, "external-id", "", "Your unique ID for this transfer")
	pixSendCmd.Flags().StringVar(&pixSendOTP, "otp", "", "Two-factor code from the store owner's authenticator app")
	pixSendCmd.Flags().StringVar(&pixSendDescription, "description", "", "Optional description shown to the recipient")
	pixSendCmd.Flags().BoolVarP(&pixSendYes, "yes", "y", false, "Skip the confirmation prompt (required when not running in a terminal)")

	for _, flag := range []string{"amount", "pix-key-type", "pix-key", "external-id", "otp"} {
		_ = pixSendCmd.MarkFlagRequired(flag)
	}

	pixCmd.AddCommand(pixSendCmd)
}

type pixSendPix struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type pixSendRequest struct {
	ExternalID  string     `json:"externalId"`
	Amount      int        `json:"amount"`
	Pix         pixSendPix `json:"pix"`
	Description string     `json:"description,omitempty"`
	OTP         string     `json:"otp"`
}

type pixSendResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    *struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		DevMode     bool   `json:"devMode"`
		ReceiptURL  string `json:"receiptUrl"`
		Amount      int    `json:"amount"`
		PlatformFee int    `json:"platformFee"`
		ExternalID  string `json:"externalId"`
	} `json:"data"`
}

func pixSend() error {
	keyType := strings.ToUpper(strings.TrimSpace(pixSendKeyType))
	key := strings.TrimSpace(pixSendKey)
	externalID := strings.TrimSpace(pixSendExternalID)

	if pixSendAmount < pixSendMinAmount {
		return fmt.Errorf("--amount must be at least %d cents (%s)", pixSendMinAmount, formatCents(pixSendMinAmount))
	}
	if !isPixKeyType(keyType) {
		return fmt.Errorf("--pix-key-type must be one of %s", strings.Join(pixKeyTypes, ", "))
	}
	if key == "" {
		return errors.New("--pix-key cannot be empty")
	}
	if externalID == "" {
		return errors.New("--external-id cannot be empty")
	}
	if strings.TrimSpace(pixSendOTP) == "" {
		return errors.New("--otp cannot be empty")
	}

	if !pixSendYes {
		if !tui.IsInteractive() {
			return errors.New("refusing to send without confirmation: pass --yes when not running in a terminal")
		}

		confirmed := false
		prompt := fmt.Sprintf("Send %s to the %s key %s?", formatCents(pixSendAmount), keyType, key)
		if err := style.Confirm(prompt, &confirmed); err != nil {
			return err
		}
		if !confirmed {
			style.PrintInfo("Cancelled, nothing was sent")
			return nil
		}
	}

	deps, err := utils.SetupClient(Local, Verbose)
	if err != nil {
		return err
	}

	request := pixSendRequest{
		ExternalID:  externalID,
		Amount:      pixSendAmount,
		Pix:         pixSendPix{Type: keyType, Key: key},
		Description: strings.TrimSpace(pixSendDescription),
		OTP:         strings.TrimSpace(pixSendOTP),
	}

	return output.RunTask("Sending PIX...", func() (output.Result, error) {
		var response pixSendResponse

		resp, err := deps.Client.R().
			SetBody(request).
			SetResult(&response).
			SetError(&response).
			Post(deps.Config.APIBaseURL + "/cli/pix/send")
		if err != nil {
			return output.Result{}, fmt.Errorf("failed to send PIX: %w", err)
		}

		if !response.Success || response.Data == nil {
			if response.Error != "" {
				return output.Result{}, errors.New(response.Error)
			}
			return output.Result{}, fmt.Errorf("failed to send PIX: unexpected response (status %d)", resp.StatusCode())
		}

		sent := response.Data
		fields := map[string]string{
			"ID":          sent.ID,
			"Status":      sent.Status,
			"Amount":      formatCents(sent.Amount),
			"Fee":         formatCents(sent.PlatformFee),
			"External ID": sent.ExternalID,
			"Receipt":     sent.ReceiptURL,
		}
		if sent.DevMode {
			fields["Mode"] = "dev"
		}

		return output.Result{Title: "PIX submitted", Fields: fields}, nil
	})
}

func isPixKeyType(keyType string) bool {
	for _, valid := range pixKeyTypes {
		if keyType == valid {
			return true
		}
	}
	return false
}

func formatCents(cents int) string {
	return fmt.Sprintf("R$ %d,%02d", cents/100, cents%100)
}
