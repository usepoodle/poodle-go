package poodle

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// mockHTTPClient is a mock implementation of the HTTPDoer interface for testing.
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestNewClient(t *testing.T) {
	apiKey := "test_api_key_123"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	config := client.GetConfig()
	if config.APIKey != apiKey {
		t.Errorf("Expected API key to be '%s', got '%s'", apiKey, config.APIKey)
	}

	if config.BaseURL != DefaultBaseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", DefaultBaseURL, config.BaseURL)
	}
}

func TestNewClientFromEnv(t *testing.T) {
	// Set environment variables
	apiKey := "env_test_api_key_123"
	baseURL := "https://test.api.usepoodle.com"

	os.Setenv("POODLE_API_KEY", apiKey)
	os.Setenv("POODLE_BASE_URL", baseURL)
	os.Setenv("POODLE_DEBUG", "true")
	defer func() {
		os.Unsetenv("POODLE_API_KEY")
		os.Unsetenv("POODLE_BASE_URL")
		os.Unsetenv("POODLE_DEBUG")
	}()

	client := NewClientFromEnv()

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	config := client.GetConfig()
	if config.APIKey != apiKey {
		t.Errorf("Expected API key to be '%s', got '%s'", apiKey, config.APIKey)
	}

	if config.BaseURL != baseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", baseURL, config.BaseURL)
	}

	if !config.Debug {
		t.Error("Expected debug to be true")
	}
}

func TestNewClientWithConfig(t *testing.T) {
	config := &Config{
		APIKey:         "custom_api_key",
		BaseURL:        "https://custom.api.com",
		Timeout:        45 * time.Second,
		ConnectTimeout: 15 * time.Second,
		Debug:          true,
	}

	client := NewClientWithConfig(config)

	if client == nil {
		t.Fatal("Expected client to be created, got nil")
	}

	clientConfig := client.GetConfig()
	if clientConfig.APIKey != config.APIKey {
		t.Errorf("Expected API key to be '%s', got '%s'", config.APIKey, clientConfig.APIKey)
	}

	if clientConfig.BaseURL != config.BaseURL {
		t.Errorf("Expected base URL to be '%s', got '%s'", config.BaseURL, clientConfig.BaseURL)
	}

	if clientConfig.Timeout != config.Timeout {
		t.Errorf("Expected timeout to be %v, got %v", config.Timeout, clientConfig.Timeout)
	}
}

func TestNewClientWithInvalidConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid config, but didn't get one")
		}
	}()

	config := &Config{
		// Missing API key - should cause panic
		BaseURL: "https://api.usepoodle.com",
	}

	NewClientWithConfig(config)
}

func TestClientDebugMethods(t *testing.T) {
	client := NewClient("test_api_key")

	// Test initial debug state
	if client.IsDebug() {
		t.Error("Expected debug to be false initially")
	}

	// Test SetDebug
	client.SetDebug(true)
	if !client.IsDebug() {
		t.Error("Expected debug to be true after SetDebug(true)")
	}

	client.SetDebug(false)
	if client.IsDebug() {
		t.Error("Expected debug to be false after SetDebug(false)")
	}
}

func TestClientGetConfig(t *testing.T) {
	client := NewClient("test_api_key")

	config1 := client.GetConfig()
	config2 := client.GetConfig()

	// Should return copies, not the same instance
	if config1 == config2 {
		t.Error("GetConfig should return copies, not the same instance")
	}

	// But should have the same values
	if config1.APIKey != config2.APIKey {
		t.Error("Config copies should have the same API key")
	}
}

func TestClientConcurrency(t *testing.T) {
	client := NewClient("test_api_key")

	// Test concurrent access to debug methods
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			client.SetDebug(true)
			client.IsDebug()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			client.SetDebug(false)
			client.GetConfig()
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// If we get here without a race condition, the test passes
}

// Note: We can't easily test the actual Send methods without mocking the HTTP client
// or setting up integration tests. For now, we'll test that the methods exist and
// can be called with valid parameters.

