package types

type User struct {
	Name  string
	Email string
}

type StoreResponse struct {
	Data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
	Error   string `json:"error"`
	Success any    `json:"success"`
}

type CliAuthRequestResponse struct {
	Success bool `json:"success"`
	Data    struct {
		PublicID string `json:"publicId"`
	} `json:"data"`
	Error string `json:"error"`
}

type CliAuthStatusResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	} `json:"data"`
	Error string `json:"error"`
}

type APIError struct {
	Message string `json:"error"`
}

type SimulatePaymentRequest struct {
	Metadata map[string]any `json:"metadata"`
}

type TransparentPaymentResponse struct {
	Data struct {
		ID           string         `json:"id"`
		Amount       int            `json:"amount"`
		Status       string         `json:"status"`
		DevMode      bool           `json:"devMode"`
		BRCode       string         `json:"brCode"`
		BRCodeBase64 string         `json:"brCodeBase64"`
		PlatformFee  int            `json:"platformFee"`
		ReceiptURL   string         `json:"receiptUrl"`
		CreatedAt    string         `json:"createdAt"`
		UpdatedAt    string         `json:"updatedAt"`
		ExpiresAt    string         `json:"expiresAt"`
		Metadata     map[string]any `json:"metadata"`
	} `json:"data"`
	Error   string `json:"error"`
	Success any    `json:"success"`
}
