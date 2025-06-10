package test

import (
	"os"
	"testing"
	"time"

	"github.com/JorgeSaicoski/microservice-commons/utils"
)

// Test environment utilities
func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := utils.GetEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got %s", result)
	}

	// Test with fallback
	result = utils.GetEnv("NON_EXISTENT", "fallback")
	if result != "fallback" {
		t.Errorf("Expected 'fallback', got %s", result)
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("INT_VAR", "42")
	defer os.Unsetenv("INT_VAR")

	result := utils.GetEnvInt("INT_VAR", 0)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}

	// Test fallback
	result = utils.GetEnvInt("NON_EXISTENT", 99)
	if result != 99 {
		t.Errorf("Expected 99, got %d", result)
	}
}

// Test string utilities
func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, test := range tests {
		result := utils.IsEmpty(test.input)
		if result != test.expected {
			t.Errorf("IsEmpty(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestSplitAndTrim(t *testing.T) {
	result := utils.SplitAndTrim("a, b , c,d", ",")
	expected := []string{"a", "b", "c", "d"}

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
		}
	}
}

func TestSlugFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test@#$%String", "test-string"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"  Leading and trailing  ", "leading-and-trailing"},
	}

	for _, test := range tests {
		result := utils.SlugFromString(test.input)
		if result != test.expected {
			t.Errorf("SlugFromString(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// Test time utilities
func TestFormatDate(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	result := utils.FormatDate(testTime)
	expected := "2023-12-25"

	if result != expected {
		t.Errorf("FormatDate() = %s, expected %s", result, expected)
	}
}

func TestStartOfDay(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 123456789, time.UTC)
	result := utils.StartOfDay(testTime)

	expected := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)

	if !result.Equal(expected) {
		t.Errorf("StartOfDay() = %v, expected %v", result, expected)
	}
}

func TestIsBusinessDay(t *testing.T) {
	// Monday
	monday := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	if !utils.IsBusinessDay(monday) {
		t.Error("Monday should be a business day")
	}

	// Saturday
	saturday := time.Date(2023, 12, 23, 0, 0, 0, 0, time.UTC)
	if utils.IsBusinessDay(saturday) {
		t.Error("Saturday should not be a business day")
	}
}

// Test validation utilities
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
	}

	for _, test := range tests {
		result := utils.IsValidEmail(test.email)
		if result != test.valid {
			t.Errorf("IsValidEmail(%s) = %v, expected %v", test.email, result, test.valid)
		}
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		password string
		strong   bool
	}{
		{"Password123!", true},
		{"Weak123", false},      // No special char
		{"password123!", false}, // No uppercase
		{"PASSWORD123!", false}, // No lowercase
		{"Password!", false},    // No digit
		{"Pass1!", false},       // Too short
	}

	for _, test := range tests {
		result := utils.IsStrongPassword(test.password)
		if result != test.strong {
			t.Errorf("IsStrongPassword(%s) = %v, expected %v", test.password, result, test.strong)
		}
	}
}
