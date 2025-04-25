package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"auto-messaging/internal/model"
)

const (
	authHeaderKey = "x-ins-auth-key"
	contentType   = "application/json"
)

// WebhookClient handles communication with webhook service
type WebhookClient struct {
	url     string
	authKey string
	client  *http.Client
}

// NewWebhookClient creates a new webhook client
func NewWebhookClient(url, authKey string) *WebhookClient {
	return &WebhookClient{
		url:     url,
		authKey: authKey,
		client:  &http.Client{},
	}
}

// SendMessage sends a message to the webhook
func (c *WebhookClient) SendMessage(req *model.WebhookRequest) (*model.WebhookResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", contentType)
	httpReq.Header.Set(authHeaderKey, c.authKey)

	log.Printf("Sending webhook request to %s", c.url)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Webhook response status: %d %s", resp.StatusCode, resp.Status)

	if resp.StatusCode != http.StatusAccepted {
		return nil, errors.New("unexpected response status: " + resp.Status)
	}

	var response model.WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
