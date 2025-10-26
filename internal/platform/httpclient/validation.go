package httpclient

import (
	"fmt"
	"net/url"
	"strings"
)

// URLValidator provides secure URL validation and parameter replacement
type URLValidator struct{}

// NewURLValidator creates a new URL validator instance
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ReplacePathParams safely replaces {paramName} in URL with validated actual values
func (v *URLValidator) ReplacePathParams(templateURL string, params map[string]string) (string, error) {
	// Parse the template URL to validate it
	parsedURL, err := url.Parse(templateURL)
	if err != nil {
		return "", fmt.Errorf("invalid template URL: %w", err)
	}

	// Perform safe replacement
	result := templateURL
	for key, value := range params {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Validate the final URL
	finalURL, err := url.Parse(result)
	if err != nil {
		return "", fmt.Errorf("generated invalid URL: %w", err)
	}

	// Ensure the scheme and host haven't changed (prevent redirects)
	if parsedURL.Scheme != finalURL.Scheme || parsedURL.Host != finalURL.Host {
		return "", fmt.Errorf("URL replacement changed scheme or host")
	}

	return result, nil
}

// Convenience function using a default validator instance
func ReplacePathParams(templateURL string, params map[string]string) (string, error) {
	validator := NewURLValidator()
	return validator.ReplacePathParams(templateURL, params)
}
