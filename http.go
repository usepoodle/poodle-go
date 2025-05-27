package poodle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPDoer is an interface for making HTTP requests.
// It is implemented by *http.Client.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClient handles HTTP communication with the Poodle API
type HTTPClient struct {
	config     *Config
	httpClient HTTPDoer // Changed from *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(config *Config) *HTTPClient {
	// Create a custom dialer for connection timeout
	dialer := &net.Dialer{
		Timeout:   config.ConnectTimeout, // This is the connection timeout
		KeepAlive: 30 * time.Second,      // Default keep-alive, can be configured if needed
	}

	transport := &http.Transport{
		// DialContext is preferred, but Dial is used for Go 1.20 compatibility.
		// The timeout is handled by the net.Dialer.
		Dial: func(network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		MaxIdleConns:          100, // Default, can be configured
		IdleConnTimeout:       90 * time.Second, // Default, can be configured
		TLSHandshakeTimeout:   10 * time.Second, // Default, can be configured
		ExpectContinueTimeout: 1 * time.Second, // Default, can be configured
	}

	return &HTTPClient{
		config: config,
		httpClient: &http.Client{
			Timeout:   config.Timeout, // This is the total request timeout
			Transport: transport,
		},
	}
}

// SendEmail sends an email via the API
func (c *HTTPClient) SendEmail(email *Email) (*EmailResponse, error) {
	// Validate email before sending
	if err := email.Validate(); err != nil {
		return nil, err
	}

	// Prepare request body
	requestBody, err := json.Marshal(email)
	if err != nil {
		return nil, NewNetworkError("Failed to encode request body", "")
	}

	// Build URL
	url := strings.TrimRight(c.config.BaseURL, "/") + "/v1/send-email"

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, NewNetworkError("Failed to create request", url)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("User-Agent", c.config.GetUserAgent())

	// Debug logging
	if c.config.Debug {
		log.Printf("Poodle API Request: %s %s", req.Method, req.URL.String())
		log.Printf("Request Body: %s", string(requestBody))
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Handle timeout errors
		if strings.Contains(err.Error(), "timeout") {
			timeout := int(c.config.Timeout.Seconds())
			return nil, NewConnectionTimeoutError(timeout, url)
		}
		return nil, NewNetworkError("Request failed: "+err.Error(), url)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewNetworkError("Failed to read response body", url)
	}

	// Debug logging
	if c.config.Debug {
		log.Printf("Poodle API Response: %d %s", resp.StatusCode, string(responseBody))
	}

	// Handle different status codes
	switch resp.StatusCode {
	case http.StatusAccepted: // 202 - Success
		return c.parseSuccessResponse(responseBody)

	case http.StatusBadRequest: // 400 - Validation error
		return nil, c.parseValidationError(responseBody)

	case http.StatusUnauthorized: // 401 - Authentication error
		return nil, c.parseAuthenticationError(responseBody)

	case http.StatusPaymentRequired: // 402 - Subscription error
		return nil, c.parseSubscriptionError(responseBody)

	case http.StatusForbidden: // 403 - Account suspended
		return nil, c.parseAccountSuspendedError(responseBody)

	case http.StatusUnprocessableEntity: // 422 - Job queue error
		return nil, c.parseValidationError(responseBody)

	case http.StatusTooManyRequests: // 429 - Rate limit
		return nil, c.parseRateLimitError(resp, responseBody)

	default:
		// Generic HTTP error
		return nil, c.parseGenericError(resp.StatusCode, responseBody, url)
	}
}

// parseSuccessResponse parses a successful API response
func (c *HTTPClient) parseSuccessResponse(body []byte) (*EmailResponse, error) {
	var response EmailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, NewNetworkError("Failed to parse response", "")
	}
	return &response, nil
}

// parseValidationError parses validation error responses
func (c *HTTPClient) parseValidationError(body []byte) error {
	var apiResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return NewValidationError("Validation failed", map[string][]string{
			"request": {"Invalid request format"},
		})
	}

	// Create a simple validation error
	errors := map[string][]string{
		"request": {apiResponse.Message},
	}

	if apiResponse.Error != "" {
		errors["details"] = []string{apiResponse.Error}
	}

	return NewValidationError(apiResponse.Message, errors)
}

// parseAuthenticationError parses authentication error responses
func (c *HTTPClient) parseAuthenticationError(body []byte) error {
	var apiResponse struct {
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return NewAuthenticationError("Invalid or missing API key")
	}

	return NewAuthenticationError(apiResponse.Message)
}

// parseSubscriptionError parses subscription error responses
func (c *HTTPClient) parseSubscriptionError(body []byte) error {
	var apiResponse struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return NewSubscriptionError("Subscription error", "unknown")
	}

	// Determine subscription error type from message
	errorType := "unknown"
	message := apiResponse.Message
	if strings.Contains(message, "expired") {
		errorType = "subscription_expired"
	} else if strings.Contains(message, "trial") {
		errorType = "trial_limit_reached"
	} else if strings.Contains(message, "limit") {
		errorType = "limit_reached"
	}

	return NewSubscriptionError(message, errorType)
}

// parseAccountSuspendedError parses account suspended error responses
func (c *HTTPClient) parseAccountSuspendedError(body []byte) error {
	var apiResponse struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return NewAccountSuspendedError("Account suspended", "unknown")
	}

	return NewAccountSuspendedError(apiResponse.Message, apiResponse.Error)
}

// parseRateLimitError parses rate limit error responses
func (c *HTTPClient) parseRateLimitError(resp *http.Response, body []byte) error {
	var apiResponse struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	// Parse response body
	json.Unmarshal(body, &apiResponse)

	// Extract rate limit information from headers
	retryAfter := 0
	if retryAfterStr := resp.Header.Get("retry-after"); retryAfterStr != "" {
		if val, err := strconv.Atoi(retryAfterStr); err == nil {
			retryAfter = val
		}
	}

	limit := 0
	if limitStr := resp.Header.Get("ratelimit-limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			limit = val
		}
	}

	remaining := 0
	if remainingStr := resp.Header.Get("ratelimit-remaining"); remainingStr != "" {
		if val, err := strconv.Atoi(remainingStr); err == nil {
			remaining = val
		}
	}

	reset := int64(0)
	if resetStr := resp.Header.Get("ratelimit-reset"); resetStr != "" {
		if val, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			reset = val
		}
	}

	message := apiResponse.Message
	if message == "" {
		message = fmt.Sprintf("Rate limit exceeded. Retry after %d seconds.", retryAfter)
	}

	return NewRateLimitError(message, retryAfter, limit, remaining, reset)
}

// parseGenericError parses generic HTTP error responses
func (c *HTTPClient) parseGenericError(statusCode int, body []byte, url string) error {
	var apiResponse struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	message := fmt.Sprintf("HTTP %d error", statusCode)
	if err := json.Unmarshal(body, &apiResponse); err == nil && apiResponse.Message != "" {
		message = apiResponse.Message
	}

	return NewHTTPError(statusCode, message, url, string(body))
}
