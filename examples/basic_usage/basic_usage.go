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

	// Example 1: Send a simple HTML email
	fmt.Println("Sending HTML email...")
	response, err := client.SendHTML(
		"sender@yourdomain.com",
		"recipient@example.com",
		"Hello from Poodle Go SDK!",
		"<h1>Welcome!</h1><p>This is a test email sent using the Poodle Go SDK.</p>",
	)

	if err != nil {
		log.Printf("Failed to send HTML email: %v", err)
	} else {
		fmt.Printf("HTML email sent successfully! Success: %t, Message: %s\n",
			response.Success, response.Message)
	}

	// Example 2: Send a plain text email
	fmt.Println("\nSending text email...")
	response, err = client.SendText(
		"sender@yourdomain.com",
		"recipient@example.com",
		"Plain Text Email",
		"This is a plain text email sent using the Poodle Go SDK.",
	)

	if err != nil {
		log.Printf("Failed to send text email: %v", err)
	} else {
		fmt.Printf("Text email sent successfully! Success: %t, Message: %s\n",
			response.Success, response.Message)
	}

	// Example 3: Send an email with both HTML and text content
	fmt.Println("\nSending email with both HTML and text...")
	response, err = client.SendWithBoth(
		"sender@yourdomain.com",
		"recipient@example.com",
		"Multi-format Email",
		"<h1>Hello!</h1><p>This email has both HTML and text versions.</p>",
		"Hello! This email has both HTML and text versions.",
	)

	if err != nil {
		log.Printf("Failed to send multi-format email: %v", err)
	} else {
		fmt.Printf("Multi-format email sent successfully! Success: %t, Message: %s\n",
			response.Success, response.Message)
	}

	// Example 4: Using the Email model
	fmt.Println("\nSending email using Email model...")
	email := poodle.NewEmailWithBoth(
		"sender@yourdomain.com",
		"recipient@example.com",
		"Email Model Example",
		"<h1>Using Email Model</h1><p>This email was created using the Email model.</p>",
		"Using Email Model\n\nThis email was created using the Email model.",
	)

	response, err = client.Send(email)
	if err != nil {
		log.Printf("Failed to send email using model: %v", err)
	} else {
		fmt.Printf("Email model sent successfully! Success: %t, Message: %s\n",
			response.Success, response.Message)
	}

	fmt.Println("\nAll examples completed!")
}
