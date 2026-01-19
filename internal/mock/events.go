package mock

import (
	"fmt"
	"time"

	"abacatepay-cli/internal/types"

	"github.com/brianvoe/gofakeit/v7"
)

func MockBillingPaidEvent() *types.BillingPaidEvent {
	amount := gofakeit.Number(100, 1000)
	id := fmt.Sprintf("evt_%s", gofakeit.LetterN(10))

	return &types.BillingPaidEvent{
		ID: id,
		Data: struct {
			Payment *types.EventPayment `json:"payment,omitempty"`
			Billing *types.EventBilling `json:"billing,omitempty"`
		}{
			Payment: &types.EventPayment{
				Amount: amount,
				Fee:    gofakeit.Number(10, 100),
				Method: "PIX",
			},
			Billing: &types.EventBilling{
				ID:         fmt.Sprintf("bill_%s", gofakeit.LetterN(10)),
				ExternalID: gofakeit.UUID(),
				Amount:     amount,
				URL:        "https://docs.abacatepay.com/pages/webhooks#billing-paid",
				Status:     "PAID",
			},
		},
		DevMode: true,
		Event:   "billing.paid",
	}
}

func MockPayoutEvent(isDone bool) *types.PayoutEvent {
	status := "CANCELLED"
	event := "payout.failed"
	if isDone {
		status = "COMPLETE"
		event = "payout.done"
	}

	amount := gofakeit.Number(1000, 50000)
	id := fmt.Sprintf("evt_%s", gofakeit.LetterN(10))

	return &types.PayoutEvent{
		ID: id,
		Data: struct {
			Transaction *types.EventTransaction `json:"transaction,omitempty"`
		}{
			Transaction: &types.EventTransaction{
				ID:          fmt.Sprintf("tran_%s", gofakeit.LetterN(16)),
				Status:      status,
				DevMode:     true,
				ReceiptURL:  "https://abacatepay.com/receipt/mock",
				Kind:        "WITHDRAW",
				Amount:      amount,
				PlatformFee: 0,
				ExternalID:  gofakeit.UUID(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		DevMode: true,
		Event:   event,
	}
}
