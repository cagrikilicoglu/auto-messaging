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

// WebhookClient defines the interface for webhook operations
type WebhookClient interface {
	SendMessage(req *model.WebhookRequest) (*model.WebhookResponse, error)
}

// webhookClient implements WebhookClient interface
type webhookClient struct {
	url     string
	authKey string
	client  *http.Client
}

// NewWebhookClient creates a new webhook client
func NewWebhookClient(url, authKey string) WebhookClient {
	return &webhookClient{
		url:     url,
		authKey: authKey,
		client:  &http.Client{},
	}
}

// SendMessage sends a message to the webhook
func (c *webhookClient) SendMessage(req *model.WebhookRequest) (*model.WebhookResponse, error) {
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

	log.Printf("Sending webhook request to %s with body: %s", c.url, string(body))
	resp, err := c.client.Do(httpReq)
	if err != nil {
		log.Printf("Webhook request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Webhook response status: %d %s", resp.StatusCode, resp.Status)

	// Accept any 2xx status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("unexpected response status: " + resp.Status)
	}

	// For 202 Accepted, we don't need to parse the response body
	if resp.StatusCode == http.StatusAccepted {
		return &model.WebhookResponse{
			Message:   "Message accepted",
			MessageID: "accepted",
		}, nil
	}

	var response model.WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Failed to decode webhook response: %v", err)
		return nil, err
	}

	log.Printf("Webhook response: %+v", response)
	return &response, nil
}
