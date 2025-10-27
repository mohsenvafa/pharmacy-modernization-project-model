# Central Sanitizer Package

This package provides centralized sanitization functions for various input types to prevent security vulnerabilities like log injection, SQL injection, XSS, and other attacks.

## Features

- **Log Sanitization**: Prevents log injection attacks by removing control characters
- **URL Sanitization**: Safely handles URLs for logging and processing
- **HTML Sanitization**: Escapes HTML special characters to prevent XSS
- **JSON Sanitization**: Ensures safe JSON serialization
- **Database Sanitization**: Removes SQL injection vectors
- **Filename Sanitization**: Creates safe filenames
- **Email Validation**: Validates and sanitizes email addresses
- **Phone Number Sanitization**: Normalizes phone numbers
- **General Utilities**: Truncation and control character removal

## Usage

### Basic Usage

```go
import "pharmacy-modernization-project-model/internal/platform/sanitizer"

// Use convenience functions with the global instance
safeURL := sanitizer.ForLogging(userInput)
safeHTML := sanitizer.ForHTML(userInput)
safeFilename := sanitizer.ForFilename(userInput)
```

### Advanced Usage

```go
import "pharmacy-modernization-project-model/internal/platform/sanitizer"

// Create your own sanitizer instance
s := sanitizer.New()

// Use instance methods
safeURL := s.ForLogging(userInput)
safeHTML := s.ForHTML(userInput)
```

### Examples

#### Log Injection Prevention
```go
// Before (vulnerable)
logger.Error("User action", zap.String("url", r.URL.String()))

// After (secure)
logger.Error("User action", zap.String("url", sanitizer.ForLogging(r.URL.String())))
```

#### HTML Output Sanitization
```go
// Before (vulnerable to XSS)
fmt.Fprintf(w, "<p>%s</p>", userInput)

// After (secure)
fmt.Fprintf(w, "<p>%s</p>", sanitizer.ForHTML(userInput))
```

#### Database Query Sanitization
```go
// Before (vulnerable to SQL injection)
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)

// After (secure)
query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", sanitizer.ForDatabase(userInput))
```

#### Filename Sanitization
```go
// Before (vulnerable)
filename := userInput + ".txt"

// After (secure)
filename := sanitizer.ForFilename(userInput) + ".txt"
```

## Available Functions

### Sanitizer Methods
- `ForLogging(input string) string` - Sanitizes for log entries
- `ForURL(input string) string` - Sanitizes URLs
- `ForHTML(input string) string` - Escapes HTML characters
- `ForJSON(input string) string` - Ensures safe JSON encoding
- `ForDatabase(input string) string` - Removes SQL injection vectors
- `ForFilename(input string) string` - Creates safe filenames
- `ForEmail(input string) string` - Validates and sanitizes emails
- `ForPhone(input string) string` - Normalizes phone numbers
- `Truncate(input string, maxLength int, suffix string) string` - Truncates strings
- `RemoveControlChars(input string) string` - Removes control characters

### Convenience Functions
All methods are also available as package-level functions using the global `Default` instance:
- `sanitizer.ForLogging(input)`
- `sanitizer.ForHTML(input)`
- `sanitizer.ForDatabase(input)`
- etc.

## Security Considerations

- **Log Injection**: The `ForLogging` function removes control characters and limits length
- **SQL Injection**: The `ForDatabase` function removes dangerous characters and SQL keywords
- **XSS Prevention**: The `ForHTML` function escapes HTML special characters
- **Path Traversal**: The `ForFilename` function removes directory traversal characters
- **Email Validation**: The `ForEmail` function validates email format and removes control characters

## Testing

Run the tests with:
```bash
go test ./internal/platform/sanitizer/ -v
```

The package includes comprehensive tests for all sanitization functions with various edge cases and attack vectors.
