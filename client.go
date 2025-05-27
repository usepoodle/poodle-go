package poodle

import (
	"sync"
)

// Client is the main Poodle SDK client
type Client struct {
	config     *Config
	httpClient *HTTPClient
	mutex      sync.RWMutex
}

// NewClient creates a new Poodle client with the provided API key
func NewClient(apiKey string) *Client {
	config := NewConfig()
	config.APIKey = apiKey
	return NewClientWithConfig(config)
}

// NewClientFromEnv creates a new Poodle client using environment variables
func NewClientFromEnv() *Client {
	config := NewConfigFromEnv()
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates a new Poodle client with custom configuration
func NewClientWithConfig(config *Config) *Client {
	if err := config.Validate(); err != nil {
		panic(err) // In Go 1.20, we don't have better error handling for constructors
	}

	return &Client{
		config:     config,
		httpClient: NewHTTPClient(config),
	}
}

// Send sends an email using the Email model
func (c *Client) Send(email *Email) (*EmailResponse, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.httpClient.SendEmail(email)
}

// SendHTML sends an HTML email
func (c *Client) SendHTML(from, to, subject, html string) (*EmailResponse, error) {
	email := NewHTMLEmail(from, to, subject, html)
	return c.Send(email)
}

// SendText sends a plain text email
func (c *Client) SendText(from, to, subject, text string) (*EmailResponse, error) {
	email := NewTextEmail(from, to, subject, text)
	return c.Send(email)
}

// SendWithBoth sends an email with both HTML and text content
func (c *Client) SendWithBoth(from, to, subject, html, text string) (*EmailResponse, error) {
	email := NewEmailWithBoth(from, to, subject, html, text)
	return c.Send(email)
}

// GetConfig returns the client configuration (read-only)
func (c *Client) GetConfig() *Config {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *c.config
	return &configCopy
}

// SetDebug enables or disables debug logging
func (c *Client) SetDebug(debug bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.config.Debug = debug
}

// IsDebug returns whether debug logging is enabled
func (c *Client) IsDebug() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.config.Debug
}
