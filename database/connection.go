package database

import (
	"fmt"
	"time"

	"github.com/JorgeSaicoski/microservice-commons/config"
	"github.com/JorgeSaicoski/microservice-commons/utils"
	"github.com/JorgeSaicoski/pgconnect"
	"gorm.io/gorm/logger"
)

// ConnectionManager manages database connections
type ConnectionManager struct {
	db     *pgconnect.DB
	config config.DatabaseConfig
}

// NewConnectionManager creates a new database connection manager
func NewConnectionManager(cfg config.DatabaseConfig) *ConnectionManager {
	return &ConnectionManager{
		config: cfg,
	}
}

// Connect establishes database connection with retry logic
func (cm *ConnectionManager) Connect() (*pgconnect.DB, error) {
	pgConfig := pgconnect.Config{
		Host:         cm.config.Host,
		Port:         cm.config.Port,
		User:         cm.config.User,
		Password:     cm.config.Password,
		DatabaseName: cm.config.DatabaseName,
		SSLMode:      cm.config.SSLMode,
		TimeZone:     cm.config.TimeZone,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogLevel:     logger.Silent, // Override with environment if needed
	}

	// Set log level based on environment
	if logLevel := utils.GetEnv("DB_LOG_LEVEL", "silent"); logLevel != "silent" {
		switch logLevel {
		case "error":
			pgConfig.LogLevel = logger.Error
		case "warn":
			pgConfig.LogLevel = logger.Warn
		case "info":
			pgConfig.LogLevel = logger.Info
		default:
			pgConfig.LogLevel = logger.Silent
		}
	}

	// Connection retry configuration
	maxRetries := utils.GetEnvInt("DB_MAX_RETRIES", 3)
	retryDelay := time.Duration(utils.GetEnvInt("DB_RETRY_DELAY_SECONDS", 30)) * time.Second

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempting to connect to database (attempt %d of %d)\n", attempt, maxRetries)

		cm.db, err = pgconnect.New(pgConfig)
		if err == nil {
			fmt.Println("Successfully connected to database")
			return cm.db, nil
		}

		fmt.Printf("Failed to connect to database: %v\n", err)

		if attempt < maxRetries {
			fmt.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// GetConnection returns the current database connection
func (cm *ConnectionManager) GetConnection() *pgconnect.DB {
	return cm.db
}

// Close closes the database connection
func (cm *ConnectionManager) Close() error {
	if cm.db != nil {
		return cm.db.Close()
	}
	return nil
}

// Reconnect attempts to reconnect to the database
func (cm *ConnectionManager) Reconnect() error {
	if cm.db != nil {
		cm.db.Close()
	}

	var err error
	cm.db, err = cm.Connect()
	return err
}

// ConnectWithConfig is a convenience function for quick database connection
func ConnectWithConfig(cfg config.DatabaseConfig) (*pgconnect.DB, error) {
	manager := NewConnectionManager(cfg)
	return manager.Connect()
}

// MustConnect connects to database or panics on failure
func MustConnect(cfg config.DatabaseConfig) *pgconnect.DB {
	db, err := ConnectWithConfig(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}
