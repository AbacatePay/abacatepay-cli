package types

import "time"

type BillingPaidEvent struct {
	ID   string `json:"id"`
	Data struct {
		Payment *EventPayment `json:"payment,omitempty"`
		Billing *EventBilling `json:"billing,omitempty"`
	} `json:"data"`
	DevMode bool   `json:"devMode"`
	Event   string `json:"event"`
}

type PayoutEvent struct {
	ID   string `json:"id"`
	Data struct {
		Transaction *EventTransaction `json:"transaction,omitempty"`
	} `json:"data"`
	DevMode bool   `json:"devMode"`
	Event   string `json:"event"`
}

type EventPayment struct {
	Amount int    `json:"amount"`
	Fee    int    `json:"fee"`
	Method string `json:"method"`
}

type EventBilling struct {
	ID         string `json:"id"`
	ExternalID string `json:"externalId"`
	URL        string `json:"url"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
}

type EventTransaction struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	DevMode     bool      `json:"devMode"`
	ReceiptURL  string    `json:"receiptUrl"`
	Kind        string    `json:"kind"`
	Amount      int       `json:"amount"`
	PlatformFee int       `json:"platformFee"`
	ExternalID  string    `json:"externalId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

