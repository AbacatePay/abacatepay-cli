package cmd

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"abacatepay-cli/internal/style"

	"github.com/spf13/cobra"
)

var (
	verifySecret    string
	verifySignature string
	verifyPayload   string
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a webhook signature locally",
	Long: `Debug and verify webhook signatures offline.

This command calculates the expected signature for a given payload and secret,
then compares it with the signature header provided.

Example:
  abacatepay verify \
    --secret "whsec_..." \
    --payload '{"id":"evt_..."}' \
    --signature "t=123456,v1=abcdef..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify()
	},
}

func init() {
	verifyCmd.Flags().StringVar(&verifySecret, "secret", "", "Webhook signing secret (starts with whsec_)")
	verifyCmd.Flags().StringVar(&verifyPayload, "payload", "", "Raw JSON payload body")
	verifyCmd.Flags().StringVar(&verifySignature, "signature", "", "The value of X-Abacate-Signature header")

	_ = verifyCmd.MarkFlagRequired("secret")
	_ = verifyCmd.MarkFlagRequired("payload")
	_ = verifyCmd.MarkFlagRequired("signature")

	rootCmd.AddCommand(verifyCmd)
}

func verify() error {
	parts := strings.Split(verifySignature, ",")
	var timestampStr, receivedSig string

	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key, val := kv[0], kv[1]
		switch key {
		case "t":
			timestampStr = val
		case "v1":
			receivedSig = val
		}
	}

	if timestampStr == "" || receivedSig == "" {
		style.PrintError("Invalid signature format. Expected format: t=TIMESTAMP,v1=SIGNATURE")
		return nil
	}

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		style.PrintError("Invalid timestamp in signature header.")
		return nil
	}

	diff := time.Since(time.Unix(ts, 0))
	var timeWarning string
	if diff > 5*time.Minute {
		timeWarning = fmt.Sprintf(" (Warning: Timestamp is %s old)", diff.Round(time.Second))
	}

	signedPayload := fmt.Sprintf("%s.%s", timestampStr, verifyPayload)
	mac := hmac.New(sha256.New, []byte(verifySecret))
	mac.Write([]byte(signedPayload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	isValid := hmac.Equal([]byte(receivedSig), []byte(expectedSig))

	fields := map[string]string{
		"Timestamp": fmt.Sprintf("%d%s", ts, timeWarning),
		"Secret":    maskSecret(verifySecret),
	}

	if !isValid {
		style.PrintError("Signature mismatch ❌")

		fmt.Println("Debug Analysis:")
		fmt.Println("---------------")
		fmt.Printf("Expected: %s\n", expectedSig)
		fmt.Printf("Received: %s\n", receivedSig)
		fmt.Println("\nCommon causes for mismatch:")
		fmt.Println("1. Payload content differs (check for extra spaces, newlines, or formatting).")
		fmt.Println("2. Wrong secret key used.")
		fmt.Println("3. Timestamp manipulation.")
		return nil
	}

	fields["Status"] = "VALID ✅"
	style.PrintSuccess("Signature Verified", fields)

	return nil
}

func maskSecret(s string) string {
	if len(s) < 8 {
		return s
	}
	return s[:6] + "..." + s[len(s)-4:]
}
