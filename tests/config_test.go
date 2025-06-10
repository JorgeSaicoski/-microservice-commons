package test

import (
	"os"
	"testing"

	"github.com/JorgeSaicoski/microservice-commons/config"
)

func TestLoadFromEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"PORT", "SERVICE_NAME", "ENVIRONMENT", "LOG_LEVEL", "ALLOWED_ORIGINS"}
	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
	}

	// Clean up after test
	defer func() {
		for env, value := range originalEnv {
			if value == "" {
				os.Unsetenv(env)
			} else {
				os.Setenv(env, value)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("PORT", "9000")
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ALLOWED_ORIGINS", "http://test1.com,http://test2.com")

	config := config.LoadFromEnv()

	if config.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", config.Port)
	}

	if config.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got %s", config.ServiceName)
	}

	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got %s", config.Environment)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.LogLevel)
	}

	expectedOrigins := []string{"http://test1.com", "http://test2.com"}
	if len(config.AllowedOrigins) != 2 {
		t.Errorf("Expected 2 origins, got %d", len(config.AllowedOrigins))
	}

	for i, origin := range expectedOrigins {
		if config.AllowedOrigins[i] != origin {
			t.Errorf("Expected origin %s, got %s", origin, config.AllowedOrigins[i])
		}
	}
}

func TestLoadWithServiceInfo(t *testing.T) {
	config := config.LoadWithServiceInfo("my-service", "2.0.0")

	if config.ServiceName != "my-service" {
		t.Errorf("Expected service name 'my-service', got %s", config.ServiceName)
	}

	if config.ServiceVersion != "2.0.0" {
		t.Errorf("Expected service version '2.0.0', got %s", config.ServiceVersion)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      config.Config
		expectError bool
	}{
		{
			name: "valid config",
			config: config.Config{
				Port:        "8080",
				ServiceName: "test-service",
				DatabaseConfig: config.DatabaseConfig{
					Host:         "localhost",
					Port:         "5432",
					User:         "user",
					Password:     "pass",
					DatabaseName: "db",
				},
				KeycloakConfig: config.KeycloakConfig{
					PublicKeyBase64: "test-key",
				},
			},
			expectError: false,
		},
		{
			name: "missing port",
			config: config.Config{
				ServiceName: "test-service",
			},
			expectError: true,
		},
		{
			name: "missing service name",
			config: config.Config{
				Port: "8080",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestEnvironmentCheckers(t *testing.T) {
	tests := []struct {
		env          string
		isDev        bool
		isProduction bool
		isStaging    bool
	}{
		{"dev", true, false, false},
		{"development", true, false, false},
		{"prod", false, true, false},
		{"production", false, true, false},
		{"staging", false, false, true},
		{"test", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			config := &config.Config{Environment: tt.env}

			if config.IsDevelopment() != tt.isDev {
				t.Errorf("IsDevelopment() = %v, want %v", config.IsDevelopment(), tt.isDev)
			}

			if config.IsProduction() != tt.isProduction {
				t.Errorf("IsProduction() = %v, want %v", config.IsProduction(), tt.isProduction)
			}

			if config.IsStaging() != tt.isStaging {
				t.Errorf("IsStaging() = %v, want %v", config.IsStaging(), tt.isStaging)
			}
		})
	}
}
