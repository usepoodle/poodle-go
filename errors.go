package poodle

import (
	"fmt"
	"net/http"
)

// PoodleError is the base interface for all Poodle SDK errors
type PoodleError interface {
	error
	StatusCode() int
	Context() map[string]interface{}
}

// BaseError provides common functionality for all error types
type BaseError struct {
	Message    string
	Code       int
	ContextMap map[string]interface{}
}

func (e *BaseError) Error() string {
	return e.Message
}

func (e *BaseError) StatusCode() int {
	return e.Code
}

func (e *BaseError) Context() map[string]interface{} {
	if e.ContextMap == nil {
		return make(map[string]interface{})
	}
	return e.ContextMap
}

// ValidationError represents validation errors (400 Bad Request)
type ValidationError struct {
	BaseError
	Errors map[string][]string
}

func NewValidationError(message string, errors map[string][]string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusBadRequest,
			ContextMap: map[string]interface{}{
				"error_type": "validation_error",
				"errors":     errors,
			},
		},
		Errors: errors,
	}
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "Validation failed"
}

// AuthenticationError represents authentication errors (401 Unauthorized)
type AuthenticationError struct {
	BaseError
}

func NewAuthenticationError(message string) *AuthenticationError {
	if message == "" {
		message = "Invalid or missing API key"
	}
	return &AuthenticationError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusUnauthorized,
			ContextMap: map[string]interface{}{
				"error_type": "authentication_error",
			},
		},
	}
}

// AccountSuspendedError represents account suspension errors (403 Forbidden)
type AccountSuspendedError struct {
	BaseError
	Reason string
}

func NewAccountSuspendedError(message, reason string) *AccountSuspendedError {
	if message == "" {
		message = "Account suspended"
	}
	return &AccountSuspendedError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusForbidden,
			ContextMap: map[string]interface{}{
				"error_type": "account_suspended",
				"reason":     reason,
			},
		},
		Reason: reason,
	}
}

// SubscriptionError represents subscription-related errors (402 Payment Required)
type SubscriptionError struct {
	BaseError
	ErrorType string
}

func NewSubscriptionError(message, errorType string) *SubscriptionError {
	if message == "" {
		message = "Subscription error"
	}
	return &SubscriptionError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusPaymentRequired,
			ContextMap: map[string]interface{}{
				"error_type":        "subscription_error",
				"subscription_type": errorType,
			},
		},
		ErrorType: errorType,
	}
}

// RateLimitError represents rate limiting errors (429 Too Many Requests)
type RateLimitError struct {
	BaseError
	RetryAfter int
	Limit      int
	Remaining  int
	Reset      int64
}

func NewRateLimitError(message string, retryAfter, limit, remaining int, reset int64) *RateLimitError {
	if message == "" {
		message = fmt.Sprintf("Rate limit exceeded. Retry after %d seconds.", retryAfter)
	}
	return &RateLimitError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusTooManyRequests,
			ContextMap: map[string]interface{}{
				"error_type":  "rate_limit_exceeded",
				"retry_after": retryAfter,
				"limit":       limit,
				"remaining":   remaining,
				"reset":       reset,
			},
		},
		RetryAfter: retryAfter,
		Limit:      limit,
		Remaining:  remaining,
		Reset:      reset,
	}
}

// NetworkError represents network connectivity errors
type NetworkError struct {
	BaseError
	URL string
}

func NewNetworkError(message, url string) *NetworkError {
	if message == "" {
		message = "Network error occurred"
	}
	return &NetworkError{
		BaseError: BaseError{
			Message: message,
			Code:    0, // No specific HTTP status for network errors
			ContextMap: map[string]interface{}{
				"error_type": "network_error",
				"url":        url,
			},
		},
		URL: url,
	}
}

func NewConnectionTimeoutError(timeout int, url string) *NetworkError {
	message := fmt.Sprintf("Connection timeout after %d seconds", timeout)
	return &NetworkError{
		BaseError: BaseError{
			Message: message,
			Code:    http.StatusRequestTimeout,
			ContextMap: map[string]interface{}{
				"error_type": "connection_timeout",
				"timeout":    timeout,
				"url":        url,
			},
		},
		URL: url,
	}
}

// HTTPError represents generic HTTP errors
type HTTPError struct {
	BaseError
	URL          string
	ResponseBody string
}

func NewHTTPError(statusCode int, message, url, responseBody string) *HTTPError {
	if message == "" {
		message = fmt.Sprintf("HTTP %d error", statusCode)
	}
	return &HTTPError{
		BaseError: BaseError{
			Message: message,
			Code:    statusCode,
			ContextMap: map[string]interface{}{
				"error_type":    "http_error",
				"url":           url,
				"response_body": responseBody,
			},
		},
		URL:          url,
		ResponseBody: responseBody,
	}
}
