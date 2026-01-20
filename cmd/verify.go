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

type signatureParts struct {
	timestamp int64
	signature string
}

func verify() error {
	parts, err := parseSignatureHeader(verifySignature)
	if err != nil {
		style.PrintError(err.Error())
		return err
	}

	expectedSig := computeSignature(verifySecret, parts.timestamp, verifyPayload)
	isValid := hmac.Equal([]byte(parts.signature), []byte(expectedSig))

	if !isValid {
		style.PrintVerifyError(expectedSig, parts.signature)
		return fmt.Errorf("signature mismatch")
	}

	printVerifySuccess(parts.timestamp, verifySecret)
	return nil
}

func parseSignatureHeader(header string) (*signatureParts, error) {
	var timestampStr, signature string

	for part := range strings.SplitSeq(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "t":
			timestampStr = kv[1]
		case "v1":
			signature = kv[1]
		}
	}

	if timestampStr == "" || signature == "" {
		return nil, fmt.Errorf("invalid signature format. Expected: t=TIMESTAMP,v1=SIGNATURE")
	}

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp in signature header")
	}

	return &signatureParts{
		timestamp: ts,
		signature: signature,
	}, nil
}

func computeSignature(secret string, timestamp int64, payload string) string {
	signedPayload := fmt.Sprintf("%d.%s", timestamp, payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func printVerifySuccess(timestamp int64, secret string) {
	fields := map[string]string{
		"Timestamp": formatTimestamp(timestamp),
		"Secret":    maskSecret(secret),
		"Status":    "VALID",
	}
	style.PrintSuccess("Signature Verified", fields)
}

func formatTimestamp(ts int64) string {
	diff := time.Since(time.Unix(ts, 0))
	if diff > 5*time.Minute {
		return fmt.Sprintf("%d (Warning: %s old)", ts, diff.Round(time.Second))
	}
	return fmt.Sprintf("%d", ts)
}

func maskSecret(s string) string {
	if len(s) < 8 {
		return s
	}
	return s[:6] + "..." + s[len(s)-4:]
}
