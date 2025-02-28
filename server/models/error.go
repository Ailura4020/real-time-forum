package models

// ErrorResponse represents the structure of the error response sent to the client.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
