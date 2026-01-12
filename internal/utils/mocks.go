package utils

type WebhookEvent struct {
	ID      string `json:"id"`
	Event   string `json:"event"`
	DevMode bool   `json:"devMode"`
	Data    any    `json:"data"`
}

func GetMockEvent() ([]byte, error) {
	return nil, nil
}
