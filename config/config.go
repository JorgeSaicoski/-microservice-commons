// config/config.go
package config

import (
	"fmt"
	"strings"

	"github.com/JorgeSaicoski/microservice-commons/utils"
)

// Config holds common configuration for microservices
// Config holds common configuration for microservices
type Config struct {
	Port           string
	AllowedOrigins []string
	DatabaseConfig DatabaseConfig
	KeycloakConfig KeycloakConfig
	LogLevel       string
	Environment    string
	ServiceName    string
	ServiceVersion string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	return &Config{
		Port:           utils.GetEnv("PORT", "8000"),
		AllowedOrigins: parseOrigins(utils.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
		DatabaseConfig: LoadDatabaseConfig(),
		KeycloakConfig: LoadKeycloakConfig(),
		LogLevel:       utils.GetEnv("LOG_LEVEL", "info"),
		Environment:    utils.GetEnv("ENVIRONMENT", "dev"),
		ServiceName:    utils.GetEnv("SERVICE_NAME", "microservice"),
		ServiceVersion: utils.GetEnv("SERVICE_VERSION", "1.0.0"),
	}
}

// LoadWithServiceInfo loads config with service-specific information
func LoadWithServiceInfo(serviceName, version string) *Config {
	config := LoadFromEnv()

	// Override with provided values if not set in environment
	if config.ServiceName == "microservice" {
		config.ServiceName = serviceName
	}
	if config.ServiceVersion == "1.0.0" {
		config.ServiceVersion = version
	}

	return config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}

	if err := c.DatabaseConfig.Validate(); err != nil {
		return fmt.Errorf("database config: %w", err)
	}

	if err := c.KeycloakConfig.Validate(); err != nil {
		return fmt.Errorf("keycloak config: %w", err)
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev" || c.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "prod" || c.Environment == "production"
}

// IsStaging returns true if running in staging environment
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// GetLogLevel returns the log level for the application
func (c *Config) GetLogLevel() LogLevel {
	switch strings.ToLower(c.LogLevel) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// parseOrigins parses the origins string into a slice
func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}

	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return []string{"http://localhost:3000"}
	}

	return result
}

// LogLevel represents logging levels
type LogLevel string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)
