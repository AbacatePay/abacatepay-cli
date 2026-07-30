package crypto

import "testing"

func TestSignWebhookPayload_MatchesKnownVector(t *testing.T) {
	body := []byte(`{"event":"billing.paid","data":{"id":"test"}}`)
	got := SignWebhookPayload(body)
	want := "+2VQ4f40+Ffmkpzlb4WFb0rcp3PDG2kSUWdEmWOVzDU="
	if got != want {
		t.Fatalf("signature mismatch: got %q, want %q", got, want)
	}
}
