package config

import (
	"fmt"
	"strconv"

	"github.com/JorgeSaicoski/microservice-commons/utils"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
	TimeZone     string
	MaxIdleConns int
	MaxOpenConns int
	LogLevel     string
}

// LoadDatabaseConfig loads database configuration from environment
func LoadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:         utils.GetEnv("POSTGRES_HOST", "localhost"),
		Port:         utils.GetEnv("POSTGRES_PORT", "5432"),
		User:         utils.GetEnv("POSTGRES_USER", "postgres"),
		Password:     utils.GetEnv("POSTGRES_PASSWORD", "postgres"),
		DatabaseName: utils.GetEnv("POSTGRES_DB", "defaultdb"),
		SSLMode:      utils.GetEnv("POSTGRES_SSLMODE", "disable"),
		TimeZone:     utils.GetEnv("POSTGRES_TIMEZONE", "UTC"),
		MaxIdleConns: utils.GetEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
		MaxOpenConns: utils.GetEnvInt("POSTGRES_MAX_OPEN_CONNS", 100),
		LogLevel:     utils.GetEnv("POSTGRES_LOG_LEVEL", "silent"),
	}
}

// ConnectionString returns the PostgreSQL connection string
func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		dc.Host, dc.Port, dc.User, dc.Password, dc.DatabaseName, dc.SSLMode, dc.TimeZone,
	)
}

// Validate validates the database configuration
func (dc *DatabaseConfig) Validate() error {
	if dc.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if dc.Port == "" {
		return fmt.Errorf("database port is required")
	}

	// Validate port is a number
	if _, err := strconv.Atoi(dc.Port); err != nil {
		return fmt.Errorf("database port must be a number: %w", err)
	}

	if dc.User == "" {
		return fmt.Errorf("database user is required")
	}

	if dc.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if dc.DatabaseName == "" {
		return fmt.Errorf("database name is required")
	}

	if dc.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}

	if dc.MaxOpenConns < 0 {
		return fmt.Errorf("max open connections cannot be negative")
	}

	if dc.MaxOpenConns > 0 && dc.MaxIdleConns > dc.MaxOpenConns {
		return fmt.Errorf("max idle connections cannot exceed max open connections")
	}

	return nil
}

// IsSSLEnabled returns true if SSL is enabled
func (dc *DatabaseConfig) IsSSLEnabled() bool {
	return dc.SSLMode != "disable"
}

// GetLogLevel returns the database log level
func (dc *DatabaseConfig) GetLogLevel() string {
	validLevels := []string{"silent", "error", "warn", "info"}
	for _, level := range validLevels {
		if dc.LogLevel == level {
			return dc.LogLevel
		}
	}
	return "silent"
}