func TestClientSendMethods(t *testing.T) {
	// Common setup for tests
	from := "from@example.com"
	to := "to@example.com"
	subject := "Test Subject"
	htmlBody := "<h1>Hello</h1>"
	textBody := "Hello"

	tests := []struct {
		name          string
		mockResponse  *http.Response
		mockErr       error
		clientSetup   func(c *Client, mock *mockHTTPClient)
		sendAction    func(c *Client) (*EmailResponse, error)
		expectSuccess bool
		expectError   bool
		errorType     interface{}
	}{
		{
			name: "Send - Success",
			mockResponse: &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(strings.NewReader(`{"success": true, "message": "Email queued"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectSuccess: true,
		},
		{
			name: "SendHTML - Success",
			mockResponse: &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(strings.NewReader(`{"success": true, "message": "Email queued"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				return c.SendHTML(from, to, subject, htmlBody)
			},
			expectSuccess: true,
		},
		{
			name: "SendText - Success",
			mockResponse: &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(strings.NewReader(`{"success": true, "message": "Email queued"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				return c.SendText(from, to, subject, textBody)
			},
			expectSuccess: true,
		},
		{
			name: "SendWithBoth - Success",
			mockResponse: &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(strings.NewReader(`{"success": true, "message": "Email queued"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				return c.SendWithBoth(from, to, subject, htmlBody, textBody)
			},
			expectSuccess: true,
		},
		{
			name: "Send - Validation Error (from email model)",
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail("invalid", to, subject, htmlBody) // Invalid from
				return c.Send(email)
			},
			expectError: true,
			errorType:  &ValidationError{},
		},
		{
			name: "Send - API Validation Error",
			mockResponse: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"success": false, "message": "Invalid request", "error": "validation details"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectError: true,
			errorType:  &ValidationError{},
		},
		{
			name: "Send - Authentication Error",
			mockResponse: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"message": "Invalid API Key"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectError: true,
			errorType:  &AuthenticationError{},
		},
		{
			name: "Send - Rate Limit Error",
			mockResponse: &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(strings.NewReader(`{"message": "Rate limit exceeded"}`)),
				Header: http.Header{
					"Retry-After":         {"60"},
					"Ratelimit-Limit":     {"100"},
					"Ratelimit-Remaining": {"0"},
					"Ratelimit-Reset":     {"1678886400"},
				},
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectError: true,
			errorType:  &RateLimitError{},
		},
		{
			name: "Send - Network Error (simulated by mockErr)",
			mockErr: NewNetworkError("simulated network problem", "http://fakeurl.com"),
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectError: true,
			errorType:  &NetworkError{},
		},
		{
			name: "Send - HTTP Error (generic)",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{"message": "Internal Server Error"}`)),
			},
			sendAction: func(c *Client) (*EmailResponse, error) {
				email := NewHTMLEmail(from, to, subject, htmlBody)
				return c.Send(email)
			},
			expectError: true,
			errorType:  &HTTPError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("test_api_key")
			mock := &mockHTTPClient{
				response: tt.mockResponse,
				err:      tt.mockErr,
			}
			client.httpClient.httpClient = mock

			if tt.clientSetup != nil {
				tt.clientSetup(client, mock)
			}

			resp, err := tt.sendAction(client)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tt.errorType != nil {
					switch tt.errorType.(type) {
					case *ValidationError:
						if _, ok := err.(*ValidationError); !ok {
							t.Errorf("Expected ValidationError, got %T", err)
						}
					case *AuthenticationError:
						if _, ok := err.(*AuthenticationError); !ok {
							t.Errorf("Expected AuthenticationError, got %T", err)
						}
					case *RateLimitError:
						if _, ok := err.(*RateLimitError); !ok {
							t.Errorf("Expected RateLimitError, got %T", err)
						}
					case *NetworkError:
						if _, ok := err.(*NetworkError); !ok {
							t.Errorf("Expected NetworkError, got %T", err)
						}
					case *HTTPError:
						if _, ok := err.(*HTTPError); !ok {
							t.Errorf("Expected HTTPError, got %T", err)
						}
					default:
						t.Errorf("Unhandled expected error type: %T", tt.errorType)
					}
				}
			} else if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if tt.expectSuccess {
				if resp == nil {
					t.Errorf("Expected response, got nil")
					return
				}
				if !resp.Success {
					t.Errorf("Expected response to be successful, but was not. Message: %s, Error: %s", resp.Message, resp.Error)
				}
			} else if resp != nil && resp.Success {
				t.Errorf("Expected response to not be successful, but it was.")
			}
		})
	}
}
