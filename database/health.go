// database/health.go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/JorgeSaicoski/pgconnect"
)

// HealthStatus represents database health status
type HealthStatus struct {
	Status       string        `json:"status"`
	Database     string        `json:"database"`
	ResponseTime time.Duration `json:"responseTime"`
	Error        string        `json:"error,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// HealthChecker handles database health checks
type HealthChecker struct {
	db      *pgconnect.DB
	timeout time.Duration
}

// NewHealthChecker creates a new database health checker
func NewHealthChecker(db *pgconnect.DB) *HealthChecker {
	return &HealthChecker{
		db:      db,
		timeout: 5 * time.Second, // Default timeout
	}
}

// SetTimeout sets the health check timeout
func (hc *HealthChecker) SetTimeout(timeout time.Duration) *HealthChecker {
	hc.timeout = timeout
	return hc
}

// Check performs a database health check
func (hc *HealthChecker) Check() HealthStatus {
	start := time.Now()
	status := HealthStatus{
		Timestamp: start,
		Database:  "postgres",
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), hc.timeout)
	defer cancel()

	// Perform health check
	err := hc.performCheck(ctx)

	status.ResponseTime = time.Since(start)

	if err != nil {
		status.Status = "unhealthy"
		status.Error = err.Error()
	} else {
		status.Status = "healthy"
	}

	return status
}

// performCheck executes the actual health check
func (hc *HealthChecker) performCheck(ctx context.Context) error {
	// Check if database connection is alive
	if err := hc.db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// Perform a simple query to ensure database is responsive
	var result int
	query := "SELECT 1"

	// Execute query with context timeout
	done := make(chan error, 1)
	go func() {
		err := hc.db.DB.WithContext(ctx).Raw(query).Scan(&result).Error
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		if result != 1 {
			return fmt.Errorf("unexpected query result: %d", result)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("health check timeout")
	}
}

// IsHealthy returns true if database is healthy
func (hc *HealthChecker) IsHealthy() bool {
	status := hc.Check()
	return status.Status == "healthy"
}

// CheckWithDetails returns detailed health information
func (hc *HealthChecker) CheckWithDetails() (HealthStatus, map[string]interface{}) {
	status := hc.Check()

	details := map[string]interface{}{
		"timeout":    hc.timeout.String(),
		"database":   status.Database,
		"checked_at": status.Timestamp.Format(time.RFC3339),
	}

	// Add connection pool stats if available
	if sqlDB, err := hc.db.DB.DB(); err == nil {
		stats := sqlDB.Stats()
		details["connections"] = map[string]interface{}{
			"open":     stats.OpenConnections,
			"in_use":   stats.InUse,
			"idle":     stats.Idle,
			"max_open": stats.MaxOpenConnections,
		}
	}

	return status, details
}

// QuickHealthCheck is a convenience function for simple health checks
func QuickHealthCheck(db *pgconnect.DB) bool {
	checker := NewHealthChecker(db)
	return checker.IsHealthy()
}

// DetailedHealthCheck returns comprehensive health information
func DetailedHealthCheck(db *pgconnect.DB) (HealthStatus, map[string]interface{}) {
	checker := NewHealthChecker(db)
	return checker.CheckWithDetails()
}
