package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// PublicWebhookKey is AbacatePay's public webhook-signing key, documented at
// https://docs.abacatepay.com/pages/webhooks/security#2-assinatura-hmac.
// It is fixed and intentionally not configurable: every webhook AbacatePay
// relays is signed with this exact key, so verification code written
// against the public docs works unmodified against events forwarded by
// this CLI.
const PublicWebhookKey = "t9dXRhHHo3yDEj5pVDYz0frf7q6bMKyMRmxxCPIPp3RCplBfXRxqlC6ZpiWmOqj4L63qEaeUOtrCI8P0VMUgo6iIga2ri9ogaHFs0WIIywSMg0q7RmBfybe1E5XJcfC4IW3alNqym0tXoAKkzvfEjZxV6bE0oG2zJrNNYmUCKZyV0KZ3JS8Votf9EAWWYdiDkMkpbMdPggfh1EqHlVkMiTady6jOR3hyzGEHrIz2Ret0xHKMbiqkr9HS1JhNHDX9"

// SignWebhookPayload computes the signature AbacatePay uses for webhook
// delivery: HMAC-SHA256 over the raw body, base64-encoded.
func SignWebhookPayload(body []byte) string {
	h := hmac.New(sha256.New, []byte(PublicWebhookKey))
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
