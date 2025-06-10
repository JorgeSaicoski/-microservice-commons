package test

import (
	"testing"
	"time"

	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/JorgeSaicoski/microservice-commons/database"
)

func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      config.DatabaseConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: config.DatabaseConfig{
				Host:         "localhost",
				Port:         "5432",
				User:         "user",
				Password:     "pass",
				DatabaseName: "db",
				MaxIdleConns: 10,
				MaxOpenConns: 100,
			},
			expectError: false,
		},
		{
			name: "missing host",
			config: config.DatabaseConfig{
				Port:         "5432",
				User:         "user",
				Password:     "pass",
				DatabaseName: "db",
			},
			expectError: true,
		},
		{
			name: "invalid port",
			config: config.DatabaseConfig{
				Host:         "localhost",
				Port:         "invalid",
				User:         "user",
				Password:     "pass",
				DatabaseName: "db",
			},
			expectError: true,
		},
		{
			name: "negative max idle conns",
			config: config.DatabaseConfig{
				Host:         "localhost",
				Port:         "5432",
				User:         "user",
				Password:     "pass",
				DatabaseName: "db",
				MaxIdleConns: -1,
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

func TestHealthChecker(t *testing.T) {

	checker := database.NewHealthChecker(nil) // nil DB for testing
	checker.SetTimeout(1 * time.Second)

	status := checker.Check()

	if status.Database != "postgres" {
		t.Errorf("Expected database type 'postgres', got %s", status.Database)
	}

	if status.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}
