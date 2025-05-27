package poodle

import (
	"strings"
	"testing"
)

func TestNewEmail(t *testing.T) {
	email := NewEmail("from@example.com", "to@example.com", "Test Subject")

	if email.From != "from@example.com" {
		t.Errorf("Expected From to be 'from@example.com', got '%s'", email.From)
	}
	if email.To != "to@example.com" {
		t.Errorf("Expected To to be 'to@example.com', got '%s'", email.To)
	}
	if email.Subject != "Test Subject" {
		t.Errorf("Expected Subject to be 'Test Subject', got '%s'", email.Subject)
	}
}

func TestNewHTMLEmail(t *testing.T) {
	html := "<h1>Hello</h1>"
	email := NewHTMLEmail("from@example.com", "to@example.com", "Test Subject", html)

	if email.HTML != html {
		t.Errorf("Expected HTML to be '%s', got '%s'", html, email.HTML)
	}
	if email.Text != "" {
		t.Errorf("Expected Text to be empty, got '%s'", email.Text)
	}
}

func TestNewTextEmail(t *testing.T) {
	text := "Hello World"
	email := NewTextEmail("from@example.com", "to@example.com", "Test Subject", text)

	if email.Text != text {
		t.Errorf("Expected Text to be '%s', got '%s'", text, email.Text)
	}
	if email.HTML != "" {
		t.Errorf("Expected HTML to be empty, got '%s'", email.HTML)
	}
}

func TestNewEmailWithBoth(t *testing.T) {
	html := "<h1>Hello</h1>"
	text := "Hello World"
	email := NewEmailWithBoth("from@example.com", "to@example.com", "Test Subject", html, text)

	if email.HTML != html {
		t.Errorf("Expected HTML to be '%s', got '%s'", html, email.HTML)
	}
	if email.Text != text {
		t.Errorf("Expected Text to be '%s', got '%s'", text, email.Text)
	}
}

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name        string
		email       *Email
		expectError bool
		errorFields []string
	}{
		{
			name: "Valid HTML email",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
			},
			expectError: false,
		},
		{
			name: "Valid text email",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
				Text:    "Hello World",
			},
			expectError: false,
		},
		{
			name: "Valid email with both HTML and text",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
				Text:    "Hello World",
			},
			expectError: false,
		},
		{
			name: "Missing from address",
			email: &Email{
				To:      "to@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
			},
			expectError: true,
			errorFields: []string{"from"},
		},
		{
			name: "Invalid from address",
			email: &Email{
				From:    "invalid-email",
				To:      "to@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
			},
			expectError: true,
			errorFields: []string{"from"},
		},
		{
			name: "Missing to address",
			email: &Email{
				From:    "from@example.com",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
			},
			expectError: true,
			errorFields: []string{"to"},
		},
		{
			name: "Invalid to address",
			email: &Email{
				From:    "from@example.com",
				To:      "invalid-email",
				Subject: "Test Subject",
				HTML:    "<h1>Hello</h1>",
			},
			expectError: true,
			errorFields: []string{"to"},
		},
		{
			name: "Missing subject",
			email: &Email{
				From: "from@example.com",
				To:   "to@example.com",
				HTML: "<h1>Hello</h1>",
			},
			expectError: true,
			errorFields: []string{"subject"},
		},
		{
			name: "Missing content",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
			},
			expectError: true,
			errorFields: []string{"content"},
		},
		{
			name: "HTML content too large",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
				HTML:    strings.Repeat("a", MaxContentSize+1),
			},
			expectError: true,
			errorFields: []string{"html"},
		},
		{
			name: "Text content too large",
			email: &Email{
				From:    "from@example.com",
				To:      "to@example.com",
				Subject: "Test Subject",
				Text:    strings.Repeat("a", MaxContentSize+1),
			},
			expectError: true,
			errorFields: []string{"text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.email.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected validation error, but got none")
					return
				}

				validationErr, ok := err.(*ValidationError)
				if !ok {
					t.Errorf("Expected ValidationError, got %T", err)
					return
				}

				for _, field := range tt.errorFields {
					if _, exists := validationErr.Errors[field]; !exists {
						t.Errorf("Expected error for field '%s', but not found in errors: %v", field, validationErr.Errors)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, but got: %v", err)
				}
			}
		})
	}
}

func TestEmailMethods(t *testing.T) {
	email := NewEmail("from@example.com", "to@example.com", "Test Subject")

	// Test SetHTML
	email.SetHTML("<h1>Hello</h1>")
	if !email.HasHTML() {
		t.Error("Expected email to have HTML content")
	}
	if email.HasText() {
		t.Error("Expected email to not have text content")
	}

	// Test SetText
	email.SetText("Hello World")
	if !email.HasText() {
		t.Error("Expected email to have text content")
	}

	// Test SetBoth
	email2 := NewEmail("from@example.com", "to@example.com", "Test Subject")
	email2.SetBoth("<h1>Hello</h1>", "Hello World")
	if !email2.HasHTML() {
		t.Error("Expected email to have HTML content")
	}
	if !email2.HasText() {
		t.Error("Expected email to have text content")
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@example.com", true},
		{"user+tag@example.com", true},
		{"user123@example-domain.com", true},
		{"", false},
		{"invalid", false},
		{"@example.com", false},
		{"test@", false},
		{"test@.com", false},
		{"test..test@example.com", false},
		{strings.Repeat("a", 250) + "@example.com", false}, // Too long
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := isValidEmail(tt.email)
			if result != tt.valid {
				t.Errorf("isValidEmail(%s) = %v, want %v", tt.email, result, tt.valid)
			}
		})
	}
}
