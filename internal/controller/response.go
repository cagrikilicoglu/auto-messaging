package controller

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a success response for messaging operations
type MessageResponse struct {
	Message string `json:"message"`
}
