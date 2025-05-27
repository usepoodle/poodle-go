package poodle

import (
	"encoding/json"
)

// EmailResponse represents the API response after sending an email
type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// NewEmailResponse creates a new EmailResponse
func NewEmailResponse(success bool, message string) *EmailResponse {
	return &EmailResponse{
		Success: success,
		Message: message,
	}
}

// NewEmailResponseWithError creates a new EmailResponse with an error
func NewEmailResponseWithError(success bool, message, errorMsg string) *EmailResponse {
	return &EmailResponse{
		Success: success,
		Message: message,
		Error:   errorMsg,
	}
}

// IsSuccessful returns true if the email was successfully queued
func (r *EmailResponse) IsSuccessful() bool {
	return r.Success
}

// HasError returns true if the response contains an error
func (r *EmailResponse) HasError() bool {
	return r.Error != ""
}

// ToJSON converts the response to JSON string
func (r *EmailResponse) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON creates an EmailResponse from JSON string
func FromJSON(jsonStr string) (*EmailResponse, error) {
	var response EmailResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
