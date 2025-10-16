package httpclient

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// URLValidator provides secure URL validation and parameter replacement
type URLValidator struct{}

// NewURLValidator creates a new URL validator instance
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// ValidateParamValue validates that a parameter value is safe for URL insertion
func (v *URLValidator) ValidateParamValue(value string) error {
	// Check for common injection patterns
	if strings.Contains(value, "://") || strings.Contains(value, "javascript:") || strings.Contains(value, "data:") {
		return fmt.Errorf("invalid parameter value: contains protocol scheme")
	}

	// Check for path traversal attempts
	if strings.Contains(value, "..") || strings.Contains(value, "~") {
		return fmt.Errorf("invalid parameter value: contains path traversal characters")
	}

	// Check for control characters and non-printable characters
	controlCharRegex := regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`)
	if controlCharRegex.MatchString(value) {
		return fmt.Errorf("invalid parameter value: contains control characters")
	}

	// URL encode the value to ensure it's safe
	encoded := url.QueryEscape(value)
	if encoded != value {
		return fmt.Errorf("invalid parameter value: contains characters that need URL encoding")
	}

	return nil
}

// ValidateURL performs security validation on a request URL
func (v *URLValidator) ValidateURL(requestURL string) error {
	// Parse the URL to validate it
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Ensure the URL has a scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http/https)")
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	// Ensure the URL has a host
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	// Check for dangerous patterns in the URL
	if strings.Contains(requestURL, "javascript:") || strings.Contains(requestURL, "data:") {
		return fmt.Errorf("URL contains dangerous protocol")
	}

	return nil
}

// ReplacePathParams safely replaces {paramName} in URL with validated actual values
func (v *URLValidator) ReplacePathParams(templateURL string, params map[string]string) (string, error) {
	// Parse the template URL to validate it
	parsedURL, err := url.Parse(templateURL)
	if err != nil {
		return "", fmt.Errorf("invalid template URL: %w", err)
	}

	// Validate all parameter values before replacement
	for key, value := range params {
		if err := v.ValidateParamValue(value); err != nil {
			return "", fmt.Errorf("invalid parameter %s: %w", key, err)
		}
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

// Global validator instance for convenience
var DefaultURLValidator = NewURLValidator()

// Convenience functions using the default validator
func ValidateParamValue(value string) error {
	return DefaultURLValidator.ValidateParamValue(value)
}

func ValidateURL(requestURL string) error {
	return DefaultURLValidator.ValidateURL(requestURL)
}

func ReplacePathParams(templateURL string, params map[string]string) (string, error) {
	return DefaultURLValidator.ReplacePathParams(templateURL, params)
}
