package model

// WebhookRequest represents the payload sent to webhook
type WebhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// WebhookResponse represents the response from webhook
type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}
