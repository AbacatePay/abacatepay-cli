package webhook

import "testing"

func TestParseEnvelope(t *testing.T) {
	tests := []struct {
		name      string
		message   string
		wantEvent string
		wantID    string
	}{
		{
			name:      "transparent v2",
			message:   `{"event":"transparent.completed","apiVersion":2,"data":{"transparent":{"id":"char_xyz789"},"customer":{"id":"cust_def456"}}}`,
			wantEvent: "transparent.completed",
			wantID:    "char_xyz789",
		},
		{
			name:      "checkout v2",
			message:   `{"event":"checkout.refunded","data":{"checkout":{"id":"bill_abc123"},"customer":{"id":"cust_abc123"}}}`,
			wantEvent: "checkout.refunded",
			wantID:    "bill_abc123",
		},
		{
			name:      "subscription prefers resource over top-level log id",
			message:   `{"id":"log_abc123","event":"subscription.renewed","data":{"subscription":{"id":"subs_123"},"payment":{"id":"char_xyz"}}}`,
			wantEvent: "subscription.renewed",
			wantID:    "subs_123",
		},
		{
			name:      "payout uses withdraw alias",
			message:   `{"event":"payout.completed","data":{"withdraw":{"id":"tran_xxx"}}}`,
			wantEvent: "payout.completed",
			wantID:    "tran_xxx",
		},
		{
			name:      "transfer v2",
			message:   `{"event":"transfer.failed","data":{"transfer":{"id":"tran_yyy"}}}`,
			wantEvent: "transfer.failed",
			wantID:    "tran_yyy",
		},
		{
			name:      "flat data.id fallback",
			message:   `{"event":"billing.paid","data":{"id":"bill_flat"}}`,
			wantEvent: "billing.paid",
			wantID:    "bill_flat",
		},
		{
			name:      "top-level id fallback",
			message:   `{"id":"evt_top","event":"something.new","data":{"other":{"id":"x"}}}`,
			wantEvent: "something.new",
			wantID:    "evt_top",
		},
		{
			name:      "no id anywhere",
			message:   `{"event":"something.new","data":{}}`,
			wantEvent: "something.new",
			wantID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, id, err := parseEnvelope([]byte(tt.message))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if event != tt.wantEvent {
				t.Errorf("event = %q, want %q", event, tt.wantEvent)
			}
			if id != tt.wantID {
				t.Errorf("id = %q, want %q", id, tt.wantID)
			}
		})
	}
}

func TestParseEnvelope_InvalidJSON(t *testing.T) {
	if _, _, err := parseEnvelope([]byte(`not json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
