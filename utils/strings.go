// utils/strings.go
package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty checks if a string is not empty and contains non-whitespace characters
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// TrimAndLower trims whitespace and converts to lowercase
func TrimAndLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// TrimAndUpper trims whitespace and converts to uppercase
func TrimAndUpper(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsIgnoreCase checks if a string slice contains a specific string (case-insensitive)
func ContainsIgnoreCase(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}

// SplitAndTrim splits a string and trims whitespace from each part
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// Join joins non-empty strings with a separator
func Join(parts []string, sep string) string {
	nonEmpty := make([]string, 0, len(parts))

	for _, part := range parts {
		if IsNotEmpty(part) {
			nonEmpty = append(nonEmpty, strings.TrimSpace(part))
		}
	}

	return strings.Join(nonEmpty, sep)
}

// Truncate truncates a string to a maximum length
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// TruncateWithEllipsis truncates a string and adds ellipsis if truncated
func TruncateWithEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	if maxLength <= 3 {
		return s[:maxLength]
	}

	return s[:maxLength-3] + "..."
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CamelCase converts a string to camelCase
func CamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += Capitalize(strings.ToLower(words[i]))
	}

	return result
}

// PascalCase converts a string to PascalCase
func PascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	var result string
	for _, word := range words {
		result += Capitalize(strings.ToLower(word))
	}

	return result
}

// SnakeCase converts a string to snake_case
func SnakeCase(s string) string {
	// Insert underscores before uppercase letters (except the first character)
	var result []rune
	runes := []rune(s)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}

	// Replace spaces and other separators with underscores
	str := string(result)
	str = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(str, "_")

	// Remove leading/trailing underscores
	return strings.Trim(str, "_")
}

// KebabCase converts a string to kebab-case
func KebabCase(s string) string {
	// Insert hyphens before uppercase letters (except the first character)
	var result []rune
	runes := []rune(s)

	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '-')
		}
		result = append(result, unicode.ToLower(r))
	}

	// Replace spaces and other separators with hyphens
	str := string(result)
	str = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(str, "-")

	// Remove leading/trailing hyphens
	return strings.Trim(str, "-")
}

// RemoveSpecialChars removes special characters, keeping only letters, numbers, and spaces
func RemoveSpecialChars(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	return reg.ReplaceAllString(s, "")
}

// RemoveNonAlphanumeric removes all non-alphanumeric characters
func RemoveNonAlphanumeric(s string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return reg.ReplaceAllString(s, "")
}

// Reverse reverses a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Pad pads a string to a specific length with spaces
func Pad(s string, length int) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(" ", length-len(s))
	return s + padding
}

// PadLeft pads a string to a specific length with spaces on the left
func PadLeft(s string, length int) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(" ", length-len(s))
	return padding + s
}

// PadCenter centers a string within a specific length
func PadCenter(s string, length int) string {
	if len(s) >= length {
		return s
	}

	totalPadding := length - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
}

// SlugFromString creates a URL-friendly slug from a string
func SlugFromString(s string) string {
	// Convert to lowercase
	slug := strings.ToLower(s)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}

// WordCount counts words in a string
func WordCount(s string) int {
	words := strings.Fields(s)
	return len(words)
}

// CharCount counts characters in a string (excluding whitespace)
func CharCount(s string) int {
	count := 0
	for _, r := range s {
		if !unicode.IsSpace(r) {
			count++
		}
	}
	return count
}
