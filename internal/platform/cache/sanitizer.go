package cache

import (
	"regexp"
	"strings"
)

// SanitizeKey sanitizes a cache key to prevent NoSQL injection and other security issues
// It removes or escapes characters that could be used maliciously in MongoDB queries
func SanitizeKey(key string) string {
	if key == "" {
		return ""
	}

	// Remove any characters that could be problematic in MongoDB queries
	// This includes MongoDB operators like $, ., and other special characters
	sanitized := regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(key, "_")

	// Ensure the key doesn't start or end with underscores
	sanitized = strings.Trim(sanitized, "_")

	// If the key becomes empty after sanitization, use a default
	if sanitized == "" {
		sanitized = "invalid_key"
	}

	// Limit key length to prevent excessively long keys
	if len(sanitized) > 100 {
		sanitized = sanitized[:100]
	}

	return sanitized
}

// ValidateID validates that an ID matches expected format patterns
// This helps prevent injection by ensuring IDs follow expected patterns
func ValidateID(id string) bool {
	if id == "" {
		return false
	}

	// Allow alphanumeric characters, hyphens, and underscores
	// This covers most common ID formats (UUIDs, MongoDB ObjectIds, etc.)
	validIDPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return validIDPattern.MatchString(id)
}
