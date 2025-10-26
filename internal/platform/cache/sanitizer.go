package cache

import (
	"fmt"
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
// This helps prevent injection by ensuring IDs follow expected patterns.
// The validation is strict to ensure only safe alphanumeric characters, hyphens, and underscores are allowed.
func ValidateID(id string) bool {
	if id == "" {
		return false
	}

	// Reject keys that are too short (less than 1 character) or too long (over 200 characters)
	// This prevents both empty keys and excessively long keys that could be used for DoS
	if len(id) < 1 || len(id) > 200 {
		return false
	}

	// Allow alphanumeric characters, hyphens, and underscores only
	// This covers most common ID formats (UUIDs, MongoDB ObjectIds, slugs, etc.)
	// The pattern must match the entire string (^...$) to prevent partial matches
	validIDPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	if !validIDPattern.MatchString(id) {
		return false
	}

	// Additional validation: ensure the key doesn't contain MongoDB operators
	// This is a defense-in-depth measure even though the regex above should catch these
	mongoOperators := []string{"$", ".", "(", ")", "[", "]", "{", "}", "*", "+", "?", "|", "^", "\\"}
	for _, op := range mongoOperators {
		if strings.Contains(id, op) {
			return false
		}
	}

	return true
}

// ValidateIDWithReason performs validation and returns a detailed reason for failure
// This is useful for logging and debugging purposes
func ValidateIDWithReason(id string) (bool, string) {
	if id == "" {
		return false, "key is empty"
	}

	if len(id) < 1 || len(id) > 200 {
		return false, fmt.Sprintf("key length invalid: %d (must be 1-200 characters)", len(id))
	}

	validIDPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validIDPattern.MatchString(id) {
		return false, "key contains invalid characters (only alphanumeric, hyphens, and underscores allowed)"
	}

	mongoOperators := []string{"$", ".", "(", ")", "[", "]", "{", "}", "*", "+", "?", "|", "^", "\\"}
	for _, op := range mongoOperators {
		if strings.Contains(id, op) {
			return false, fmt.Sprintf("key contains forbidden MongoDB operator: %s", op)
		}
	}

	return true, ""
}
