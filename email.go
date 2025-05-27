package poodle

import (
	"regexp"
	"strings"
)

// Email represents an email to be sent
type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	HTML    string `json:"html,omitempty"`
	Text    string `json:"text,omitempty"`
}

// Email validation constants
const (
	MaxContentSize = 10 * 1024 * 1024 // 10MB
)

// Email address validation regex (RFC 5322 compliant)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email instance
func NewEmail(from, to, subject string) *Email {
	return &Email{
		From:    from,
		To:      to,
		Subject: subject,
	}
}

// NewHTMLEmail creates a new Email instance with HTML content
func NewHTMLEmail(from, to, subject, html string) *Email {
	return &Email{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
	}
}

// NewTextEmail creates a new Email instance with text content
func NewTextEmail(from, to, subject, text string) *Email {
	return &Email{
		From:    from,
		To:      to,
		Subject: subject,
		Text:    text,
	}
}

// NewEmailWithBoth creates a new Email instance with both HTML and text content
func NewEmailWithBoth(from, to, subject, html, text string) *Email {
	return &Email{
		From:    from,
		To:      to,
		Subject: subject,
		HTML:    html,
		Text:    text,
	}
}

// Validate validates the email data
func (e *Email) Validate() error {
	errors := make(map[string][]string)

	// Validate required fields
	if strings.TrimSpace(e.From) == "" {
		errors["from"] = append(errors["from"], "From address is required")
	} else if !isValidEmail(e.From) {
		errors["from"] = append(errors["from"], "From address is not a valid email")
	}

	if strings.TrimSpace(e.To) == "" {
		errors["to"] = append(errors["to"], "To address is required")
	} else if !isValidEmail(e.To) {
		errors["to"] = append(errors["to"], "To address is not a valid email")
	}

	if strings.TrimSpace(e.Subject) == "" {
		errors["subject"] = append(errors["subject"], "Subject is required")
	}

	// Validate content - at least one of HTML or Text must be provided
	if strings.TrimSpace(e.HTML) == "" && strings.TrimSpace(e.Text) == "" {
		errors["content"] = append(errors["content"], "At least one content type (html or text) is required")
	}

	// Validate content size
	if len(e.HTML) > MaxContentSize {
		errors["html"] = append(errors["html"], "HTML content exceeds maximum size limit")
	}

	if len(e.Text) > MaxContentSize {
		errors["text"] = append(errors["text"], "Text content exceeds maximum size limit")
	}

	if len(errors) > 0 {
		return NewValidationError("Email validation failed", errors)
	}

	return nil
}

// SetHTML sets the HTML content
func (e *Email) SetHTML(html string) *Email {
	e.HTML = html
	return e
}

// SetText sets the text content
func (e *Email) SetText(text string) *Email {
	e.Text = text
	return e
}

// SetBoth sets both HTML and text content
func (e *Email) SetBoth(html, text string) *Email {
	e.HTML = html
	e.Text = text
	return e
}

// HasHTML returns true if the email has HTML content
func (e *Email) HasHTML() bool {
	return strings.TrimSpace(e.HTML) != ""
}

// HasText returns true if the email has text content
func (e *Email) HasText() bool {
	return strings.TrimSpace(e.Text) != ""
}

// isValidEmail validates email address format
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 254 {
		return false
	}

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	localPart := parts[0]
	domainPart := parts[1]

	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") || strings.Contains(localPart, "..") {
		return false
	}

	if strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") ||
		strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") || // Technically suffix dot might be for FQDN but generally not for email context
		strings.Contains(domainPart, "..") { // Prevent domain..com
		return false
	}
	domainLabels := strings.Split(domainPart, ".")
	if len(domainLabels) < 2 { // e.g. @localhost
		return false
	}
	for _, label := range domainLabels {
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") || len(label) == 0 {
			return false // Each label cannot start or end with a hyphen or be empty
		}
	}

	return true
}
