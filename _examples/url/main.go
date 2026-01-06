package main

import (
	"fmt"

	"github.com/hotfixfirst/go-xvalidator"
)

// Website demonstrates URL validation
type Website struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	URL         string `json:"url" validate:"required,url"`
	SecureURL   string `json:"secure_url" validate:"omitempty,https_url"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}

// APIEndpoint demonstrates HTTPS-only validation
type APIEndpoint struct {
	Name        string `json:"name" validate:"required"`
	BaseURL     string `json:"base_url" validate:"required,https_url"`
	WebhookURL  string `json:"webhook_url" validate:"required,https_url"`
	CallbackURL string `json:"callback_url" validate:"omitempty,https_url"`
}

func main() {
	fmt.Println("=== URL Validation Examples ===\n")

	// Create validator instance
	v, err := xvalidator.NewValidator()
	if err != nil {
		panic(err)
	}

	// Example 1: Valid HTTP URL
	fmt.Println("Example 1: Valid HTTP URL")
	site1 := Website{
		Name:        "Example Site",
		URL:         "http://example.com",
		Description: "This is an example website",
	}
	if err := v.Validate(site1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 2: Valid HTTPS URL
	fmt.Println("Example 2: Valid HTTPS URL")
	site2 := Website{
		Name:        "Secure Site",
		URL:         "https://secure.example.com",
		SecureURL:   "https://api.example.com",
		Description: "This is a secure website",
	}
	if err := v.Validate(site2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 3: Valid URL with path
	fmt.Println("Example 3: Valid URL with Path")
	site3 := Website{
		Name:        "Blog Site",
		URL:         "https://blog.example.com/posts/hello-world",
		Description: "Blog post URL",
	}
	if err := v.Validate(site3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 4: Valid URL with query parameters
	fmt.Println("Example 4: Valid URL with Query Parameters")
	site4 := Website{
		Name:        "Search Page",
		URL:         "https://example.com/search?q=golang&page=1",
		Description: "Search results page",
	}
	if err := v.Validate(site4); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 5: Valid URL with port
	fmt.Println("Example 5: Valid URL with Port")
	site5 := Website{
		Name:        "Local Server",
		URL:         "http://localhost:8080/api",
		Description: "Local development server",
	}
	if err := v.Validate(site5); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 6: Invalid - missing protocol
	fmt.Println("Example 6: Invalid - Missing Protocol")
	invalid1 := Website{
		Name:        "Invalid Site",
		URL:         "example.com",
		Description: "Missing protocol",
	}
	if err := v.Validate(invalid1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 7: Invalid - malformed URL
	fmt.Println("Example 7: Invalid - Malformed URL")
	invalid2 := Website{
		Name:        "Invalid Site",
		URL:         "http://",
		Description: "Incomplete URL",
	}
	if err := v.Validate(invalid2); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 8: Invalid - spaces in URL
	fmt.Println("Example 8: Invalid - Spaces in URL")
	invalid3 := Website{
		Name:        "Invalid Site",
		URL:         "http://example .com",
		Description: "URL with space",
	}
	if err := v.Validate(invalid3); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 9: Valid API endpoint (HTTPS only)
	fmt.Println("Example 9: Valid API Endpoint (HTTPS Only)")
	api1 := APIEndpoint{
		Name:        "Payment Gateway",
		BaseURL:     "https://api.payment.com",
		WebhookURL:  "https://myapp.com/webhook",
		CallbackURL: "https://myapp.com/callback",
	}
	if err := v.Validate(api1); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 10: Invalid - HTTP not allowed for API
	fmt.Println("Example 10: Invalid - HTTP Not Allowed for API")
	invalidAPI := APIEndpoint{
		Name:       "Insecure API",
		BaseURL:    "http://api.example.com",
		WebhookURL: "https://myapp.com/webhook",
	}
	if err := v.Validate(invalidAPI); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Validation passed!")
	}
	fmt.Println()

	// Example 11: Single URL validation
	fmt.Println("Example 11: Single URL Field Validation")
	url := "https://golang.org"
	if err := v.Var(url, "url"); err != nil {
		fmt.Printf("❌ URL '%s' validation failed: %v\n", url, err)
	} else {
		fmt.Printf("✅ URL '%s' is valid!\n", url)
	}
	fmt.Println()

	// Example 12: Single HTTPS URL validation
	fmt.Println("Example 12: Single HTTPS URL Validation")
	httpsURL := "https://secure.example.com"
	if err := v.Var(httpsURL, "https_url"); err != nil {
		fmt.Printf("❌ HTTPS URL '%s' validation failed: %v\n", httpsURL, err)
	} else {
		fmt.Printf("✅ HTTPS URL '%s' is valid!\n", httpsURL)
	}
	fmt.Println()

	// Example 13: Single invalid HTTPS validation (HTTP provided)
	fmt.Println("Example 13: Invalid - HTTP URL for HTTPS Validator")
	httpURL := "http://example.com"
	if err := v.Var(httpURL, "https_url"); err != nil {
		fmt.Printf("❌ HTTPS URL '%s' validation failed: %v\n", httpURL, err)
	} else {
		fmt.Printf("✅ HTTPS URL '%s' is valid!\n", httpURL)
	}
}
