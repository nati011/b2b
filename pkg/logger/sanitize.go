package logger

import (
	"regexp"
	"strings"

	fluent_validator "marketplace/pkg/validate"
)

var (
	// Sensitive field patterns
	sensitiveFields = map[string]bool{
		"password":        true,
		"passwd":          true,
		"secret":          true,
		"token":           true,
		"apikey":          true,
		"api_key":         true,
		"access_token":    true,
		"auth":            true,
		"authorization":   true,
		"credit_card":     true,
		"card_number":     true,
		"ssn":             true,
		"social_security": true,
	}

	// Patterns for sensitive data - using regex to find potential matches, then validate
	// Email pattern to find potential email addresses in strings
	emailPattern = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	// Credit card pattern (basic format detection)
	creditCardPattern = regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`)
	// SSN pattern (format: XXX-XX-XXXX)
	ssnPattern = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
)

const redactedValue = "[REDACTED]"

// SanitizeFields processes a map of fields and redacts sensitive values
// Only applies sanitization if enabled in configuration
func SanitizeFields(fields map[string]any) map[string]any {
	// Skip sanitization if not enabled
	if !sanitizeEnabled.Load() {
		return fields
	}

	sanitized := make(map[string]any, len(fields))
	for k, v := range fields {
		key := strings.ToLower(k)
		if isSensitiveField(key) {
			sanitized[k] = redactedValue
		} else {
			sanitized[k] = sanitizeValue(v)
		}
	}
	return sanitized
}

// Sanitize processes log arguments and redacts sensitive data
// Only applies sanitization if enabled in configuration
func Sanitize(args ...any) []any {
	// Skip sanitization if not enabled
	if !sanitizeEnabled.Load() {
		return args
	}

	if len(args) == 0 {
		return args
	}

	sanitized := make([]any, 0, len(args))
	for i := 0; i < len(args); i++ {
		if i+1 < len(args) {
			// Check if this is a key-value pair
			if key, ok := args[i].(string); ok {
				keyLower := strings.ToLower(key)
				if isSensitiveField(keyLower) {
					sanitized = append(sanitized, key, redactedValue)
					i++ // Skip the value
					continue
				}
				// Sanitize the value
				sanitized = append(sanitized, key, sanitizeValue(args[i+1]))
				i++ // Skip the value
				continue
			}
		}
		// Not a key-value pair, sanitize the value itself
		sanitized = append(sanitized, sanitizeValue(args[i]))
	}
	return sanitized
}

// sanitizeValue recursively sanitizes a value
func sanitizeValue(v any) any {
	switch val := v.(type) {
	case string:
		return sanitizeString(val)
	case map[string]any:
		return SanitizeFields(val)
	case []any:
		sanitized := make([]any, len(val))
		for i, item := range val {
			sanitized[i] = sanitizeValue(item)
		}
		return sanitized
	default:
		return v
	}
}

// sanitizeString redacts sensitive patterns in strings
// Only applies sanitization if enabled in configuration
func sanitizeString(s string) string {
	// Skip sanitization if not enabled
	if !sanitizeEnabled.Load() {
		return s
	}

	// Redact email addresses using validate package for validation
	s = emailPattern.ReplaceAllStringFunc(s, func(match string) string {
		// Use validate package to verify it's actually a valid email
		result := fluent_validator.New().
			And(fluent_validator.EmailValid(match)).
			Validate()
		if result.IsValid {
			return redactedValue
		}
		return match
	})

	// Redact credit card numbers (basic pattern matching)
	// Note: Credit card validation would require Luhn algorithm check
	// For now, we redact based on pattern format
	s = creditCardPattern.ReplaceAllString(s, redactedValue)

	// Redact SSN numbers (format: XXX-XX-XXXX)
	s = ssnPattern.ReplaceAllString(s, redactedValue)

	return s
}

// isSensitiveField checks if a field name indicates sensitive data
func isSensitiveField(fieldName string) bool {
	fieldName = strings.ToLower(fieldName)
	for sensitive := range sensitiveFields {
		if strings.Contains(fieldName, sensitive) {
			return true
		}
	}
	return false
}
