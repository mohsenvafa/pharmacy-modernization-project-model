package sanitizer

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides centralized sanitization functions for various input types
type Sanitizer struct {
	// Control characters and potentially dangerous characters for logging
	logDangerousChars *regexp.Regexp
	// Characters that could be used for injection attacks
	injectionChars *regexp.Regexp
}

// New creates a new sanitizer instance
func New() *Sanitizer {
	return &Sanitizer{
		// Match control characters (0x00-0x1F, 0x7F) and other potentially dangerous characters
		logDangerousChars: regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`),
		// Match characters commonly used in injection attacks
		injectionChars: regexp.MustCompile(`[<>'"&\\]`),
	}
}

// ForLogging sanitizes input specifically for logging purposes
// Removes control characters and escapes newlines to prevent log injection
func (s *Sanitizer) ForLogging(input string) string {
	if input == "" {
		return ""
	}

	// Remove control characters and other dangerous characters
	sanitized := s.logDangerousChars.ReplaceAllString(input, "")

	// Replace newlines with escaped version to prevent log injection
	sanitized = strings.ReplaceAll(sanitized, "\n", "\\n")
	sanitized = strings.ReplaceAll(sanitized, "\r", "\\r")

	// Limit length to prevent excessively long log entries
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000] + "...[truncated]"
	}

	return sanitized
}

// ForURL sanitizes input for use in URLs
// Escapes special characters and validates URL structure
func (s *Sanitizer) ForURL(input string) string {
	if input == "" {
		return ""
	}

	// Parse and escape the URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		// If parsing fails, return a safe placeholder
		return "[invalid-url]"
	}

	// Reconstruct the URL with proper escaping
	return parsedURL.String()
}

// ForHTML sanitizes input for HTML output
// Escapes HTML special characters
func (s *Sanitizer) ForHTML(input string) string {
	if input == "" {
		return ""
	}

	// Replace HTML special characters
	sanitized := strings.ReplaceAll(input, "&", "&amp;")
	sanitized = strings.ReplaceAll(sanitized, "<", "&lt;")
	sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")
	sanitized = strings.ReplaceAll(sanitized, "\"", "&quot;")
	sanitized = strings.ReplaceAll(sanitized, "'", "&#x27;")

	return sanitized
}

// ForJSON sanitizes input for JSON serialization
// Ensures the string is safe for JSON encoding
func (s *Sanitizer) ForJSON(input string) string {
	if input == "" {
		return ""
	}

	// Use json.Marshal to properly escape the string
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "[invalid-json]"
	}

	// Remove the surrounding quotes that json.Marshal adds
	return string(jsonBytes[1 : len(jsonBytes)-1])
}

// ForDatabase sanitizes input for database queries
// Removes characters that could be used for SQL injection
func (s *Sanitizer) ForDatabase(input string) string {
	if input == "" {
		return ""
	}

	// Remove characters commonly used in SQL injection
	sanitized := s.injectionChars.ReplaceAllString(input, "")

	// Remove SQL keywords that could be dangerous
	sqlKeywords := []string{"DROP", "DELETE", "INSERT", "UPDATE", "SELECT", "UNION", "ALTER", "CREATE"}
	for _, keyword := range sqlKeywords {
		// Case-insensitive replacement
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		sanitized = re.ReplaceAllString(sanitized, "[removed]")
	}

	return sanitized
}

// ForFilename sanitizes input for use as a filename
// Removes dangerous characters and ensures valid filename
func (s *Sanitizer) ForFilename(input string) string {
	if input == "" {
		return "unnamed"
	}

	// Remove or replace dangerous characters
	sanitized := regexp.MustCompile(`[<>:"/\\|?*]`).ReplaceAllString(input, "_")

	// Remove control characters
	sanitized = regexp.MustCompile(`[\x00-\x1F\x7F]`).ReplaceAllString(sanitized, "")

	// Remove leading/trailing dots and spaces
	sanitized = strings.Trim(sanitized, ". ")

	// Ensure it's not empty after sanitization
	if sanitized == "" {
		sanitized = "unnamed"
	}

	// Limit length
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}

// ForEmail sanitizes and validates email input
func (s *Sanitizer) ForEmail(input string) string {
	if input == "" {
		return ""
	}

	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input) {
		return "[invalid-email]"
	}

	// Remove any control characters
	sanitized := s.logDangerousChars.ReplaceAllString(input, "")

	return sanitized
}

// ForPhone sanitizes phone number input
func (s *Sanitizer) ForPhone(input string) string {
	if input == "" {
		return ""
	}

	// Remove all non-digit characters except + at the beginning
	sanitized := regexp.MustCompile(`[^\d+]`).ReplaceAllString(input, "")

	// Ensure it starts with + if it contains international format
	if len(sanitized) > 10 && !strings.HasPrefix(sanitized, "+") {
		sanitized = "+" + sanitized
	}

	// Limit length
	if len(sanitized) > 20 {
		sanitized = sanitized[:20]
	}

	return sanitized
}

// Truncate truncates a string to the specified length with optional suffix
func (s *Sanitizer) Truncate(input string, maxLength int, suffix string) string {
	if input == "" || maxLength <= 0 {
		return ""
	}

	if len(input) <= maxLength {
		return input
	}

	if suffix == "" {
		suffix = "..."
	}

	truncateLength := maxLength - len(suffix)
	if truncateLength <= 0 {
		return suffix
	}

	return input[:truncateLength] + suffix
}

// RemoveControlChars removes all control characters from input
func (s *Sanitizer) RemoveControlChars(input string) string {
	if input == "" {
		return ""
	}

	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1 // Remove the character
		}
		return r
	}, input)
}

// Global sanitizer instance for convenience
var Default = New()

// Convenience functions using the global instance
func ForLogging(input string) string {
	return Default.ForLogging(input)
}

func ForURL(input string) string {
	return Default.ForURL(input)
}

func ForHTML(input string) string {
	return Default.ForHTML(input)
}

func ForJSON(input string) string {
	return Default.ForJSON(input)
}

func ForDatabase(input string) string {
	return Default.ForDatabase(input)
}

func ForFilename(input string) string {
	return Default.ForFilename(input)
}

func ForEmail(input string) string {
	return Default.ForEmail(input)
}

func ForPhone(input string) string {
	return Default.ForPhone(input)
}

func Truncate(input string, maxLength int, suffix string) string {
	return Default.Truncate(input, maxLength, suffix)
}

func RemoveControlChars(input string) string {
	return Default.RemoveControlChars(input)
}
