package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func SignWebhookPayload(secret string, timestamp int64, body []byte) string {
	payload := fmt.Sprintf("%d.%s", timestamp, string(body))
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
