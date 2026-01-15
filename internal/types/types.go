package types

type Customer struct {
	Name      string `json:"name,omitempty"`
	Cellphone string `json:"cellphone,omitempty"`
	Email     string `json:"email,omitempty"`
	TaxID     string `json:"taxId,omitempty"`
}

type CheckoutResponse struct {
	Data struct {
		ID     string `json:"id"`
		URL    string `json:"url"`
		Status string `json:"status"`
		Amount int    `json:"amount"`
	} `json:"data"`
}
