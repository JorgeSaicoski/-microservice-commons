// utils/validation.go
package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/mail"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

// Validator provides validation utilities
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error
func (v *Validator) AddError(field, value, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// Clear clears all validation errors
func (v *Validator) Clear() {
	v.errors = make([]ValidationError, 0)
}

// Email validation
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// URL validation
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Phone number validation (basic)
func IsValidPhone(phone string) bool {
	// Remove common separators
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	// Check if all remaining characters are digits
	if len(cleaned) < 10 || len(cleaned) > 15 {
		return false
	}

	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

// Password strength validation
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// Username validation
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	// Allow letters, numbers, underscore, and hyphen
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

// Slug validation
func IsValidSlug(slug string) bool {
	if len(slug) == 0 {
		return false
	}

	// Allow lowercase letters, numbers, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, slug)
	return matched && !strings.HasPrefix(slug, "-") && !strings.HasSuffix(slug, "-")
}

// UUID validation
func IsValidUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// Credit card validation (basic Luhn algorithm)
func IsValidCreditCard(number string) bool {
	// Remove spaces and hyphens
	cleaned := strings.ReplaceAll(number, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// Check if all characters are digits
	for _, r := range cleaned {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	// Check length (13-19 digits for most cards)
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	alternate := false

	for i := len(cleaned) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cleaned[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// IP address validation
func IsValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// Port validation
func IsValidPort(port string) bool {
	num, err := strconv.Atoi(port)
	return err == nil && num > 0 && num <= 65535
}

// Alphanumeric validation
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Alphabetic validation
func IsAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Numeric validation
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Integer validation
func IsInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// Length validation
func IsValidLength(s string, min, max int) bool {
	length := len(s)
	return length >= min && length <= max
}

// Range validation for numbers
func IsInRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// File extension validation
func HasValidExtension(filename string, allowedExtensions []string) bool {
	if len(allowedExtensions) == 0 {
		return true
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))

	for _, allowed := range allowedExtensions {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}

	return false
}

// JSON validation
func IsValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// Base64 validation
func IsValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// Hex validation
func IsValidHex(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}

// Date validation
func IsValidDate(dateStr, format string) bool {
	_, err := time.Parse(format, dateStr)
	return err == nil
}

// Required field validation
func IsRequired(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return IsNotEmpty(v)
	case []string:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	default:
		return true
	}
}

// Custom validation functions
type ValidatorFunc func(value interface{}) bool

// ValidateField validates a single field with multiple rules
func ValidateField(field, value string, rules ...func(string) bool) []ValidationError {
	var errors []ValidationError

	for _, rule := range rules {
		if !rule(value) {
			errors = append(errors, ValidationError{
				Field:   field,
				Value:   value,
				Message: "Validation failed",
			})
		}
	}

	return errors
}

// Common validation rule creators
func MinLength(min int) func(string) bool {
	return func(s string) bool {
		return len(s) >= min
	}
}

func MaxLength(max int) func(string) bool {
	return func(s string) bool {
		return len(s) <= max
	}
}

func ExactLength(length int) func(string) bool {
	return func(s string) bool {
		return len(s) == length
	}
}

func MatchesPattern(pattern string) func(string) bool {
	regex := regexp.MustCompile(pattern)
	return func(s string) bool {
		return regex.MatchString(s)
	}
}

func NotEmpty() func(string) bool {
	return func(s string) bool {
		return IsNotEmpty(s)
	}
}
