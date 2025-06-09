// config/validation.go
package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// Validator provides configuration validation utilities
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// ValidateRequired validates that a field is not empty
func (v *Validator) ValidateRequired(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "field is required",
		})
	}
	return v
}

// ValidatePort validates that a value is a valid port number
func (v *Validator) ValidatePort(field, value string) *Validator {
	if value == "" {
		return v
	}

	port, err := strconv.Atoi(value)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "must be a valid port number",
		})
		return v
	}

	if port < 1 || port > 65535 {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "port must be between 1 and 65535",
		})
	}

	return v
}

// ValidateURL validates that a value is a valid URL
func (v *Validator) ValidateURL(field, value string) *Validator {
	if value == "" {
		return v
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "must be a valid URL",
		})
		return v
	}

	if parsedURL.Scheme == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "URL must include scheme (http:// or https://)",
		})
	}

	return v
}

// ValidateOneOf validates that a value is one of the allowed values
func (v *Validator) ValidateOneOf(field, value string, allowed []string) *Validator {
	if value == "" {
		return v
	}

	for _, allowedValue := range allowed {
		if value == allowedValue {
			return v
		}
	}

	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Value:   value,
		Message: fmt.Sprintf("must be one of: %s", strings.Join(allowed, ", ")),
	})

	return v
}

// ValidateMinMax validates that a numeric value is within min/max range
func (v *Validator) ValidateMinMax(field, value string, min, max int) *Validator {
	if value == "" {
		return v
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: "must be a valid number",
		})
		return v
	}

	if num < min || num > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Value:   value,
			Message: fmt.Sprintf("must be between %d and %d", min, max),
		})
	}

	return v
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidationError {
	return v.errors
}

// Error returns a formatted error message with all validation errors
func (v *Validator) Error() error {
	if !v.HasErrors() {
		return nil
	}

	var messages []string
	for _, err := range v.errors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(messages, "\n  - "))
}

// ValidateConfig validates a complete configuration
func ValidateConfig(config *Config) error {
	validator := NewValidator()

	// Validate basic fields
	validator.ValidateRequired("PORT", config.Port)
	validator.ValidatePort("PORT", config.Port)
	validator.ValidateRequired("SERVICE_NAME", config.ServiceName)
	validator.ValidateOneOf("ENVIRONMENT", config.Environment, []string{"dev", "development", "staging", "prod", "production"})
	validator.ValidateOneOf("LOG_LEVEL", config.LogLevel, []string{"debug", "info", "warn", "error"})

	// Validate allowed origins
	for i, origin := range config.AllowedOrigins {
		field := fmt.Sprintf("ALLOWED_ORIGINS[%d]", i)
		validator.ValidateURL(field, origin)
	}

	return validator.Error()
}
