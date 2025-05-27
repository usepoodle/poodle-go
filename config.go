package poodle

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Default configuration values
const (
	DefaultBaseURL        = "https://api.usepoodle.com"
	DefaultTimeout        = 30 * time.Second
	DefaultConnectTimeout = 10 * time.Second
	SDKVersion            = "1.0.0"
)

// Config holds the configuration for the Poodle client
type Config struct {
	APIKey         string
	BaseURL        string
	Timeout        time.Duration
	ConnectTimeout time.Duration
	Debug          bool
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		BaseURL:        DefaultBaseURL,
		Timeout:        DefaultTimeout,
		ConnectTimeout: DefaultConnectTimeout,
		Debug:          false,
	}
}

// NewConfigFromEnv creates a new configuration from environment variables
func NewConfigFromEnv() *Config {
	config := NewConfig()

	if apiKey := os.Getenv("POODLE_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	if baseURL := os.Getenv("POODLE_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if timeoutStr := os.Getenv("POODLE_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.Timeout = timeout
		}
	}

	if connectTimeoutStr := os.Getenv("POODLE_CONNECT_TIMEOUT"); connectTimeoutStr != "" {
		if connectTimeout, err := time.ParseDuration(connectTimeoutStr); err == nil {
			config.ConnectTimeout = connectTimeout
		}
	}

	if debugStr := os.Getenv("POODLE_DEBUG"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			config.Debug = debug
		}
	}

	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return &ValidationError{
			BaseError: BaseError{Message: "API key is required"},
			Errors: map[string][]string{
				"api_key": {"API key is required"},
			},
		}
	}

	if c.BaseURL == "" {
		return &ValidationError{
			BaseError: BaseError{Message: "Base URL is required"},
			Errors: map[string][]string{
				"base_url": {"Base URL is required"},
			},
		}
	}

	if c.Timeout <= 0 {
		return &ValidationError{
			BaseError: BaseError{Message: "Timeout must be greater than 0"},
			Errors: map[string][]string{
				"timeout": {"Timeout must be greater than 0"},
			},
		}
	}

	if c.ConnectTimeout <= 0 {
		return &ValidationError{
			BaseError: BaseError{Message: "Connect timeout must be greater than 0"},
			Errors: map[string][]string{
				"connect_timeout": {"Connect timeout must be greater than 0"},
			},
		}
	}

	return nil
}

// GetUserAgent returns the User-Agent string for HTTP requests
func (c *Config) GetUserAgent() string {
	return fmt.Sprintf("poodle-go/%s", SDKVersion)
}
