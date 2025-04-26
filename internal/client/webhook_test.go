package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"auto-messaging/internal/model"
)

func TestWebhookClient_SendMessage(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		req           *model.WebhookRequest
		expectedError bool
	}{
		{
			name: "successful message send",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message_id": "test-123", "message": "Message sent successfully"}`))
			},
			req: &model.WebhookRequest{
				Content: "Test message",
				To:      "test@example.com",
			},
			expectedError: false,
		},
		{
			name: "webhook returns 202 accepted",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			},
			req: &model.WebhookRequest{
				Content: "Test message",
				To:      "test@example.com",
			},
			expectedError: false,
		},
		{
			name: "webhook returns error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			req: &model.WebhookRequest{
				Content: "Test message",
				To:      "test@example.com",
			},
			expectedError: true,
		},
		{
			name: "invalid response format",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`invalid json`))
			},
			req: &model.WebhookRequest{
				Content: "Test message",
				To:      "test@example.com",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			// Create webhook client with test server URL
			client := NewWebhookClient(server.URL, "test-auth-key")

			// Send message
			resp, err := client.SendMessage(tt.req)

			// Check error
			if (err != nil) != tt.expectedError {
				t.Errorf("SendMessage() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			// If no error expected, check response
			if !tt.expectedError {
				if resp == nil {
					t.Error("Expected non-nil response")
					return
				}

				// For 202 Accepted, we expect a default response
				if tt.handler == nil {
					if resp.MessageID != "accepted" {
						t.Errorf("Expected MessageID 'accepted', got %s", resp.MessageID)
					}
				}
			}
		})
	}
}

func TestWebhookClient_Headers(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("x-ins-auth-key") != "test-key" {
			t.Errorf("Expected x-ins-auth-key: test-key, got %s", r.Header.Get("x-ins-auth-key"))
		}

		// Read and verify request body
		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if reqBody["to"] != "test@example.com" || reqBody["content"] != "Test message" {
			t.Errorf("Unexpected request body: %v", reqBody)
		}

		// Send a valid response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Message sent successfully",
		})
	}))
	defer server.Close()

	// Create webhook client
	client := NewWebhookClient(server.URL, "test-key")

	// Create test request
	request := &model.WebhookRequest{
		Content: "Test message",
		To:      "test@example.com",
	}

	// Send message
	response, err := client.SendMessage(request)
	if err != nil {
		t.Errorf("SendMessage() error = %v", err)
	}
	if response.Message != "Message sent successfully" {
		t.Errorf("Expected response message 'Message sent successfully', got '%s'", response.Message)
	}
}
