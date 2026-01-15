package types

type Customer struct {
	Name      string `json:"name,omitempty"`
	Cellphone string `json:"cellphone,omitempty"`
	Email     string `json:"email,omitempty"`
	TaxID     string `json:"taxId,omitempty"`
}

type CreateCheckoutRequest struct {
	Items         []Item    `json:"items"`
	Method        string    `json:"method,omitempty"` // PIX | CARD (default: PIX)
	ReturnURL     string    `json:"returnUrl,omitempty"`
	CompletionURL string    `json:"completionUrl,omitempty"`
	CustomerID    string    `json:"customerId,omitempty"`
	Customer      *Customer `json:"customer,omitempty"`
	Coupons       []string  `json:"coupons,omitempty"`
	ExternalID    string    `json:"externalId,omitempty"`
}

type Item struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type CheckoutResponse struct {
	Data struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Status string `json:"status"`
		Amount int    `json:"amount"`
	} `json:"data"`
}

type PixResponse struct {
	Data struct {
		ID     string `json:"id"`
		BRCode string `json:"brCode"`
		Status string `json:"status"`
	} `json:"data"`
}

type User struct {
	Name  string
	Email string
}

type DeviceLoginResponse struct {
	DeviceCode      string `json:"deviceCode"`
	VerificationURI string `json:"verificationUri"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type APIError struct {
	Message string `json:"error"`
}
