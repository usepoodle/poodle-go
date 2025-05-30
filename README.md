# Poodle Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.20%2B-blue.svg)](https://golang.org)
[![Build Status](https://github.com/usepoodle/poodle-go/workflows/CI/badge.svg)](https://github.com/usepoodle/poodle-go/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/usepoodle/poodle-go/blob/main/LICENSE)

Go SDK for Poodle's email sending API.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Features](#features)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Error Handling](#error-handling)
- [Contributing](#contributing)
- [License](#license)

## Installation

Install the SDK using Go modules:

```bash
go get github.com/usepoodle/poodle-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/usepoodle/poodle-go"
)

func main() {
    // Initialize the client
    client := poodle.NewClient("your_api_key_here")

    // Send an email
    response, err := client.SendHTML(
        "sender@yourdomain.com",
        "recipient@example.com",
        "Hello from Poodle!",
        "<h1>Hello World!</h1><p>This is a test email.</p>",
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Email sent! Success: %t, Message: %s\n",
        response.Success, response.Message)
}
```

## Features

- Simple and intuitive API
- HTML and plain text email support
- Comprehensive error handling
- Built-in input validation
- Go 1.20+ support
- Goroutine-safe client design
- Extensive test suite
- Zero external dependencies

## Configuration

### API Key

Set your API key in one of these ways:

**1. Pass directly to constructor:**

```go
client := poodle.NewClient("your_api_key_here")
```

**2. Use environment variable:**

```bash
export POODLE_API_KEY=your_api_key_here
```

```go
client := poodle.NewClientFromEnv()
```

**3. Use Configuration object:**

```go
config := &poodle.Config{
    APIKey:         "your_api_key_here",
    BaseURL:        "https://api.usepoodle.com",
    Timeout:        30 * time.Second,
    ConnectTimeout: 10 * time.Second,
    Debug:          true,
}

client := poodle.NewClientWithConfig(config)
```

### Environment Variables

| Variable                 | Default                     | Description          |
| ------------------------ | --------------------------- | -------------------- |
| `POODLE_API_KEY`         | -                           | Your Poodle API key  |
| `POODLE_BASE_URL`        | `https://api.usepoodle.com` | API base URL         |
| `POODLE_TIMEOUT`         | `30s`                       | Request timeout      |
| `POODLE_CONNECT_TIMEOUT` | `10s`                       | Connection timeout   |
| `POODLE_DEBUG`           | `false`                     | Enable debug logging |

## Usage Examples

### Basic Email Sending

```go
client := poodle.NewClient("your_api_key")

// HTML email
response, err := client.SendHTML(
    "sender@yourdomain.com",
    "recipient@example.com",
    "Welcome!",
    "<h1>Welcome to our service!</h1>",
)

// Plain text email
response, err := client.SendText(
    "sender@yourdomain.com",
    "recipient@example.com",
    "Welcome!",
    "Welcome to our service!",
)

// Both HTML and text
response, err := client.SendWithBoth(
    "sender@yourdomain.com",
    "recipient@example.com",
    "Welcome!",
    "<h1>Welcome!</h1>",
    "Welcome!",
)
```

### Using the Email Model

```go
email := &poodle.Email{
    From:    "sender@yourdomain.com",
    To:      "recipient@example.com",
    Subject: "Welcome Email",
    HTML:    "<h1>Hello!</h1><p>Welcome to our service!</p>",
    Text:    "Hello! Welcome to our service!",
}

response, err := client.Send(email)
if err != nil {
    log.Fatal(err)
}

if response.Success {
    fmt.Println("Email queued successfully!")
}
```

### Error Handling

```go
response, err := client.SendHTML(
    "sender@yourdomain.com",
    "recipient@example.com",
    "Test Email",
    "<h1>Hello!</h1>",
)

if err != nil {
    switch e := err.(type) {
    case *poodle.ValidationError:
        fmt.Printf("Validation error: %s\n", e.Error())
        for field, errors := range e.Errors {
            fmt.Printf("  %s: %v\n", field, errors)
        }
    case *poodle.AuthenticationError:
        fmt.Printf("Authentication failed: %s\n", e.Error())
    case *poodle.RateLimitError:
        fmt.Printf("Rate limit exceeded. Retry after: %d seconds\n", e.RetryAfter)
    case *poodle.NetworkError:
        fmt.Printf("Network error: %s\n", e.Error())
    default:
        fmt.Printf("Unknown error: %s\n", err.Error())
    }
    return
}

fmt.Printf("Email sent successfully! Message: %s\n", response.Message)
```

## API Reference

### Client

#### `NewClient(apiKey string) *Client`

Creates a new client with the provided API key.

#### `NewClientFromEnv() *Client`

Creates a new client using environment variables.

#### `NewClientWithConfig(config *Config) *Client`

Creates a new client with custom configuration.

### Methods

#### `Send(email *Email) (*EmailResponse, error)`

Sends an email using the Email model.

#### `SendHTML(from, to, subject, html string) (*EmailResponse, error)`

Sends an HTML email.

#### `SendText(from, to, subject, text string) (*EmailResponse, error)`

Sends a plain text email.

#### `SendWithBoth(from, to, subject, html, text string) (*EmailResponse, error)`

Sends an email with both HTML and text content.

### Types

#### `Email`

```go
type Email struct {
    From    string `json:"from"`
    To      string `json:"to"`
    Subject string `json:"subject"`
    HTML    string `json:"html,omitempty"`
    Text    string `json:"text,omitempty"`
}
```

#### `EmailResponse`

```go
type EmailResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}
```

#### `Config`

```go
type Config struct {
    APIKey         string
    BaseURL        string
    Timeout        time.Duration
    ConnectTimeout time.Duration
    Debug          bool
}
```

## Error Handling

The SDK provides specific error types for different scenarios:

- `ValidationError` - Invalid request data (400)
- `AuthenticationError` - Invalid or missing API key (401)
- `AccountSuspendedError` - Account suspended (403)
- `SubscriptionError` - Subscription issues (402)
- `RateLimitError` - Rate limit exceeded (429)
- `NetworkError` - Network connectivity issues

Each error type provides additional context and methods for handling specific scenarios.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](https://github.com/usepoodle/poodle-go/blob/main/CONTRIBUTING.md) for details on the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/usepoodle/poodle-go/blob/main/LICENSE) file for details.
