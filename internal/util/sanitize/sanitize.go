// Package sanitize provides input sanitization utilities for security
package sanitize

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	// MaxStringLength is the maximum allowed length for string inputs
	MaxStringLength = 1000
	// MaxQueryLength is the maximum allowed length for query strings
	MaxQueryLength = 500
)

var (
	// dangerousChars contains characters that could be used in injection attacks
	dangerousChars = []string{";", "'", "\"", "\n", "\r", "\x00", "\x1a", "\x1b", "\x1c", "\x1d", "\x1e", "\x1f"}
	// sqlKeywords contains SQL keywords that should be sanitized in non-SQL contexts
	sqlKeywords = []string{"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE", "TRUNCATE", "UNION", "EXEC", "SCRIPT"}
	// shellMetachars contains shell metacharacters that could be dangerous
	shellMetachars = []string{"|", "&", ";", "$", "(", ")", "<", ">", "`", "$(", "${"}
)

// SanitizeString sanitizes a string input by removing potentially dangerous characters
func SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// Remove dangerous characters
	result := input
	for _, char := range dangerousChars {
		result = strings.ReplaceAll(result, char, "")
	}

	// Remove shell metacharacters
	for _, meta := range shellMetachars {
		result = strings.ReplaceAll(result, meta, "")
	}

	// Limit length
	if len(result) > MaxStringLength {
		result = result[:MaxStringLength]
	}

	// Trim whitespace
	result = strings.TrimSpace(result)

	return result
}

// SanitizeQuery sanitizes a query string (more aggressive sanitization)
func SanitizeQuery(input string) string {
	if input == "" {
		return input
	}

	result := SanitizeString(input)

	// Remove SQL keywords (case-insensitive)
	for _, keyword := range sqlKeywords {
		regex := regexp.MustCompile(`(?i)\b` + keyword + `\b`)
		result = regex.ReplaceAllString(result, "")
	}

	// Limit query length to shorter max
	if len(result) > MaxQueryLength {
		result = result[:MaxQueryLength]
	}

	return strings.TrimSpace(result)
}

// SanitizeFilterValue sanitizes a filter value used in API queries
func SanitizeFilterValue(value interface{}) string {
	if value == nil {
		return ""
	}

	str := fmt.Sprintf("%v", value)

	// Remove potential injection attempts
	str = strings.ReplaceAll(str, ";", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\t", "")

	// Remove potential command injection
	str = strings.ReplaceAll(str, "|", "")
	str = strings.ReplaceAll(str, "&", "")
	str = strings.ReplaceAll(str, "$", "")
	str = strings.ReplaceAll(str, "`", "")
	str = strings.ReplaceAll(str, "$(", "")
	str = strings.ReplaceAll(str, "${", "")

	// Limit length
	if len(str) > MaxStringLength {
		str = str[:MaxStringLength]
	}

	return strings.TrimSpace(str)
}

// SanitizeIdentifier sanitizes an identifier (like a name, ID, etc.)
func SanitizeIdentifier(input string) string {
	if input == "" {
		return input
	}

	// Only allow alphanumeric, hyphens, underscores, and dots
	var result strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SanitizeJSONPath sanitizes a JSONPath expression
func SanitizeJSONPath(input string) string {
	if input == "" {
		return input
	}

	// Remove any non-JSONPath characters
	// JSONPath allows: alphanumeric, ., [, ], {, }, :, @, -, _, space, and quotes
	allowedChars := regexp.MustCompile(`[a-zA-Z0-9\.\[\]\{\}\:\-\_\s'\""]+`)
	result := allowedChars.FindAllString(input, -1)

	return strings.Join(result, "")
}

// ValidateString validates a string input against security rules
func ValidateString(input string, maxLength int) error {
	if len(input) > maxLength {
		return fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}

	// Check for dangerous patterns
	if strings.ContainsAny(input, "\x00\x1a\x1b\x1c\x1d\x1e\x1f") {
		return fmt.Errorf("input contains invalid control characters")
	}

	return nil
}

// ValidateIdentifier validates an identifier (name, ID, etc.)
func ValidateIdentifier(input string) error {
	if input == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	// Only allow alphanumeric, hyphens, underscores, and dots
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, input)
	if !matched {
		return fmt.Errorf("identifier contains invalid characters")
	}

	return nil
}
