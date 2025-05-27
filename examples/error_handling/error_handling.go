package main

import (
	"fmt"
	"log"
	"os"

	"github.com/usepoodle/poodle-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("POODLE_API_KEY")
	if apiKey == "" {
		log.Fatal("POODLE_API_KEY environment variable is required")
	}

	// Initialize the Poodle client
	client := poodle.NewClient(apiKey)

	// Example 1: Handle validation errors
	fmt.Println("Example 1: Validation Error")
	response, err := client.SendHTML(
		"invalid-email", // Invalid email address
		"recipient@example.com",
		"Test Email",
		"<h1>Hello!</h1>",
	)

	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Email sent successfully: %s\n", response.Message)
	}

	// Example 2: Handle missing content error
	fmt.Println("\nExample 2: Missing Content Error")
	email := &poodle.Email{
		From:    "sender@yourdomain.com",
		To:      "recipient@example.com",
		Subject: "Test Email",
		// No HTML or Text content - should cause validation error
	}

	response, err = client.Send(email)
	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Email sent successfully: %s\n", response.Message)
	}

	// Example 3: Handle authentication error (with invalid API key)
	fmt.Println("\nExample 3: Authentication Error")
	invalidClient := poodle.NewClient("invalid_api_key_123")

	response, err = invalidClient.SendHTML(
		"sender@yourdomain.com",
		"recipient@example.com",
		"Test Email",
		"<h1>Hello!</h1>",
	)

	if err != nil {
		handleError(err)
	} else {
		fmt.Printf("Email sent successfully: %s\n", response.Message)
	}

	fmt.Println("\nError handling examples completed!")
}

// handleError demonstrates how to handle different types of Poodle errors
func handleError(err error) {
	fmt.Printf("Error occurred: %s\n", err.Error())

	// Handle specific error types
	switch e := err.(type) {
	case *poodle.ValidationError:
		fmt.Println("  Type: Validation Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Println("  Field Errors:")
		for field, errors := range e.Errors {
			fmt.Printf("    %s: %v\n", field, errors)
		}
		fmt.Println("  Context:", e.Context())

	case *poodle.AuthenticationError:
		fmt.Println("  Type: Authentication Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Println("  Suggestion: Check your API key and ensure it's valid")
		fmt.Println("  Context:", e.Context())

	case *poodle.RateLimitError:
		fmt.Println("  Type: Rate Limit Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Printf("  Retry After: %d seconds\n", e.RetryAfter)
		fmt.Printf("  Limit: %d\n", e.Limit)
		fmt.Printf("  Remaining: %d\n", e.Remaining)
		fmt.Printf("  Reset Time: %d\n", e.Reset)
		fmt.Println("  Suggestion: Wait before retrying")
		fmt.Println("  Context:", e.Context())

	case *poodle.SubscriptionError:
		fmt.Println("  Type: Subscription Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Printf("  Error Type: %s\n", e.ErrorType)
		fmt.Println("  Suggestion: Check your subscription status")
		fmt.Println("  Context:", e.Context())

	case *poodle.AccountSuspendedError:
		fmt.Println("  Type: Account Suspended Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Printf("  Reason: %s\n", e.Reason)
		fmt.Println("  Suggestion: Contact support")
		fmt.Println("  Context:", e.Context())

	case *poodle.NetworkError:
		fmt.Println("  Type: Network Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Printf("  URL: %s\n", e.URL)
		fmt.Println("  Suggestion: Check your internet connection and try again")
		fmt.Println("  Context:", e.Context())

	case *poodle.HTTPError:
		fmt.Println("  Type: HTTP Error")
		fmt.Printf("  Status Code: %d\n", e.StatusCode())
		fmt.Printf("  URL: %s\n", e.URL)
		fmt.Printf("  Response Body: %s\n", e.ResponseBody)
		fmt.Println("  Context:", e.Context())

	default:
		fmt.Println("  Type: Unknown Error")
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println()
}
